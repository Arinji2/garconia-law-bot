package articles

import (
	"fmt"
	"log"
	"strings"

	commands_utils "github.com/arinji2/law-bot/commands"
	"github.com/arinji2/law-bot/pb"
	"github.com/bwmarrin/discordgo"
)

type ArticleCommand struct {
	ArticleData []pb.BaseCollection
	PbAdmin     pb.PocketbaseAdmin
}

func (a *ArticleCommand) HandleArticleResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	switch len(data.Options) {
	case 0:
		a.handleAllArticles(s, i)
	case 1:
		a.handleSpecificArticle(s, i, data)
	}
}

func (a *ArticleCommand) HandleArticleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice

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
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func (a *ArticleCommand) handleAllArticles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	description := "Showing **All Articles** \n"

	if len(a.ArticleData) == 0 {
		commands_utils.RespondWithEphemeralError(s, i, "No articles found ")
		return
	}

	for _, v := range a.ArticleData {
		description += fmt.Sprintf(
			"**Article %s**: %s\n\n", v.Number, v.Description)
	}

	commands_utils.RespondWithEmbed(s, i, "Constitution Article Details", description)
}

func (a *ArticleCommand) handleSpecificArticle(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	articleNumber := data.Options[0].StringValue()
	description := fmt.Sprintf(
		"Showing **Article Number: %s**\n", articleNumber)

	articleData, err := a.PbAdmin.GetArticleByNumber(articleNumber)
	if err != nil {
		log.Printf("Error fetching specific article: %v", err)
		commands_utils.RespondWithEphemeralError(s, i, "Could not retrieve article data")
		return
	}

	clauseData, err := a.PbAdmin.GetClausesByArticle(articleNumber)
	if err != nil {
		log.Printf("Error fetching specific article: %v", err)
		clauseData = make([]pb.ClauseCollection, 0)
	}
	description += fmt.Sprintf(
		"**Article %s**: %s\n\n **Total Clauses**: %d \n", articleData.Number, articleData.Description, len(clauseData))
	if len(clauseData) > 0 {
		description += fmt.Sprintf("Showing All Clauses For **Article Number: %s**\n\n", articleNumber)
		for _, v := range clauseData {
			description += fmt.Sprintf(
				"**A %s, Clause %s**: %s\n\n", v.Expand.Article.Number, v.Number, v.Description)
		}
	}
	err = commands_utils.RespondWithEmbed(s, i, "Constitution Article Details", description)
	if err != nil {
		log.Printf("Error sending article response: %v", err)
	}
}
