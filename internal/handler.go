package internal

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"redmine-automatization-bot/internal/bolt"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/handlers"
	"redmine-automatization-bot/internal/redmine"
)

var waiter *WaiterStorage
var redmineApis *RedmineApis
var users *UserStorage
var storage Storage

func init() {
	storage = bolt.NewStorage("bot")

	var err error
	users, err = NewUserStorage(storage)
	if err != nil {
		panic(err)
	}
	redmineApis = NewRedmineApis(users)
	waiter = NewWaiters()
}

func handle(message *tgbotapi.Message) {
	api, status := authorize(message)
	if status == false {
		return
	}

	if message.IsCommand() {
		exists := handleCommand(message, api)
		if exists {
			return
		}
	}

	exists := handleWaiters(message)
	if exists == true {
		return
	}

	handleText(message, api)
}

func authorize(message *tgbotapi.Message) (*redmine.Api, bool) {
	api, exists := redmineApis.Find(message.From.ID)
	if !exists {
		exists := handleWaiters(message)
		if exists == true {
			return nil, false
		}

		user := users.Find(message.From.ID)
		if user == nil {
			requestRedmineUrl(message)
			return nil, false
		}

		if user.RedmineApiKey == "" {
			requestRedmineApiKey(message)
			return nil, false
		}

		api, err := redmine.NewApi(user.RedmineUrl, user.RedmineApiKey)
		if err != nil {
			log.Panic(err)
			return nil, false
		}

		redmineApis.Save(message.From.ID, api)
	}

	return api, true
}

func handleWaiters(message *tgbotapi.Message) bool {
	waiter, exists := waiter.Find(message.From.ID)
	if exists {
		waiter(message)

		return true
	}

	return false
}

func handleCommand(message *tgbotapi.Message, api *redmine.Api) bool {
	handler, exists := global.CommandHandlers[message.Command()]
	if !exists {
		return false
	}

	msg, err := handler.Handle(message, api)
	if err != nil {
		log.Panic("Error on handle", err)
	}

	_, err = Bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	return true
}

func handleText(message *tgbotapi.Message, api *redmine.Api) {
	handler, exists := global.TextHandlers[message.Text]
	if !exists {
		return
	}

	msg, err := handler.Handle(message, api)
	if err != nil {
		log.Panic("Error on handle", err)
	}

	_, err = Bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

func requestRedmineUrl(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please, send base url of your redmine")
	_, err := Bot.Send(msg)
	if err != nil {
		log.Panic(err)
		return
	}

	userId := message.From.ID
	waiter.Set(userId, func(message *tgbotapi.Message) {
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
	waiter.Set(userId, func(message *tgbotapi.Message) {
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

		msg, err := (&handlers.Start{}).Handle(message, api)
		if err != nil {
			log.Panic(userId, err)
			return
		}

		_, err = Bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}

		waiter.Remove(userId)
	})
}

func response(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := Bot.Send(msg)

	return err
}

func SimpleResponse(inputMessage *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, text)

	_, err := Bot.Send(msg)
	if err != nil {
		log.Panic("Error on send response", err)
	}
}
