package sentence

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/MrFlynn/thesaurus-bot/pkg/thesaurus"
)

// ThesaurizeSentence takes a sentence of words and replaces each with a related word.
func ThesaurizeSentence(sentence string, api thesaurus.API) (string, error) {
	regx, _ := regexp.Compile("[^a-zA-Z0-9]")
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	stringArray := strings.Split(sentence, " ")
	output := make([]string, 0, len(stringArray))

	for i := range stringArray {
		word := regx.ReplaceAllString(stringArray[i], "")
		synonymResponse, _ := api.Call(word)

		words := collectWords(synonymResponse)

		// If nothing could be appended
		if len(words) == 0 {
			output = append(output, word)
		} else {
			idx := r.Intn(len(words))
			output = append(output, words[idx])
		}
	}

	return strings.Join(output, " "), nil
}

func collectWords(resp thesaurus.Response) []string {
	words := make([]string, 0, 10)

	words = append(words, resp.Adjective.Related...)
	words = append(words, resp.Adjective.Similar...)
	words = append(words, resp.Adjective.Synonym...)

	words = append(words, resp.Adverb.Related...)
	words = append(words, resp.Adverb.Similar...)
	words = append(words, resp.Adverb.Synonym...)

	words = append(words, resp.Noun.Related...)
	words = append(words, resp.Noun.Similar...)
	words = append(words, resp.Noun.Synonym...)

	words = append(words, resp.Verb.Related...)
	words = append(words, resp.Verb.Similar...)
	words = append(words, resp.Verb.Synonym...)

	return words
}
