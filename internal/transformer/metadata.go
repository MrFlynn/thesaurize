package transformer

import (
	"regexp"
	"strings"
)

type capitalization int

const (
	noCapitalization capitalization = iota
	firstCapitzlized
	allCapitalized
)

var words = regexp.MustCompile(`\b[^\s-]+\b[^\w\s]*`)
var punctuation = regexp.MustCompile(`[.,!?-]+\B`)
var capital = regexp.MustCompile(`\b[A-Z]+`)

type wordMetadata struct {
	Word           string
	HasPunctuation bool
	Punctuation    string
	Capitalization capitalization
}

func newWordMetadata(word string) *wordMetadata {
	punctuationChars := punctuation.FindString(word)
	hasPunctuation := false

	if punctuationChars != "" {
		hasPunctuation = true
		word = punctuation.ReplaceAllLiteralString(word, "")
	}

	var capitalizationMode capitalization
	switch len(capital.FindString(word)) {
	case 0:
		capitalizationMode = noCapitalization
	case 1:
		capitalizationMode = firstCapitzlized
		word = strings.ToLower(word)
	default:
		capitalizationMode = allCapitalized
		word = strings.ToLower(word)
	}

	return &wordMetadata{
		Word:           word,
		HasPunctuation: hasPunctuation,
		Punctuation:    punctuationChars,
		Capitalization: capitalizationMode,
	}
}

func generateMetadataFromSentence(sentence string) []*wordMetadata {
	sentenceSlice := words.FindAllString(sentence, -1)
	result := make([]*wordMetadata, 0, len(sentenceSlice))

	for _, word := range sentenceSlice {
		result = append(result, newWordMetadata(word))
	}

	return result
}

func constructSentence(metadataSlice []*wordMetadata) string {
	var result strings.Builder

	for _, word := range metadataSlice {
		switch word.Capitalization {
		case firstCapitzlized:
			result.WriteString(strings.Title(word.Word))
		case allCapitalized:
			result.WriteString(strings.ToUpper(word.Word))
		default:
			result.WriteString(word.Word)
		}

		if word.HasPunctuation {
			result.WriteString(word.Punctuation)
		}

		result.WriteString(" ")
	}

	return result.String()
}
