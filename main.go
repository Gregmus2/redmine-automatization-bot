package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"redmine-automatization-bot/bolt"
	"redmine-automatization-bot/handler"
	"redmine-automatization-bot/redmine"
)

var bot *tgbotapi.BotAPI
var users *UserStorage
var waiters *Waiters
var storage Storage
var redmineApis *RedmineApis

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	storage = bolt.NewStorage("bot")

	users, err = NewUserStorage(storage)
	if err != nil {
		panic(err)
	}

	redmineApis = NewRedmineApis(users)
	waiters = NewWaiters()
}

func main() {
	defer storage.Close()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
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

func handle(message *tgbotapi.Message) {
	waiter, exists := waiters.Find(message.From.ID)
	if exists {
		waiter(message, bot)
		return
	}

	api, exists := redmineApis.Find(message.From.ID)
	if !exists {
		user := users.Find(message.From.ID)
		if user == nil {
			requestRedmineUrl(message)
			return
		}

		if user.RedmineApiKey == "" {
			requestRedmineApiKey(message)
			return
		}

		api, err := redmine.NewApi(user.RedmineUrl, user.RedmineApiKey)
		if err != nil {
			log.Panic(err)
			return
		}

		redmineApis.Save(message.From.ID, api)
	}

	handler.Handle(message, bot, api)
}

func requestRedmineUrl(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please, send base url of your redmine")
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
		return
	}

	userId := message.From.ID
	waiters.Set(userId, func(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
		redmineUrl := message.Text
		user, err := users.Register(userId, redmineUrl)
		if err != nil {
			log.Panic(user, err)
			return
		}

		requestRedmineApiKey(message)
	})
}

func requestRedmineApiKey(message *tgbotapi.Message) {
	err := response(message.Chat.ID, "Please, send your redmine api key")
	if err != nil {
		log.Panic(err)
		return
	}

	userId := message.From.ID
	waiters.Set(userId, func(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
		apiKey := message.Text
		// todo validate key
		err := users.AddApiKey(userId, apiKey)
		if err != nil {
			log.Panic(userId, err)
			return
		}

		user := users.Find(userId)
		api, err := redmine.NewApi(user.RedmineUrl, user.RedmineApiKey)
		if err != nil {
			err := response(message.Chat.ID, err.Error())
			if err != nil {
				log.Panic(err)
				return
			}
		}
		redmineApis.Save(userId, api)

		err = (&handler.Start{}).Handle(message, bot, api)
		if err != nil {
			log.Panic(userId, err)
			return
		}

		waiters.Remove(userId)
	})
}

func response(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := bot.Send(msg)

	return err
}
