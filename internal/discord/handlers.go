package discord

import (
	"strings"

	"github.com/MrFlynn/thesaurize-bot/internal/database"
	"github.com/MrFlynn/thesaurize-bot/internal/transformer"
	"github.com/bwmarrin/discordgo"
)

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate, d database.Database) error {
	if strings.HasPrefix(m.Content, "!thesaurize") {
		message := discordgo.MessageSend{}

		if m.Content == "!thesaurize help" {
			// Display help dialog.
			message.Embed = &helpEmbed
		} else {
			content := m.ContentWithMentionsReplaced()

			message.Content = transformer.Transform(
				content[12:len(content)],
				d,
			)
		}

		s.ChannelMessageSendComplex(m.ChannelID, &message)
	}

	return nil
}
