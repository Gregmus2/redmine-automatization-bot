package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
)

type Hide struct{}

func init() {
	global.RegisterCommand(&Hide{}, "hide")
}

func (d *Hide) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	return d.HandleCommandRow(message, api), nil
}

func (d *Hide) HandleCommandRow(message *tgbotapi.Message, api *redmine.Api) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Done")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	return msg
}

func (_ *Hide) GetRequiredArgs() []string {
	return []string{}
}
