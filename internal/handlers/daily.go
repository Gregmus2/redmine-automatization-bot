package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
)

type Daily struct{}

func init() {
	var handler Daily
	global.RegisterText(&handler, "daily")
}

func (d *Daily) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	_, err := api.CreateTimeEntry(35298, 0.25, 14, "Дейли")
	if err != nil {
		return nil, err
	}

	return tgbotapi.NewMessage(message.Chat.ID, "Done"), nil
}
