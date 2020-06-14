package transformer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBasicNewWordMetadata(t *testing.T) {
	meta := newWordMetadata("hello")
	expected := &wordMetadata{
		Word:           "hello",
		HasPunctuation: false,
		Punctuation:    "",
		Capitalization: 0,
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}
}

func TestCapitalizationNewWordMetadata(t *testing.T) {
	meta := newWordMetadata("Hello")
	expected := &wordMetadata{
		Word:           "hello",
		HasPunctuation: false,
		Punctuation:    "",
		Capitalization: 1,
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}
}

func TestAllCapitalizedNewWordMetadata(t *testing.T) {
	meta := newWordMetadata("HELLO")
	expected := &wordMetadata{
		Word:           "hello",
		HasPunctuation: false,
		Punctuation:    "",
		Capitalization: 2,
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}
}

func TestQuoteNewWordMetadata(t *testing.T) {
	meta := newWordMetadata("\"hello\"")
	expected := &wordMetadata{
		Word:           "\"hello\"",
		HasPunctuation: false,
		Punctuation:    "",
		Capitalization: 0,
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}
}

func TestPunctuationNewWordMetadata(t *testing.T) {
	meta := newWordMetadata("hello,.!?---!")
	expected := &wordMetadata{
		Word:           "hello",
		HasPunctuation: true,
		Punctuation:    ",.!?---!",
		Capitalization: 0,
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}
}

func TestBasicGenerateMetadataFromSentence(t *testing.T) {
	sentenceMeta := generateMetadataFromSentence("Hello, world! How are you doing?")
	expected := []*wordMetadata{
		&wordMetadata{
			Word:           "hello",
			HasPunctuation: true,
			Punctuation:    ",",
			Capitalization: 1,
		},
		&wordMetadata{
			Word:           "world",
			HasPunctuation: true,
			Punctuation:    "!",
			Capitalization: 0,
		},
		&wordMetadata{
			Word:           "how",
			HasPunctuation: false,
			Punctuation:    "",
			Capitalization: 1,
		},
		&wordMetadata{
			Word:           "are",
			HasPunctuation: false,
			Punctuation:    "",
			Capitalization: 0,
		},
		&wordMetadata{
			Word:           "you",
			HasPunctuation: false,
			Punctuation:    "",
			Capitalization: 0,
		},
		&wordMetadata{
			Word:           "doing",
			HasPunctuation: true,
			Punctuation:    "?",
			Capitalization: 0,
		},
	}

	if !cmp.Equal(sentenceMeta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, sentenceMeta)
	}
}
