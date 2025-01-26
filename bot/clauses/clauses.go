package clauses

import (
	"fmt"
	"log"
	"strings"

	commands_utils "github.com/arinji2/law-bot/commands"
	"github.com/arinji2/law-bot/pb"
	"github.com/bwmarrin/discordgo"
)

type ClauseCommand struct {
	ClauseData  []pb.ClauseCollection
	ArticleData []pb.BaseCollection
	PbAdmin     pb.PocketbaseAdmin
}

func (c *ClauseCommand) HandleClauseResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	switch len(data.Options) {
	case 1:
		c.handleArticleClauses(s, i, data)
	case 2:
		c.handleSpecificClause(s, i, data)
	}
}

func (c *ClauseCommand) HandleClauseAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice

	localClauseData := make([]pb.ClauseCollection, 0)
	if data.Options[0].StringValue() != "" {
		for _, v := range c.ClauseData {
			if v.Expand.Article.Number == data.Options[0].StringValue() {
				localClauseData = append(localClauseData, v)
			}
		}
	}

	switch {
	case data.Options[0].Focused:
		searchTerm := strings.ToLower(data.Options[0].StringValue())
		for i, v := range c.ArticleData {
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
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func (c *ClauseCommand) handleArticleClauses(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	articleNumber := data.Options[0].StringValue()
	description := fmt.Sprintf("Showing All Clauses For **Article Number: %s**\n", articleNumber)

	clauseData, err := c.PbAdmin.GetClausesByArticle(articleNumber)
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

func (c *ClauseCommand) handleSpecificClause(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	articleNumber := data.Options[0].StringValue()
	clauseNumber := data.Options[1].StringValue()
	description := fmt.Sprintf(
		"Showing Clause Number: **%s** For Article Number: **%s**\n\n", clauseNumber, articleNumber)
	clauseData, err := c.PbAdmin.GetClauseByNumber(clauseNumber, articleNumber, true)
	if err != nil {
		log.Printf("Error fetching specific clause: %v", err)
		commands_utils.RespondWithEphemeralError(s, i, "Could not retrieve clause data")
		return
	}

	amendmentsData, err := c.PbAdmin.GetAmendmentsByClause(clauseNumber)
	if err != nil {
		log.Printf("Error fetching amendments: %v", err)
		commands_utils.RespondWithEphemeralError(s, i, "Could not retrieve amendments data")
		return
	}

	description += fmt.Sprintf(
		"**A %s, Clause %s**: %s\n\n", clauseData.Expand.Article.Number, clauseData.Number, clauseData.Description)
	if len(amendmentsData) > 0 {
		amendmentNumbers := make([]string, 0)
		for _, v := range amendmentsData {
			amendmentNumbers = append(amendmentNumbers, v.Number)
		}
		description += fmt.Sprintf("**Amended**, See %s", strings.Join(amendmentNumbers, ", "))
	}
	commands_utils.RespondWithEmbed(s, i, "Constitution Clause Details", description)
}
