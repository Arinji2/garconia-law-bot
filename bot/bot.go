package bot

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arinji2/law-bot/pb"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session  *discordgo.Session
	GuildID  string
	Commands []*discordgo.ApplicationCommand
}

var (
	pbAdmin     *pb.PocketbaseAdmin
	articleData []pb.BaseCollection
	clauseData  []pb.ClauseCollection
)

func NewBot(token string, guildID string) (*Bot, error) {
	var err error
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid token: %v", err)
	}
	return &Bot{Session: s, GuildID: guildID}, nil
}

func (b *Bot) Run(locPbAdmin *pb.PocketbaseAdmin) {
	log.Println("Starting bot...")
	pbAdmin = locPbAdmin
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	createdCommands := b.registerCommands()
	b.Commands = createdCommands
	log.Println("Bot is now running.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	_, locArticleData, err := pbAdmin.GetAllArticles()
	if err != nil {
		log.Panicf("Cannot get articles: %v", err)
		articleData = make([]pb.BaseCollection, 0)
	} else {
		articleData = locArticleData
	}

	locClauseData, err := pbAdmin.GetAllClauses(true)
	if err != nil {
		log.Panicf("Cannot get clauses: %v", err)
		clauseData = make([]pb.ClauseCollection, 0)
	} else {
		clauseData = locClauseData
	}

	<-stop
	log.Println("\nShutting down gracefully...")

	if err := b.Session.Close(); err != nil {
		log.Printf("Error closing Discord session: %v", err)
	} else {
		log.Println("Discord session closed successfully.")
	}
	b.unregisterCommands()
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
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"get-clauses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				HandleGetClauses(s, i)
			case discordgo.InteractionApplicationCommandAutocomplete:
				handleClauseAutocomplete(s, i)
			}
		},
	}
)
