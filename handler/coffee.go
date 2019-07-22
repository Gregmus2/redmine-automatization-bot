package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/redmine"
)

type Coffee struct {}

func init() {
	var handler Coffee
	RegisterText(&handler, "coffee")
}

func (d *Coffee) Handle(message *tgbotapi.Message, bot *tgbotapi.BotAPI, api *redmine.Api) error {
	_, err := api.CreateTimeEntry(35300, 0.25, 14, "Кофе")
	if err != nil {
		return err
	}

	simpleResponse(message, "Done", bot)

	return nil
}
