package transformer

import "github.com/MrFlynn/thesaurize-bot/internal/database"

// Transform takes a sentence and runs each word through the thesaurus.
func Transform(sentence string, db database.Database) string {
	words := generateMetadataFromSentence(sentence)

	for _, meta := range words {
		meta.Word = db.GetBestCandidateWord(meta.Word)
	}

	return constructSentence(words)
}
