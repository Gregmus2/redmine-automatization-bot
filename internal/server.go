package internal

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"redmine-automatization-bot/internal/global"
)

var Bot *tgbotapi.BotAPI
var _testing = false

func init() {
	if _testing {
		return
	}

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	Bot, err = tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}
	Bot.Debug = false
	log.Printf("Authorized on account %s", Bot.Self.UserName)
}

func Serve() {
	defer global.Stor.Close()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			log.Printf("Non message update %s", string(update.UpdateID))
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		go handle(update.Message)
	}
}
