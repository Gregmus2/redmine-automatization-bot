package internal

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"redmine-automatization-bot/internal/global"
	_ "redmine-automatization-bot/internal/handlers"
	"redmine-automatization-bot/internal/redmine"
	"strings"
)

// todo unit тесты
func handle(message *tgbotapi.Message) {
	api, status := authorize(message)
	if status == false {
		return
	}

	session := global.SessionData{
		Message: message,
		Api:     api,
	}

	if message.IsCommand() {
		global.Waiter.Remove(message.From.ID)
		handleCommand(&session)
		return
	}

	exists := handleWaiters(message)
	if exists == true {
		return
	}

	handleText(&session)
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

func handleCommand(session *global.SessionData) {
	handler, exists := global.CommandHandlers[session.Message.Command()]
	if !exists {
		return
	}

	msg, err := handler.Handle(session)
	if err != nil {
		log.Panic("Error on handle", err)
	}

	_, err = Bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

func handleText(session *global.SessionData) {
	commandString, exists := global.TS.GetTemplateCommand(session.Message.From.ID, session.Message.Text)
	if !exists {
		return
	}

	commandIndex := strings.Index(commandString, "|")
	var command string
	if commandIndex == -1 {
		command = commandString
		session.Message.Text = ""
	} else {
		command = commandString[:commandIndex]
		session.Message.Text = commandString[commandIndex+1:]
	}

	handler, exists := global.CommandHandlers[command]
	if !exists {
		log.Panic(command)
		return
	}

	msg, _ := handler.HandleCommandRow(session)

	_, err := Bot.Send(msg)
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
