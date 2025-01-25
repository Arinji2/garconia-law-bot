package env

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Env struct {
	Bot struct {
		Token   string
		GuildID string
	}
	Auth struct {
		Email    string
		Password string
	}
}

func SetupEnv() *Env {
	log.Println("Loading environment variables...")
	token := os.Getenv("TOKEN")
	guildID := os.Getenv("GUILD_ID")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if token == "" {
		log.Fatal("TOKEN is empty")
	}
	if guildID == "" {
		log.Fatal("GuildID is empty")
	}
	if adminEmail == "" {
		log.Fatal("ADMIN_EMAIL is empty")
	}
	if adminPassword == "" {
		log.Fatal("ADMIN_PASSWORD is empty")
	}
	log.Println("Environment variables loaded.")
	return &Env{
		Bot: struct {
			Token   string
			GuildID string
		}{
			Token:   token,
			GuildID: guildID,
		},
		Auth: struct {
			Email    string
			Password string
		}{
			Email:    adminEmail,
			Password: adminPassword,
		},
	}
}
