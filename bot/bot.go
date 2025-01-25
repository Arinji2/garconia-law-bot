package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session  *discordgo.Session
	GuildID  string
	Commands []*discordgo.ApplicationCommand
}

func NewBot(token string, guildID string) (*Bot, error) {
	var err error
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid token: %v", err)
	}
	return &Bot{Session: s, GuildID: guildID}, nil
}

func (b *Bot) Run() {
	log.Println("Starting bot...")
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
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
	defer b.Session.Close()

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
			Name:        "multi-autocomplete",
			Description: "Showcase of multiple autocomplete option",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "autocomplete-option-1",
					Description:  "Autocomplete option 1",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "autocomplete-option-2",
					Description:  "Autocomplete option 2",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"multi-autocomplete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				data := i.ApplicationCommandData()
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(
							"Option 1: %s\nOption 2: %s",
							data.Options[0].StringValue(),
							data.Options[1].StringValue(),
						),
					},
				})
				if err != nil {
					panic(err)
				}
			case discordgo.InteractionApplicationCommandAutocomplete:
				data := i.ApplicationCommandData()
				var choices []*discordgo.ApplicationCommandOptionChoice
				switch {
				case data.Options[0].Focused:
					choices = []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Autocomplete 4 first option",
							Value: "autocomplete_default",
						},
						{
							Name:  "Choice 3",
							Value: "choice3",
						},
						{
							Name:  "Choice 4",
							Value: "choice4",
						},
						{
							Name:  "Choice 5",
							Value: "choice5",
						},
					}
					if data.Options[0].StringValue() != "" {
						choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
							Name:  data.Options[0].StringValue(),
							Value: "choice_custom",
						})
					}

				case data.Options[1].Focused:
					choices = []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Autocomplete 4 second option",
							Value: "autocomplete_1_default",
						},
						{
							Name:  "Choice 3.1",
							Value: "choice3_1",
						},
						{
							Name:  "Choice 4.1",
							Value: "choice4_1",
						},
						{
							Name:  "Choice 5.1",
							Value: "choice5_1",
						},
					}
					if data.Options[1].StringValue() != "" {
						choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
							Name:  data.Options[1].StringValue(),
							Value: "choice_custom_2",
						})
					}
				}

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionApplicationCommandAutocompleteResult,
					Data: &discordgo.InteractionResponseData{
						Choices: choices,
					},
				})
				if err != nil {
					panic(err)
				}
			}
		},
	}
)
