package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MrFlynn/thesaurize/internal/discord"
	"github.com/MrFlynn/thesaurize/internal/loader"
	"github.com/urfave/cli/v2"
)

var infoTempl = `--- Thesaurize ---
Author:   %s
Compiled: %s
Commit:   %s

Report issues to https://github.com/MrFlynn/thesaurize
`

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
				Action: func(c *cli.Context) error {
					var skip bool
					if skipCommonWords == "true" {
						skip = true
					}

					fmt.Printf("Thesaurize v%s (%s)\n\n", c.App.Version, c.App.Metadata["commit"])

					return discord.Run(c, skip)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "token",
						Aliases:  []string{"t"},
						Usage:    "Discord API key",
						EnvVars:  []string{"THESAURIZE_DISCORD_TOKEN"},
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
				Name:        "load",
				Usage:       "Load data into Redis backend",
				Description: "Load data from an OpenOffice thesaurus file into Redis",
				Action: func(ctx *cli.Context) error {
					return loader.Load(ctx)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "data",
						Aliases:  []string{"d"},
						Usage:    "OpenOffice thesaurus data file",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "datastore",
						Aliases:  []string{"s"},
						Usage:    "URI of Redis datastore. Formatted like redis://<address>:<port>",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "skip-profane-words",
						Aliases: []string{"p"},
						Value:   true,
						Usage:   "Skip profane words when loading data into redis",
					},
					&cli.StringSliceFlag{
						Name:  "profane-word-categories",
						Usage: "Categories of profane words to skip (general, lgbtq, racial, religious, sexual, and/or shock)",
						Value: cli.NewStringSlice("lgbtq", "racial", "religious"),
					},
					&cli.StringFlag{
						Name:  "profane-word-index-url",
						Usage: "Index of profane words",
						Value: "https://raw.githubusercontent.com/dsojevic/profanity-list/refs/heads/main/en.json",
					},
				},
			},
			{
				Name:        "info",
				Usage:       "Get more detailed information about the bot",
				Description: "Get build information about the bot and a link to report issues",
				Action: func(c *cli.Context) error {
					fmt.Printf(
						infoTempl,
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
