package main

import (
	"os"

	"github.com/MrFlynn/thesaurize-bot/pkg/discord"
	"github.com/MrFlynn/thesaurize-bot/pkg/thesaurus"
)

func main() {
	thesaurusAPI := thesaurus.API{Key: os.Getenv("THESAURUS_KEY")}
	discordKey := os.Getenv("DISCORD_KEY")

	discord.InitializeBot(discordKey, &thesaurusAPI)
}
