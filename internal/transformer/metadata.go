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
var punctuationRegex = regexp.MustCompile(`^(\W+)|(\W+)$`)
var capitalRegex = regexp.MustCompile(`\b[A-Z]+`)

// WordMetadata is individual word metadata information.
type WordMetadata struct {
	Capitalization capitalization
	PrePunc        string
	PostPunc       string
}

// This method might look a little contrived with the number of "if" statements
// checking if `meta` is nil. This is intentional. The point of this is to
// reduce heap allocations as much as possible by only creating meta as needed.
func createWordMetadata(word string) (*WordMetadata, string) {
	var meta *WordMetadata

	if idxs := punctuationRegex.FindAllStringIndex(word, 2); len(idxs) > 0 {
		var pre, post string

		if len(idxs) == 1 {
			if idxs[0][0] == 0 {
				pre = word[0:idxs[0][1]]
			} else {
				post = word[idxs[0][0]:idxs[0][1]]
			}
		} else {
			pre = word[0:idxs[0][1]]
			post = word[idxs[1][0]:idxs[1][1]]
		}

		meta = &WordMetadata{
			PrePunc:  pre,
			PostPunc: post,
		}

		word = punctuationRegex.ReplaceAllLiteralString(word, "")
	}

	switch len(capitalRegex.FindString(word)) {
	case 0:
		if meta != nil {
			meta.Capitalization = noCapitalization
		}
	case 1:
		if meta == nil {
			meta = &WordMetadata{}
		}

		meta.Capitalization = firstCapitzlized
		word = strings.ToLower(word)
	default:
		if meta == nil {
			meta = &WordMetadata{}
		}

		meta.Capitalization = allCapitalized
		word = strings.ToLower(word)
	}

	return meta, word
}

// MessageMetadata contains a list of words and their associated metadata.
type MessageMetadata struct {
	Words    []string
	Metadata []*WordMetadata
}

// New initializes message metadata struct from a string.
func (m *MessageMetadata) New(message string) {
	wordList := words.FindAllString(message, -1)

	m.Words = make([]string, len(wordList))
	m.Metadata = make([]*WordMetadata, len(wordList))

	for idx, word := range wordList {
		meta, normalizedWord := createWordMetadata(word)

		m.Words[idx] = normalizedWord
		m.Metadata[idx] = meta
	}
}

type wordMetadata struct {
	Word           string
	HasPunctuation bool
	Punctuation    string
	Capitalization capitalization
}

func newWordMetadata(word string) *wordMetadata {
	punctuationChars := punctuationRegex.FindString(word)
	hasPunctuation := false

	if punctuationChars != "" {
		hasPunctuation = true
		word = punctuationRegex.ReplaceAllLiteralString(word, "")
	}

	var capitalizationMode capitalization
	switch len(capitalRegex.FindString(word)) {
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
