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

// Regexes for handling how to split messages into usable components.
var wordSplitRegex = regexp.MustCompile(`\w\b-+|\S+`)
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
	size     uint32
}

// New initializes message metadata struct from a string.
func (m *MessageMetadata) New(message string) {
	wordList := wordSplitRegex.FindAllString(message, -1)

	m.size = uint32(len(message))
	m.Words = make([]string, len(wordList))
	m.Metadata = make([]*WordMetadata, len(wordList))

	for idx, word := range wordList {
		meta, normalizedWord := createWordMetadata(word)

		m.Words[idx] = normalizedWord
		m.Metadata[idx] = meta
	}
}

func (m MessageMetadata) capitalize(word string, idx int) string {
	if m.Metadata[idx] == nil {
		return word
	}

	switch m.Metadata[idx].Capitalization {
	case firstCapitzlized:
		return capitalizeFirst(word)
	case allCapitalized:
		return strings.ToUpper(word)
	default:
		return word
	}
}

func (m MessageMetadata) String() string {
	builder := ReversibleStringBuilder{}
	builder.Init()

	// Grow the buffer so that we have some headroom over the original string.
	builder.Grow(int(1.2 * float32(m.size)))

	for idx, word := range m.Words {
		if meta := m.Metadata[idx]; meta != nil {
			builder.WriteString(meta.PrePunc)
			builder.WriteString(m.capitalize(word, idx))
			builder.WriteString(meta.PostPunc)
		} else {
			builder.WriteString(word)
		}

		if builder.Len() >= 1997 {
			builder.Reverse(-1)
			builder.WriteString("....")

			break
		}

		builder.WriteString(" ")
		builder.Flush()
	}

	return builder.String()[:builder.Len()-1]
}

func capitalizeFirst(s string) string {
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
