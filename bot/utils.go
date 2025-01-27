package bot

import (
	"log"
	"slices"

	"github.com/arinji2/law-bot/pb"
	"github.com/bwmarrin/discordgo"
)

func checkPermissions(s *discordgo.Session, i *discordgo.InteractionCreate) {
	hasPermission := false
	for _, role := range i.Member.Roles {
		if slices.Contains(AllowedRoles, role) {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have permission to use this command.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

func checkChannel(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	hasPermission := slices.Contains(AllowedChannels, i.ChannelID)
	if !hasPermission {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You cannot use this command in this channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return false
	}
	return true
}

func refreshData() {
	log.Println("Refreshing data...")

	locArticleData, err := PbAdmin.GetAllArticles()
	if err != nil {
		log.Panicf("Cannot get articles: %v", err)
		locArticleData = make([]pb.BaseCollection, 0)
	}
	log.Printf("Found %d articles", len(locArticleData))

	locClauseData, err := PbAdmin.GetAllClauses(true)
	if err != nil {
		log.Panicf("Cannot get clauses: %v", err)
		locClauseData = make([]pb.ClauseCollection, 0)
	}

	log.Printf("Found %d clauses", len(locClauseData))

	locAmendmentData, err := PbAdmin.GetAllAmendments(true)
	if err != nil {
		log.Panicf("Cannot get amendments: %v", err)
		locAmendmentData = make([]pb.AmendmentCollection, 0)
	}

	log.Printf("Found %d ammendments", len(locAmendmentData))

	ClauseCommand.ArticleData = locArticleData
	ClauseCommand.ClauseData = locClauseData
	ClauseCommand.PbAdmin = *PbAdmin

	ArticleCommand.ArticleData = locArticleData
	ArticleCommand.PbAdmin = *PbAdmin

	AmendmentCommand.ClauseData = locClauseData
	AmendmentCommand.ArticleData = locArticleData
	AmendmentCommand.AmendmentData = locAmendmentData
	AmendmentCommand.PbAdmin = *PbAdmin
}

func (b *Bot) registerCommands() []*discordgo.ApplicationCommand {
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	err := b.Session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")

	createdCommands, err := b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, b.GuildID, commands)
	if err != nil {
		log.Panicf("Cannot create commands: %v", err)
	}
	return createdCommands
}

func (b *Bot) unregisterCommands() {
	log.Println("Removing commands...")

	for _, v := range b.Commands {
		err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, b.GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
