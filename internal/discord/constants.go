package discord

import "github.com/bwmarrin/discordgo"

// Structs for embedding help and bot info.
var helpEmbed = discordgo.MessageEmbed{
	Title: ":book: Thesaurize Bot for Discord :book:",
	URL:   "https://github.com/MrFlynn/thesaurize",
	Description: `Makes sentences nonsensical using a thesaurus.
		Inspired by [ThesaurizeThis](https://reddit.com/r/ThesaurizeThis) by OrionSuperman.`,
	Type: "rich",
	Footer: &discordgo.MessageEmbedFooter{
		Text: "Thesaurize Bot by MrFlynn",
	},
	Fields: []*discordgo.MessageEmbedField{
		{
			Name:  "But why?",
			Value: "Because it's amusing, that's why.",
		},
		{
			Name:  "Commands",
			Value: "Just enter your text into the command `!thesaurize <text>` to use this bot",
		},
	},
}
