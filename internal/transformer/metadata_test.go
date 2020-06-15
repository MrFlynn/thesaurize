package transformer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBasicCreateWordMetadata(t *testing.T) {
	meta, word := createWordMetadata("hello")

	if meta != nil {
		t.Errorf("Expected: nil, Got: %+v\n", meta)
	}

	if word != "hello" {
		t.Errorf("Expected \"hello\", got: %s", word)
	}
}

func TestCapitalizationCreateWordMetadata(t *testing.T) {
	meta, word := createWordMetadata("Hello")
	expected := &WordMetadata{
		Capitalization: 1,
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}

	if word != "hello" {
		t.Errorf("Expected \"hello\", got: %s", word)
	}
}

func TestAllCapitalizedCreateWordMetadata(t *testing.T) {
	meta, word := createWordMetadata("HELLO")
	expected := &WordMetadata{
		Capitalization: 2,
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}

	if word != "hello" {
		t.Errorf("Expected \"hello\", got: %s", word)
	}
}

func TestQuoteNewWordMetadata(t *testing.T) {
	meta, word := createWordMetadata("\"hello\"")
	expected := &WordMetadata{
		Capitalization: 0,
		PrePunc:        `"`,
		PostPunc:       `"`,
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}

	if word != "hello" {
		t.Errorf("Expected \"\"hello\"\", got: %s", word)
	}
}

func TestPrePuncCreateWordMetadata(t *testing.T) {
	meta, word := createWordMetadata("@hello")
	expected := &WordMetadata{
		Capitalization: 0,
		PrePunc:        "@",
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}

	if word != "hello" {
		t.Errorf("Expected \"hello\", got: %s", word)
	}
}

func TestPostPuncCreateWordMetadata(t *testing.T) {
	meta, word := createWordMetadata("hello!!!")
	expected := &WordMetadata{
		Capitalization: 0,
		PostPunc:       "!!!",
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}

	if word != "hello" {
		t.Errorf("Expected \"hello\", got: %s", word)
	}
}

func TestPunctuationNewWordMetadata(t *testing.T) {
	meta, word := createWordMetadata("hello,.!?---!")
	expected := &WordMetadata{
		Capitalization: 0,
		PostPunc:       ",.!?---!",
	}

	if !cmp.Equal(meta, expected) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected, meta)
	}

	if word != "hello" {
		t.Errorf("Expected \"hello\", got: %s", word)
	}
}

func TestBasicGenerateMetadataFromSentence(t *testing.T) {
	meta := MessageMetadata{}
	meta.New("Hello, world! How are you doing?")

	expected := &MessageMetadata{
		Words: []string{"hello", "world", "how", "are", "you", "doing"},
		Metadata: []*WordMetadata{
			{
				Capitalization: 1,
				PostPunc:       ",",
			},
			{
				Capitalization: 0,
				PostPunc:       "!",
			},
			{
				Capitalization: 1,
			},
			nil,
			nil,
			{
				Capitalization: 0,
				PostPunc:       "?",
			},
		},
	}

	if !cmp.Equal(meta.Words, expected.Words) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected.Words, meta.Words)
	}

	if !cmp.Equal(meta.Metadata, expected.Metadata) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected.Metadata, meta.Words)
	}
}

func TestComplexGenerateMetadataFromSentence(t *testing.T) {
	meta := MessageMetadata{}
	meta.New("Sentence with hyphenated-words and middle$#!punctuation.")

	expected := &MessageMetadata{
		Words: []string{"sentence", "with", "hyphenated", "words", "and", "middle$#!punctuation"},
		Metadata: []*WordMetadata{
			{
				Capitalization: 1,
			},
			nil,
			{
				Capitalization: 0,
				PostPunc:       "-",
			},
			nil,
			nil,
			{
				Capitalization: 0,
				PostPunc:       ".",
			},
		},
	}

	if !cmp.Equal(meta.Words, expected.Words) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected.Words, meta.Words)
	}

	if !cmp.Equal(meta.Metadata, expected.Metadata) {
		t.Errorf("Expected: %+v\n Got: %+v\n", expected.Metadata, meta.Words)
	}
}
