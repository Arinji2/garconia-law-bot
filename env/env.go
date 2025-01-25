package env

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Bot struct {
	Token   string
	GuildID string
}

type PB struct {
	Email      string
	Password   string
	BaseDomain string
}
type Env struct {
	Bot Bot
	PB  PB
}

func loadEnv(envName string) string {
	val := os.Getenv(envName)
	if val == "" {
		log.Fatalf("Environment variable %s is empty", envName)
	}
	return val
}

func SetupEnv() *Env {
	log.Println("Loading environment variables...")
	token := loadEnv("TOKEN")
	guildID := loadEnv("GUILD_ID")
	adminEmail := loadEnv("ADMIN_EMAIL")
	adminPassword := loadEnv("ADMIN_PASSWORD")
	baseDomain := loadEnv("BASE_DOMAIN")
	log.Println("Environment variables loaded.")
	return &Env{
		Bot: Bot{
			Token:   token,
			GuildID: guildID,
		},
		PB: PB{
			Email:      adminEmail,
			Password:   adminPassword,
			BaseDomain: baseDomain,
		},
	}
}
