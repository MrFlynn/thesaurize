package transformer

import (
	"strings"
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

func TestCapitalize(t *testing.T) {
	meta := MessageMetadata{
		Words: []string{"hello", "world", "how", "are", "you"},
		Metadata: []*WordMetadata{
			{
				Capitalization: 1,
			},
			{
				Capitalization: 0,
			},
			nil,
			nil,
			{
				Capitalization: 2,
			},
		},
	}
	expected := []string{"Hello", "world", "how", "are", "YOU"}

	for i, word := range meta.Words {
		result := meta.capitalize(word, i)

		if e := expected[i]; e != result {
			t.Errorf("Expected %s\n Got: %s\n", e, result)
		}
	}
}

func TestStringBasic(t *testing.T) {
	meta := MessageMetadata{
		Words: []string{"hello", "world"},
		Metadata: []*WordMetadata{
			nil,
			nil,
		},
	}

	if meta.String() != "hello world" {
		t.Errorf("Expected \"hello world\"\n Got %s\n", meta.String())
	}
}

func TestStringComplex(t *testing.T) {
	meta := &MessageMetadata{
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

	if meta.String() != "Hello, world! How are you doing?" {
		t.Errorf("Expected \"Hello, world! How are you doing?\"\n Got %s", meta.String())
	}
}

func TestCutoff(t *testing.T) {
	meta := &MessageMetadata{
		Words: []string{strings.Repeat("a", 1990), strings.Repeat("b", 10)},
		Metadata: []*WordMetadata{
			nil,
			nil,
		},
	}

	if meta.String() != strings.Repeat("a", 1990)+" ..." {
		t.Errorf("Expected string of length 1994\n Got length %d", len(meta.String()))
	}
}

func TestBasicIntegration(t *testing.T) {
	meta := MessageMetadata{}
	meta.New("Hey, how's it going?")

	if meta.String() != "Hey, how's it going?" {
		t.Errorf("Expected %s\n Got %s", "Hey, how's it going?", meta.String())
	}
}

func TestIntegrationWithReplacement(t *testing.T) {
	meta := MessageMetadata{}
	meta.New("Hi, my name is John.")

	meta.Words[len(meta.Words)-1] = "jane"

	if meta.String() != "Hi, my name is Jane." {
		t.Errorf("Expected %s\n Got %s", "Hi, my name is Jane.", meta.String())
	}
}
