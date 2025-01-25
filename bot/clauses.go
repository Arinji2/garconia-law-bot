package bot

import (
	"fmt"
	"log"
	"strings"

	commands_utils "github.com/arinji2/law-bot/commands"
	"github.com/arinji2/law-bot/pb"
	"github.com/bwmarrin/discordgo"
)

func handleClauseResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	switch len(data.Options) {
	case 1:
		handleArticleClauses(s, i, data)
	case 2:
		handleSpecificClause(s, i, data)
	}
}

func handleClauseAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice

	localClauseData := make([]pb.ClauseCollection, 0)
	if data.Options[0].StringValue() != "" {
		for _, v := range clauseData {
			if v.Expand.Article.Number == data.Options[0].StringValue() {
				localClauseData = append(localClauseData, v)
			}
		}
	}

	switch {
	case data.Options[0].Focused:
		searchTerm := strings.ToLower(data.Options[0].StringValue())
		for i, v := range articleData {
			if i > 25 {
				break
			}
			if strings.Contains(strings.ToLower(v.Description), searchTerm) ||
				strings.Contains(strings.ToLower(v.Number), searchTerm) {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  formatDescription(fmt.Sprintf("Article %s: %s", v.Number, v.Description)),
					Value: v.Number,
				})
			}
		}

	case data.Options[1].Focused:
		searchTerm := strings.ToLower(data.Options[1].StringValue())
		for i, v := range localClauseData {
			if i > 25 {
				break
			}
			if strings.Contains(strings.ToLower(v.Description), searchTerm) ||
				strings.Contains(strings.ToLower(v.Number), searchTerm) {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  formatDescription(fmt.Sprintf("A %s, Clause %s: %s", v.Expand.Article.Number, v.Number, v.Description)),
					Value: v.Number,
				})
			}
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func handleArticleClauses(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	articleNumber := data.Options[0].StringValue()
	description := fmt.Sprintf("Showing All Clauses For **Article Number: %s**\n", articleNumber)

	clauseData, err := pbAdmin.GetClausesByArticle(articleNumber)
	if err != nil {
		log.Printf("Error fetching clauses: %v", err)
		commands_utils.RespondWithEphemeralError(s, i, "Could not retrieve clause data")
		return
	}

	if len(clauseData) == 0 {
		commands_utils.RespondWithEphemeralError(s, i, "No clauses found for this article")
		return
	}

	for _, v := range clauseData {
		description += fmt.Sprintf(
			"**A %s, Clause %s**: %s\n\n", v.Expand.Article.Number, v.Number, v.Description)
	}

	commands_utils.RespondWithEmbed(s, i, "Constitution Clause Details", description)
}

func handleSpecificClause(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	articleNumber := data.Options[0].StringValue()
	clauseNumber := data.Options[1].StringValue()
	description := fmt.Sprintf(
		"Showing **Clause Number: %s** For **Article Number: %s**\n", clauseNumber, articleNumber)

	clauseData, err := pbAdmin.GetClauseByNumber(clauseNumber, articleNumber, true)
	if err != nil {
		log.Printf("Error fetching specific clause: %v", err)
		commands_utils.RespondWithEphemeralError(s, i, "Could not retrieve clause data")
		return
	}

	description += fmt.Sprintf(
		"**A %s, Clause %s**: %s\n", clauseData.Expand.Article.Number, clauseData.Number, clauseData.Description)

	commands_utils.RespondWithEmbed(s, i, "Constitution Clause Details", description)
}
