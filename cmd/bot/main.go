package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MrFlynn/thesaurize-bot/internal/discord"
	"github.com/urfave/cli/v2"
)

func main() {
	compiled, err := time.Parse(time.RFC3339, date)
	if err != nil {
		compiled = time.Now()
	}

	app := &cli.App{
		// Basic information.
		Name:  "thesaurize",
		Usage: "A Discord bot to make statements sound ridiculous.",
		// Commands.
		Commands: []*cli.Command{
			{
				Name:        "run",
				Usage:       "Run the bot",
				Description: "Configure and run the discord bot",
				Action:      discord.Run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "token",
						Aliases:  []string{"t"},
						Usage:    "Discord API key",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "datastore",
						Aliases:  []string{"s"},
						Usage:    "URI of Redis datastore. Formatted like redis://<address>:<port>",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "timeout",
						Aliases: []string{"w"},
						Usage:   "How long to wait for the database in seconds. A value of 0 will skip this check",
						Value:   30,
					},
				},
			},
			{
				Name:        "info",
				Usage:       "Get more detailed information about the bot",
				Description: "Get build information about the bot and a link to report issues",
				Action: func(c *cli.Context) error {
					fmt.Printf(
						"--- Thesaurize --- \nAuthor:   %s\nCompiled: %s\nCommit:   %s\n\nReport issues to https://github.com/MrFlynn/thesaurize\n",
						c.App.Authors[0].String(),
						c.App.Compiled.String(),
						c.App.Metadata["commit"],
					)

					return nil
				},
			},
		},
		// Exit handler.
		ExitErrHandler: func(context *cli.Context, err error) {
			log.Fatalf("Application ran into fatal error: %s", err)
		},
		// App information.
		Version:  version,
		Compiled: compiled,
		Metadata: map[string]interface{}{
			"commit": commit,
		},
		Authors: []*cli.Author{
			{
				Name:  "Nick Pleatsikas",
				Email: "nick@pleatsikas.me",
			},
		},
	}

	app.Run(os.Args)
}
