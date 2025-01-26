package amendments

import (
	"fmt"
	"log"
	"strings"

	commands_utils "github.com/arinji2/law-bot/commands"
	"github.com/arinji2/law-bot/pb"
	"github.com/bwmarrin/discordgo"
)

type AmendmentCommand struct {
	ClauseData    []pb.ClauseCollection
	ArticleData   []pb.BaseCollection
	AmendmentData []pb.AmendmentCollection
	PbAdmin       pb.PocketbaseAdmin
}

func (a *AmendmentCommand) HandleAmendmentResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	switch len(data.Options) {
	case 3:
		a.handleArticleClauses(s, i, data)
	}
}

func (a *AmendmentCommand) HandleAmendmentAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice

	localClauseData := make([]pb.ClauseCollection, 0)
	if data.Options[0].StringValue() != "" {
		for _, v := range a.ClauseData {
			if v.Expand.Article.Number == data.Options[0].StringValue() {
				localClauseData = append(localClauseData, v)
			}
		}
	}
	localAmendmentData := make([]pb.AmendmentCollection, 0)
	if len(data.Options) > 1 {
		if data.Options[1].StringValue() != "" {
			for _, v := range a.AmendmentData {
				if v.Expand.Clause.Number == data.Options[1].StringValue() {
					localAmendmentData = append(localAmendmentData, v)
				}
			}
		}
	}

	switch {
	case data.Options[0].Focused:
		searchTerm := strings.ToLower(data.Options[0].StringValue())
		for i, v := range a.ArticleData {
			if i > 25 {
				break
			}
			if strings.Contains(strings.ToLower(v.Description), searchTerm) ||
				strings.Contains(strings.ToLower(v.Number), searchTerm) {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  commands_utils.FormatDescription(fmt.Sprintf("Article %s: %s", v.Number, v.Description)),
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
					Name:  commands_utils.FormatDescription(fmt.Sprintf("A %s, Clause %s: %s", v.Expand.Article.Number, v.Number, v.Description)),
					Value: v.Number,
				})
			}
		}

	case data.Options[2].Focused:
		searchTerm := strings.ToLower(data.Options[0].StringValue())
		for i, v := range localAmendmentData {
			if i > 25 {
				break
			}
			if strings.Contains(strings.ToLower(v.Description), searchTerm) ||
				strings.Contains(strings.ToLower(v.Number), searchTerm) {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  commands_utils.FormatDescription(fmt.Sprintf("%s Amendment: %s", commands_utils.OrdinalRepresentation(v.Number), v.Description)),
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

func (a *AmendmentCommand) handleArticleClauses(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	articleNumber := data.Options[0].StringValue()
	clauseNumber := data.Options[1].StringValue()
	amendmentNumber := data.Options[2].StringValue()
	description := fmt.Sprintf("Showing **%s Amendment** For **Clause %s** Of **Article %s**\n\n", commands_utils.OrdinalRepresentation(amendmentNumber), clauseNumber, articleNumber)

	amendmentData, err := a.PbAdmin.GetAmendmentByNumber(amendmentNumber, clauseNumber, articleNumber, true)
	if err != nil {
		log.Printf("Error fetching amendment: %v", err)
		commands_utils.RespondWithEphemeralError(s, i, "Could not find amendment data")
		return
	}

	description += fmt.Sprintf("**%s Amendement**: %s\n\n", commands_utils.OrdinalRepresentation(amendmentData.Number), amendmentData.Description)

	description += fmt.Sprintf("**Ammends**, See Clause %s", amendmentData.Expand.Clause.Number)
	commands_utils.RespondWithEmbed(s, i, "Constitution Clause Details", description)
}
