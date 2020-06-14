package database

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"strings"
	"time"
)

// Database type acts as the control interface for the Redis datastore.
type Database struct {
	uri    string
	client *redis.Client
}

// New creates a database connection.
func New(uri string) Database {
	if strings.HasPrefix(uri, "redis://") {
		uri = uri[8:]
	}

	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to database: %s due to err: %s", uri, err)
	}

	log.Printf("Connected to database at %s", uri)

	return Database{
		uri:    uri,
		client: client,
	}
}

// GetBestCandidateWord returns the best replacement synonym for supplied word.
// The return order should be in the defined lexeme order in Lexeme.go.
func (d Database) GetBestCandidateWord(word string) string {
	r := make([]*redis.StringCmd, 0, 4)
	results := &r

	_, err := d.client.TxPipelined(func(pipe redis.Pipeliner) error {
		pipe.Expire(fmt.Sprintf("best_word_single_%s", word), 10*time.Second)

		for _, l := range ordering {
			res := pipe.SRandMember(fmt.Sprintf("%s:%s", l, word))
			*results = append(*results, res)
		}

		return nil
	})

	if err != nil {
		log.Printf("Could not access datastore for word: %s, %s", word, err)
		return word
	}

	for _, c := range *results {
		w, err := c.Result()
		if err == nil {
			return w
		}
	}

	// Fallback. Only return if nothing was found.
	return word
}
