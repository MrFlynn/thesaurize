package database

// Lexeme type defines the part of speech associated with the word lookup.
type Lexeme int

const (
	// Noun defines the noun lexeme.
	Noun Lexeme = iota
	// Verb defines the noun lexeme.
	Verb
	// Adjective defines the noun lexeme.
	Adjective
	// Adverb defines the noun lexeme.
	Adverb
)

var (
	databaseStringMap = map[Lexeme]string{
		Noun:      "noun",
		Verb:      "verb",
		Adjective: "adj",
		Adverb:    "adv",
	}
	ordering = []Lexeme{Noun, Verb, Adjective, Adverb}
)

func (l Lexeme) String() string {
	if l < Noun || l > Adverb {
		return ""
	}

	return databaseStringMap[l]
}
