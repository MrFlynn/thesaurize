package discord

import (
	"fmt"
	"strings"

	"github.com/MrFlynn/thesaurize/internal/database"
	"github.com/MrFlynn/thesaurize/internal/transformer"
	"github.com/bwmarrin/discordgo"
)

func errorHandler(s *discordgo.Session, err error, channelID string) {
	msg := "Sorry. Something went wrong with the bot. Please try again later."
	if botErr, ok := err.(botError); ok {
		if botErr.t == errorUser {
			// Only take original error message if error is classified as a user error.
			msg = botErr.Error()
		}
	}

	s.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Title:       ":x: Thesaurize Error :x:",
		Description: msg,
		Type:        "rich",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Thesaurize Bot by MrFlynn",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Report a Bug",
				Value: "Submit an issue [here](https://github.com/MrFlynn/thesaurize/issues/new).",
			},
		},
	})
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate, d database.Database) error {
	var err error

	if strings.HasPrefix(m.Content, "!thesaurize") {
		outgoingMessage := discordgo.MessageSend{}
		content := trimCommand(s, m.Message)

		if content == "help" {
			// Display help dialog.
			outgoingMessage.Embed = &helpEmbed
		} else if len(m.Mentions) > 0 && !m.MentionEveryone {
			if m.Mentions[0].Username == content[1:] {
				content, err = mentionParser(s, m.Mentions[0], m.ChannelID)

				if err != nil {
					goto terminate
				}
			}
		} else if len(content) == 0 {
			goto terminate
		}

		if len(content) > 0 {
			outgoingMessage.Content = transformer.Transform(content, d, skipCommonWords)
			_, err = s.ChannelMessageSendComplex(m.ChannelID, &outgoingMessage)
		}
	}

terminate:
	return err
}

func mentionParser(s *discordgo.Session, u *discordgo.User, channelID string) (string, error) {
	messages, err := s.ChannelMessages(channelID, 100, "", "", "")
	if err != nil {
		return "", fmt.Errorf(
			"Could not accesses messages in channel:%s, err:%s",
			channelID,
			err,
		)
	}

	for _, msg := range messages[1:] {
		if msg.Author.ID == u.ID {
			return trimCommand(s, msg), nil
		}
	}

	return "", botError{
		why: fmt.Errorf(
			"Sorry, but I couldn't find the last message from %s. Is it be more than 100 messages ago?",
			u.Username,
		),
		t: errorUser,
	}
}
