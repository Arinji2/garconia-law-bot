package main

import (
	"log"

	"github.com/arinji2/law-bot/bot"
	"github.com/arinji2/law-bot/env"
	"github.com/arinji2/law-bot/pb"
)

func main() {
	e := env.SetupEnv()
	pbAdmin := pb.SetupPocketbase(e.PB)

	discordBot, err := bot.NewBot(e.Bot.Token, e.Bot.GuildID)
	if err != nil {
		log.Panicf("Cannot create bot: %v", err)
	}
	discordBot.Run(pbAdmin)
}
