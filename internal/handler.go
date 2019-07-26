package internal

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"redmine-automatization-bot/internal/global"
	_ "redmine-automatization-bot/internal/handlers"
	"redmine-automatization-bot/internal/redmine"
)

func handle(message *tgbotapi.Message) {
	api, status := authorize(message)
	if status == false {
		return
	}

	if message.IsCommand() {
		global.Waiter.Remove(message.From.ID)
		handleCommand(message, api)
		return
	}

	exists := handleWaiters(message)
	if exists == true {
		return
	}

	handleText(message, api)
}

func authorize(message *tgbotapi.Message) (*redmine.Api, bool) {
	api, exists := global.RA.Find(message.From.ID)
	if !exists {
		exists := handleWaiters(message)
		if exists == true {
			return nil, false
		}

		user := global.Users.Find(message.From.ID)
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

		global.RA.Save(message.From.ID, api)
	}

	return api, true
}

func handleWaiters(message *tgbotapi.Message) bool {
	callable, exists := global.Waiter.Find(message.From.ID)
	if exists {
		msg := callable(message)
		if msg == nil {
			return true
		}

		_, err := Bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}

		return true
	}

	return false
}

func handleCommand(message *tgbotapi.Message, api *redmine.Api) {
	handler, exists := global.CommandHandlers[message.Command()]
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
	global.Waiter.Set(userId, func(message *tgbotapi.Message) tgbotapi.Chattable {
		redmineUrl := message.Text
		user, err := global.Users.Register(userId, redmineUrl)
		if err != nil {
			log.Panic(user, err)
			return nil
		}

		requestRedmineApiKey(message)

		return tgbotapi.NewMessage(message.Chat.ID, "Please, send your redmine api key")
	})
}

func requestRedmineApiKey(message *tgbotapi.Message) {
	userId := message.From.ID
	global.Waiter.Set(userId, func(message *tgbotapi.Message) tgbotapi.Chattable {
		apiKey := message.Text
		// todo validate key
		err := global.Users.AddApiKey(userId, apiKey)
		if err != nil {
			log.Panic(userId, err)
			return nil
		}

		user := global.Users.Find(userId)
		api, err := redmine.NewApi(user.RedmineUrl, user.RedmineApiKey)
		if err != nil {
			return tgbotapi.NewMessage(message.Chat.ID, err.Error())
		}

		global.RA.Save(userId, api)
		global.Waiter.Remove(userId)

		return tgbotapi.NewMessage(message.Chat.ID, "Welcome")
	})
}
