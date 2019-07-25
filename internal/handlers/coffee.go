package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
)

type Coffee struct{}

func init() {
	global.RegisterText(&Coffee{}, "coffee")
}

func (d *Coffee) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	_, err := api.CreateTimeEntry(35300, 0.25, 14, "Кофе")
	if err != nil {
		return nil, err
	}

	return tgbotapi.NewMessage(message.Chat.ID, "Done"), nil
}
