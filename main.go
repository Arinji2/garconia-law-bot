package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

var s *discordgo.Session

func main() {
	token := os.Getenv("TOKEN")
	guildID := os.Getenv("GUILD_ID")

	if token == "" {
		log.Fatal("TOKEN is empty")
	}

	if guildID == "" {
		log.Fatal("GuildID is empty")
	}
	var err error
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid token: %v", err)
	}
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	log.Println("Adding commands...")

	createdCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildID, commands)
	if err != nil {
		log.Panicf("Cannot create commands: %v", err)
	}
	fmt.Println("Bot is now running.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("\nShutting down gracefully...")

	if err := s.Close(); err != nil {
		log.Printf("Error closing Discord session: %v", err)
	} else {
		log.Println("Discord session closed successfully.")
	}
	log.Println("Removing commands...")

	for _, v := range createdCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "single-autocomplete",
			Description: "Showcase of single autocomplete option",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "autocomplete-option",
					Description:  "Autocomplete option",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
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
		"single-autocomplete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				data := i.ApplicationCommandData()
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(
							"You picked %q autocompletion",
							data.Options[0].StringValue(),
						),
					},
				})
				if err != nil {
					panic(err)
				}
			case discordgo.InteractionApplicationCommandAutocomplete:
				data := i.ApplicationCommandData()
				choices := []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Autocomplete",
						Value: "autocomplete",
					},
					{
						Name:  "Autocomplete is best!",
						Value: "autocomplete_is_best",
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
				// In this case there are multiple autocomplete options. The Focused field shows which option user is focused on.
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
