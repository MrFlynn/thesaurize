package loader

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/MrFlynn/thesaurize/internal/database"
	"github.com/urfave/cli/v2"
)

// Load loads data into a Redis database from a source thesaurus file.
func Load(ctx *cli.Context) error {
	var (
		uri = ctx.String("data")

		rd  io.ReadCloser
		err error
	)

	switch parts := strings.SplitN(uri, "://", 2); parts[0] {
	case "file":
		rd, err = os.Open(parts[1])
	case "https", "http":
		var resp *http.Response

		resp, err = http.Get(uri)
		if resp.StatusCode == http.StatusNotFound {
			err = fmt.Errorf("unable to get file %s", uri)
		}

		rd = resp.Body
	default:
		err = fmt.Errorf("unknown protocol %s", parts[0])
	}

	if err != nil {
		return err
	}

	defer rd.Close()

	dataFile, err := getDataReaderFromZip(rd)
	if err != nil {
		return err
	}

	defer dataFile.Close()

	var (
		ch = make(chan entry)
		wg sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := pushToRedis(database.New(ctx.String("datastore")), ch, 500); err != nil {
			log.Fatalf("Unable to push data to redis: %s", err)
		}
	}()

	var filter *profanityFilter

	if ctx.Bool("skip-profane-words") {
		categories := make(map[string]struct{}, len(ctx.StringSlice("profane-word-categories")))
		for _, category := range ctx.StringSlice("profane-word-categories") {
			categories[category] = struct{}{}
		}

		filter = &profanityFilter{categories: categories}
		if err := filter.init(ctx.String("profane-word-index-url")); err != nil {
			log.Fatalf("Unable to initialize profanity filter: %s", err)
		}
	}

	log.Println("Loading dataset into Redis backend")

	if err := scanDataFile(dataFile, ch, filter); err != nil {
		log.Fatalf("Unable to read data file: %s", err)
	}

	wg.Wait()

	log.Println("Loading complete")
	return nil
}

func getDataReaderFromZip(file io.Reader) (io.ReadCloser, error) {
	buff := bytes.NewBuffer([]byte{})
	sz, err := io.Copy(buff, file)
	if err != nil {
		return nil, err
	}

	zipFile, err := zip.NewReader(bytes.NewReader(buff.Bytes()), sz)
	if err != nil {
		return nil, err
	}

	for _, file := range zipFile.File {
		if strings.HasSuffix(file.Name, ".dat") {
			return file.Open()
		}
	}

	return nil, errors.New("thesaurus data file not present in zip archive")
}

type entry struct {
	key    string
	values []string
}

func pushToRedis(db database.Database, in chan entry, queueSize int) error {
	var (
		count    int
		pipeline = db.GetPipeline()
	)

	for e := range in {
		if count >= queueSize {
			if _, err := pipeline.Exec(); err != nil {
				return err
			}

			pipeline = db.GetPipeline()
			count = 0
		}

		pipeline.SAdd(e.key, e.values)
		count++
	}

	if _, err := pipeline.Exec(); err != nil {
		return err
	}

	return db.SendReady()
}

func scanDataFile(rd io.Reader, out chan entry, filter *profanityFilter) error {
	defer close(out)

	scanner := bufio.NewScanner(rd)
	if !scanner.Scan() {
		return scanner.Err()
	}

	for scanner.Scan() {
		word, synonyms, err := readSynonyms(scanner, filter)
		if err != nil {
			log.Printf("Unable to get synonyms for '%s': %s", word, err)
		}

		for lexeme, syns := range synonyms {
			out <- entry{key: lexeme + ":" + word, values: syns}
		}
	}

	return scanner.Err()
}

func readSynonyms(scanner *bufio.Scanner, filter *profanityFilter) (string, map[string][]string, error) {
	wordHeader := strings.SplitN(scanner.Text(), "|", 2)
	if fieldCount := len(wordHeader); fieldCount < 2 {
		return "", nil, fmt.Errorf("invalid header, expected 2 fields, got %d", fieldCount)
	}

	rowCount, err := strconv.Atoi(wordHeader[1])
	if err != nil {
		return "", nil, fmt.Errorf(
			"invalid row count in header for word '%s', %s is not a valid number",
			wordHeader[0],
			wordHeader[1],
		)
	}

	var skip bool
	if filter != nil {
		skip = filter.match(wordHeader[0])
	}

	synonyms := make(map[string][]string, rowCount)

	for i := 0; i < rowCount && scanner.Scan(); i++ {
		if skip {
			continue
		}

		rowFields := strings.Split(scanner.Text(), "|")
		if len(rowFields) < 2 {
			continue
		}

		lexeme := strings.Trim(rowFields[0], "()")
		for _, synonym := range rowFields[1:] {
			if filter != nil && filter.match(synonym) {
				continue
			}

			if synWords, ok := synonyms[lexeme]; ok {
				synonyms[lexeme] = append(synWords, synonym)
			} else {
				synonyms[lexeme] = []string{synonym}
			}
		}
	}

	return string(wordHeader[0]), synonyms, scanner.Err()
}
