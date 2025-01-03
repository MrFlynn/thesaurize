package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/MrFlynn/thesaurize/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/urfave/cli/v2"
)

var (
	// Set if common words should be skipped.
	skipCommonWords bool
)

// bot type provides methods for communicating with discord.
type bot struct {
	key            string
	database       database.Database
	serviceHandler *discordgo.Session
}

// Error handling for bot. Stores information on whether to expose error
// messages to end user.
type errorType int

const (
	errorInternal errorType = iota
	errorUser
)

type botError struct {
	why error
	t   errorType
}

func (e botError) Error() string {
	return e.why.Error()
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

func (b *bot) run(ctx *cli.Context) error {
	var err error

	err = b.database.WaitForReady(ctx.Int("timeout"))
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

	_, err = b.serviceHandler.ApplicationCommandCreate(b.serviceHandler.State.User.ID, "", command)
	if err != nil {
		log.Println("Could not register application commands. Exiting...")
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c

	fmt.Printf("\n")
	log.Println("Bot shutting down. Goodbye...")

	return err
}

// Run creates a bot and runs it. This provides the primary entrypoint into the bot. This function
// is called directly by the main function in the main package. The `skip` variable instructs the
// bot to skip common words in translation.
func Run(ctx *cli.Context, skip bool) error {
	skipCommonWords = skip

	bot, err := new(ctx)
	if err != nil {
		log.Print("Could not initialize bot")
		return err
	}

	bot.serviceHandler.AddHandler(bot.commandHandler)
	bot.serviceHandler.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) {
		s.UpdateGameStatus(0, "Reading a thesaurus")
	})

	return bot.run(ctx)
}
