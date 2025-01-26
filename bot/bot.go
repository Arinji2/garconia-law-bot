package bot

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arinji2/law-bot/bot/amendments"
	"github.com/arinji2/law-bot/bot/articles"
	"github.com/arinji2/law-bot/bot/clauses"
	"github.com/arinji2/law-bot/pb"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session  *discordgo.Session
	GuildID  string
	Commands []*discordgo.ApplicationCommand
}

var (
	ClauseCommand    clauses.ClauseCommand
	ArticleCommand   articles.ArticleCommand
	AmendmentCommand amendments.AmendmentCommand
)

func NewBot(token string, guildID string) (*Bot, error) {
	var err error
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid token: %v", err)
	}
	return &Bot{Session: s, GuildID: guildID}, nil
}

func (b *Bot) Run(pbAdmin *pb.PocketbaseAdmin) {
	log.Println("Starting bot...")
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	locArticleData, err := pbAdmin.GetAllArticles()
	if err != nil {
		log.Panicf("Cannot get articles: %v", err)
		locArticleData = make([]pb.BaseCollection, 0)
	}

	locClauseData, err := pbAdmin.GetAllClauses(true)
	if err != nil {
		log.Panicf("Cannot get clauses: %v", err)
		locClauseData = make([]pb.ClauseCollection, 0)
	}

	locAmendmentData, err := pbAdmin.GetAllAmendments(true)
	if err != nil {
		log.Panicf("Cannot get amendments: %v", err)
		locAmendmentData = make([]pb.AmendmentCollection, 0)
	}

	ClauseCommand.ArticleData = locArticleData
	ClauseCommand.ClauseData = locClauseData
	ClauseCommand.PbAdmin = *pbAdmin

	ArticleCommand.ArticleData = locArticleData
	ArticleCommand.PbAdmin = *pbAdmin

	AmendmentCommand.ClauseData = locClauseData
	AmendmentCommand.ArticleData = locArticleData
	AmendmentCommand.AmendmentData = locAmendmentData
	AmendmentCommand.PbAdmin = *pbAdmin

	createdCommands := b.registerCommands()
	b.Commands = createdCommands

	log.Println("Bot is now running.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("\nShutting down gracefully...")

	if err := b.Session.Close(); err != nil {
		log.Printf("Error closing Discord session: %v", err)
	} else {
		log.Println("Discord session closed successfully.")
	}

	b.unregisterCommands()
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "get-clauses",
			Description: "Get the Clauses of the Constitution",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "article-number",
					Description:  "Article Number of the Clause",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "clause-number",
					Description:  "Clause Number (Optional)",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "get-articles",
			Description: "Get the Articles of the Constitution",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "article-number",
					Description:  "Article Number of the Clause",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "get-amendments",
			Description: "Get the Amendments of the Constitution",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "article-number",
					Description:  "Article Number of the Constitution",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "clause-number",
					Description:  "Clause Number of the Article",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "amendment-number",
					Description:  "Amendment Number of the Clause",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"get-clauses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				ClauseCommand.HandleClauseResponse(s, i)
			case discordgo.InteractionApplicationCommandAutocomplete:
				ClauseCommand.HandleClauseAutocomplete(s, i)
			}
		},

		"get-articles": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				ArticleCommand.HandleArticleResponse(s, i)
			case discordgo.InteractionApplicationCommandAutocomplete:
				ArticleCommand.HandleArticleAutocomplete(s, i)
			}
		},

		"get-amendments": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				AmendmentCommand.HandleAmendmentResponse(s, i)
			case discordgo.InteractionApplicationCommandAutocomplete:
				AmendmentCommand.HandleAmendmentAutocomplete(s, i)
			}
		},
	}
)
