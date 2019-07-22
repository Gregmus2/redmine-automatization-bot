package handler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"redmine-automatization-bot/redmine"
)

type Handler interface {
	Handle(message *tgbotapi.Message, bot *tgbotapi.BotAPI, api *redmine.Api) error
}

var commandHandlers = make(map[string]Handler)
var textHandlers = make(map[string]Handler)

func Handle(message *tgbotapi.Message, bot *tgbotapi.BotAPI, api *redmine.Api) {
	handlers, key := resolveHandler(message)
	handler, exists := handlers[key]
	if !exists {
		return
	}

	err := handler.Handle(message, bot, api)
	if err != nil {
		log.Panic("Error on handle", err)
	}
}

func resolveHandler(message *tgbotapi.Message) (map[string]Handler, string) {
	if message.IsCommand() {
		return commandHandlers, message.Command()
	}

	return textHandlers, message.Text
}

func simpleResponse(inputMessage *tgbotapi.Message, text string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, text)

	_, err := bot.Send(msg)
	if err != nil {
		log.Panic("Error on send response", err)
	}
}

func RegisterCommand(handler Handler, command string) {
	if _, exists := commandHandlers[command]; exists {
		log.Fatalln(command, "Command already registered")
	}

	log.Println("Register", command, "command")
	commandHandlers[command] = handler
}

func RegisterText(handler Handler, text string) {
	if _, exists := textHandlers[text]; exists {
		log.Fatalln(text, "Text command already registered")
	}

	log.Println("Register", text, "text command")
	textHandlers[text] = handler
}