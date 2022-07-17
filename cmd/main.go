package main

import (
	telegramApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/reinanhs/golang-telegram-bot-example/pkg/handler"
	"log"
	"os"
)

var telegramTokenEnv = ""

func main() {

	telegramTokenEnv = os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramTokenEnv == "" {
		panic("Environment `TELEGRAM_BOT_TOKEN` variable not found")
	}

	bot, err := telegramApi.NewBotAPI(telegramTokenEnv)
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegramApi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			handler.ActionHandler(&update, bot)
		}
	}
}
