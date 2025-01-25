package commands_utils

import "github.com/bwmarrin/discordgo"

const (
	botFooter  = "Garconian Constitution Bot"
	embedColor = 0xA31621 // Garconia Red
)

func CreateBaseEmbed(title, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       embedColor,
		Footer: &discordgo.MessageEmbedFooter{
			Text: botFooter,
		},
	}
}

func RespondWithEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, title, description string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{CreateBaseEmbed(title, description)},
		},
	})
}

func RespondWithEphemeralError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
