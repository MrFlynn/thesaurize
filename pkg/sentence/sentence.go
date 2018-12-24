package sentence

import (
	"log"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/MrFlynn/thesaurize-bot/pkg/thesaurus"
)

// Initialize global randomizer
var randomizer = rand.New(rand.NewSource(time.Now().Unix()))

// Regexes for recognizing english words and punctutation.
var wordRegex = regexp.MustCompile("[^a-zA-Z0-9]")
var puncRegex = regexp.MustCompile("[;.,?!]")
var capitalRegex = regexp.MustCompile("[A-Z]")

// ThesaurizeSentence takes a sentence of words and replaces each with a related word.
func ThesaurizeSentence(sentence string, api *thesaurus.API) string {
	sentenceArray := strings.Split(sentence, " ")
	output := make([]string, 0, len(sentenceArray))

	for i := range sentenceArray {
		word, err := chooseWord(sentenceArray[i], api)
		if err != nil {
			log.Print(err, ": "+sentenceArray[i])

			if strings.Contains(err.Error(), "Usage exceeded") {
				return ":x: API Usage Exceeded! :x:"
			}
		}

		output = append(output, word)
	}

	return strings.Join(output, " ")
}

func chooseWord(word string, api *thesaurus.API) (string, error) {
	// Ignore word outright if it is in the ignore words list.
	if v, ok := ignoreWords[word]; ok {
		if v {
			return word, nil
		}
	}

	// Strip punctuation from word.
	strippedWord := wordRegex.ReplaceAllString(word, "")

	resp, err := api.Lookup(strippedWord)
	if err != nil {
		return word, err
	}

	similarWordList := compileWordBucket(resp)

	// Choose word using randomizer.
	randIndex := randomizer.Intn(len(similarWordList))
	outputWord := similarWordList[randIndex]

	// Handle punctuation and capital letters.
	outputWord = handlePunctuation(outputWord, word)
	outputWord = handleCapitalization(outputWord, strippedWord)

	return outputWord, nil
}

func handlePunctuation(newWord string, originalWord string) string {
	matchIndices := puncRegex.FindAllStringSubmatchIndex(originalWord, -1)

	for i := range matchIndices {
		// Calculate the offset required to insert the original punctuation
		// into the new word.
		offset := len(originalWord) - len(newWord)

		// Start and end indices of punctuation.
		start := matchIndices[i][0] - offset
		end := matchIndices[i][1] - offset

		// Get punctuation mark and insert into new string.
		punctuationMark := string(originalWord[matchIndices[i][0]])
		newWord = newWord[:start+1] + punctuationMark + newWord[end:]
	}

	return newWord
}

func handleCapitalization(newWord string, originalWord string) string {
	res := strings.Join(capitalRegex.FindAllString(originalWord, -1), "")

	if res == string(originalWord[0]) {
		return strings.Title(newWord)
	} else if res == originalWord {
		return strings.ToUpper(newWord)
	} else {
		return newWord
	}
}

// Compile list of synonyms, related words, etc. that will be used to randomly
// draw from later.
func compileWordBucket(resp thesaurus.Response) []string {
	words := make([]string, 0, 10)

	v := reflect.ValueOf(resp)
	for i := 0; i < v.NumField(); i++ {
		element := v.Field(i)

		if element.CanInterface() {
			entity := element.Interface()

			// Type cast element to WordClass.
			w, ok := entity.(thesaurus.WordClass)
			if ok {
				words = decideListAppend(w, words)
			}
		}
	}

	return words
}

// Decides which related words to add to bucket to be picked from later.
func decideListAppend(class thesaurus.WordClass, list []string) []string {
	if class.Synonym != nil {
		list = append(list, class.Synonym...)
		return list
	} else if class.Similar != nil && len(list) == 0 {
		list = append(list, class.Similar...)
		return list
	} else if class.Related != nil && len(list) == 0 {
		list = append(list, class.Related...)
		return list
	} else {
		return list
	}
}
