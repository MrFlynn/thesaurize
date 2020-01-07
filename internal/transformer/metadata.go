package transformer

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
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
	var loopContinue bool
	var b strings.Builder
	builder := &b

	for _, word := range metadataSlice {
		switch word.Capitalization {
		case firstCapitzlized:
			loopContinue = writeIfValid(builder, capitalize(word.Word))
		case allCapitalized:
			loopContinue = writeIfValid(builder, strings.ToUpper(word.Word))
		default:
			loopContinue = writeIfValid(builder, word.Word)
		}

		if word.HasPunctuation {
			loopContinue = writeIfValid(builder, word.Punctuation)
		}

		if !loopContinue {
			break
		}

		builder.WriteString(" ")
	}

	return builder.String()
}

func writeIfValid(builder *strings.Builder, s string) bool {
	if len(s)+builder.Len() >= 1997 {
		builder.WriteString("...")
		return false
	}

	builder.WriteString(s)
	return true
}

func capitalize(s string) string {
	if len(s) > 0 {
		r, sz := utf8.DecodeRuneInString(s)
		if r != utf8.RuneError || sz > 1 {
			upper := unicode.ToUpper(r)
			if upper != r {
				s = string(upper) + s[sz:]
			}
		}
	}

	return s
}
