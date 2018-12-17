package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/MrFlynn/thesaurus-bot/pkg/sentence"
	"github.com/MrFlynn/thesaurus-bot/pkg/thesaurus"

	"github.com/bwmarrin/discordgo"
)

// InitializeBot initializes the discord bot and creates channel to listen for commands
// and OS events.
func InitializeBot(botToken string, thesaurusAPI thesaurus.API) {
	// Initialize the bot.
	dg, err := discordgo.New("Bot " + botToken)

	// Add a handler for listening for "!thesaurize" command.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, "!thesaurize") {
			newSentence, err := sentence.ThesaurizeSentence(m.Content[11:len(m.Content)], thesaurusAPI)
			if err != nil {
				log.Fatal(err)
			}

			s.ChannelMessageSend(m.ChannelID, newSentence)
		}
	})

	// Open the websocket connection. Handle errors here.
	err = dg.Open()
	if err != nil {
		fmt.Println("Could not open websocket connection. Exiting...", err)
		return
	}

	defer dg.Close()

	// Make the channel and set the channel notifications.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
