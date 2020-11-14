package transformer

import "github.com/MrFlynn/thesaurize/internal/database"

// Transform takes a message and runs each word through the thesaurus.
func Transform(message string, db database.Database, skipCommon bool) string {
	messageMeta := MessageMetadata{}
	messageMeta.New(message)

	for idx, word := range messageMeta.Words {
		// Skip word if it's in a preconfigured list of words to ignore.
		if skipCommon {
			if _, ok := ignoreWords[word]; ok {
				continue
			}
		}

		messageMeta.Words[idx] = db.GetBestCandidateWord(word)
	}

	return messageMeta.String()
}
