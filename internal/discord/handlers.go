package discord

import (
	"fmt"
	"strings"

	"github.com/MrFlynn/thesaurize-bot/internal/database"
	"github.com/MrFlynn/thesaurize-bot/internal/transformer"
	"github.com/bwmarrin/discordgo"
)

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate, d database.Database) error {
	var err error

	if strings.HasPrefix(m.Content, "!thesaurize") {
		outgoingMessage := discordgo.MessageSend{}
		content := trimCommand(s, m.Message)

		if content == "help" {
			// Display help dialog.
			outgoingMessage.Embed = &helpEmbed
		} else if len(m.Mentions) > 0 && !m.MentionEveryone {
			user := m.Mentions[0]
			content, err = mentionParser(s, user, m.ChannelID)

			if err != nil {
				goto terminate
			}
		} else if len(content) == 0 {
			goto terminate
		}

		if len(content) > 0 {
			outgoingMessage.Content = transformer.Transform(content, d)
			_, err = s.ChannelMessageSendComplex(m.ChannelID, &outgoingMessage)
		}
	}

terminate:
	return err
}

func mentionParser(s *discordgo.Session, u *discordgo.User, channelID string) (string, error) {
	messages, err := s.ChannelMessages(channelID, 20, "", "", "")
	if err != nil {
		return "", fmt.Errorf(
			"Could not accesses messages in channel:%s, err:%s",
			channelID,
			err,
		)
	}

	for _, msg := range messages {
		if msg.Author.ID == u.ID {
			return trimCommand(s, msg), nil
		}
	}

	return "", nil
}
