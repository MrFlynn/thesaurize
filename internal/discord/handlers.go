package discord

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/MrFlynn/thesaurize/internal/transformer"
	"github.com/bwmarrin/discordgo"
)

var commandFormat = regexp.MustCompile(`^</\w+:\d+>`)

func errorHandler(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	msg := "Sorry. Something went wrong with the bot. Please try again later."
	if botErr, ok := err.(botError); ok {
		if botErr.t == errorUser {
			// Only take original error message if error is classified as a user error.
			msg = botErr.Error()
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
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
				},
			},
		},
	})
}

func (b *bot) commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Data.Name == "thesaurize" {
		if len(i.Data.Options) < 1 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Embeds: []*discordgo.MessageEmbed{
						helpEmbed,
					},
				},
			})

			return
		}

		switch i.Data.Options[0].Name {
		case "words":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: transformer.Transform(i.Data.Options[0].StringValue(), b.database, skipCommonWords),
				},
			})
		case "member":
			message, err := mentionParser(s, i.Data.Options[0].UserValue(s), i.ChannelID)
			if err != nil {
				errorHandler(s, i, err)

				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: transformer.Transform(message, b.database, skipCommonWords),
				},
			})
		default:
			errorHandler(s, i, botError{
				why: errors.New("Unknown option. Please provide some text or a username to the bot"),
				t:   errorUser,
			})
		}
	}
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

	for _, msg := range messages {
		if msg.Author.ID == u.ID && !commandFormat.MatchString(msg.Content) {
			return msg.Content, nil
		}
	}

	return "", botError{
		why: fmt.Errorf(
			"Sorry, but I couldn't find the last message from %s. Is it be more than 100 messages ago?",
			u.Mention(),
		),
		t: errorUser,
	}
}

func (b *bot) joinHandler(s *discordgo.Session, c *discordgo.GuildCreate) {
	if c.Guild.Unavailable {
		log.Printf("guild %s unavailable", c.Guild.Name)

		return
	}

	joined, err := b.database.IsServerJoined(c.Guild.ID)
	if err != nil {
		log.Printf("Could not check if server was in joined set %s", err)
	}

	if joined {
		return
	}

	for _, channel := range c.Guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
			_, err := s.ChannelMessageSendEmbed(channel.ID, helpEmbed)
			if err == nil {
				b.database.AddJoinedServer(c.Guild.ID)
				return // If channel delivery was successful, exit.
			}
		}
	}

	log.Printf("Could not find default channel in %s to send help message", c.Guild.Name)
}
