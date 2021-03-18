package discord

import "github.com/bwmarrin/discordgo"

// Structs for embedding help and bot info.
var (
	command = &discordgo.ApplicationCommand{
		Name:        "thesaurize",
		Description: "Run some words through a thesaurus",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "words",
				Description: "Words to run through thesaurus",
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Thesaurus this member's last message",
			},
		},
	}

	helpEmbed = &discordgo.MessageEmbed{
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
				Name:  "Basic Usage",
				Value: "Just enter some text into the command `/thesaurize words:<text>` to get start.",
			},
			{
				Name:  "Thesaurizing a Previous Message",
				Value: "Use the command `/thesaurize member:@member` to thesaurize their last message.",
			},
		},
	}
)
