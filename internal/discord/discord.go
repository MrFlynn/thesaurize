package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/MrFlynn/thesaurize-bot/internal/database"
	"github.com/MrFlynn/thesaurize-bot/internal/transformer"
	"github.com/bwmarrin/discordgo"
	"github.com/urfave/cli/v2"
)

// bot type provides methods for communicating with discord.
type bot struct {
	key            string
	database       database.Database
	serviceHandler *discordgo.Session
}

func new(ctx *cli.Context) (bot, error) {
	service, err := discordgo.New("Bot " + ctx.String("token"))
	if err != nil {
		log.Println("Could not initialize bot subsystem")
		return bot{}, err
	}

	return bot{
		key:            ctx.String("token"),
		database:       database.New(ctx.String("datastore")),
		serviceHandler: service,
	}, nil
}

func (b *bot) registerHandlers() {
	b.serviceHandler.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, "!thesaurize") {
			message := discordgo.MessageSend{}

			if len(m.Content) == 11 || m.Content == "!thesaurize help" {
				// Display help dialog.
				message.Embed = &helpEmbed
			} else {
				content := m.ContentWithMentionsReplaced()

				message.Content = transformer.Transform(
					content[12:len(content)],
					b.database,
				)
			}

			s.ChannelMessageSendComplex(m.ChannelID, &message)
		}
	})
}

func (b *bot) run() error {
	err := b.serviceHandler.Open()
	if err != nil {
		log.Println("Could not open connection to discord. Exiting...")
		return err
	}

	defer b.serviceHandler.Close()

	log.Print("Bot connected to discord")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c

	fmt.Printf("\n")
	log.Println("Bot shutting down. Goodbye...")

	return nil
}

// Run creates a bot and runs it. This provides the primary entrypoint into the bot. This function
// is called directly by the main function in the main package.
func Run(ctx *cli.Context) error {
	bot, err := new(ctx)

	if err != nil {
		log.Print("Could not initialize bot")
		return err
	}

	bot.registerHandlers()

	return bot.run()
}
