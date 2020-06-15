package transformer

import "github.com/MrFlynn/thesaurize-bot/internal/database"

// Transform takes a message and runs each word through the thesaurus.
func Transform(message string, db database.Database) string {
	messageMeta := MessageMetadata{}
	messageMeta.New(message)

	for idx, word := range messageMeta.Words {
		messageMeta.Words[idx] = db.GetBestCandidateWord(word)
	}

	return messageMeta.String()
}
