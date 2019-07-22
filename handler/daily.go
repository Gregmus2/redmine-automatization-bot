package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/redmine"
)

type Daily struct {}

func init() {
	var handler Daily
	RegisterText(&handler, "daily")
}

func (d *Daily) Handle(message *tgbotapi.Message, bot *tgbotapi.BotAPI, api *redmine.Api) error {
	_, err := api.CreateTimeEntry(35298, 0.25, 14, "Дейли")
	if err != nil {
		return err
	}

	simpleResponse(message, "Done", bot)

	return nil
}
