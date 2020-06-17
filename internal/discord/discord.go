package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"

	"github.com/MrFlynn/thesaurize-bot/internal/database"
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

func (b *bot) registerHandler(handler func(s *discordgo.Session, m *discordgo.MessageCreate, d database.Database) error) {
	b.serviceHandler.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		err := handler(s, m, b.database)
		if err != nil {
			name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			log.Printf("Got error: %s from handler %s", err, name)
		}
	})
}

func (b *bot) run() error {
	var err error

	err = b.database.WaitForReady(30)
	if err != nil {
		log.Println(err)
		return err
	}

	err = b.serviceHandler.Open()
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

	return err
}

func trimCommand(s *discordgo.Session, m *discordgo.Message) string {
	msg, _ := m.ContentWithMoreMentionsReplaced(s)
	return strings.TrimPrefix(msg, "!thesaurize ")
}

// Run creates a bot and runs it. This provides the primary entrypoint into the bot. This function
// is called directly by the main function in the main package.
func Run(ctx *cli.Context) error {
	bot, err := new(ctx)

	if err != nil {
		log.Print("Could not initialize bot")
		return err
	}

	bot.registerHandler(commandHandler)

	return bot.run()
}
