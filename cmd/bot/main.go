package main

import (
	"os"

	"github.com/MrFlynn/thesaurize-bot/internal/discord"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "thesaurize",
		Version: Version,
		Usage:   "A Discord bot to make statements sound ridiculous.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "Discord API key.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "datastore",
				Aliases:  []string{"s"},
				Usage:    "URI of Redis datastore. Formatted like redis://<address>:<port>",
				Required: true,
			},
		},
		Action: discord.Run,
		Authors: []*cli.Author{
			{
				Name:  "Nick Pleatsikas",
				Email: "nick@pleatsikas.me",
			},
		},
	}

	app.Run(os.Args)
}
