package sentence

import (
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/MrFlynn/thesaurus-bot/pkg/thesaurus"
)

// ThesaurizeSentence takes a sentence of words and replaces each with a related word.
func ThesaurizeSentence(sentence string, api thesaurus.API) (string, error) {
	regx, _ := regexp.Compile("[^a-zA-Z0-9]")
	randomizer := rand.New(rand.NewSource(time.Now().Unix()))

	stringArray := strings.Split(sentence, " ")
	output := make([]string, 0, len(stringArray))

	for i := range stringArray {
		word := regx.ReplaceAllString(stringArray[i], "")

		if _, ok := ignoreWords[word]; ok {
			output = append(output, word)
		} else {
			resp, _ := api.Call(word)
			words := compileWordBucket(resp)

			if len(words) == 0 {
				// If thesaurus could not find synonym then return the input word.
				output = append(output, word)
			} else {
				idx := randomizer.Intn(len(words))
				output = append(output, words[idx])
			}
		}
	}

	return strings.Join(output, " "), nil
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
