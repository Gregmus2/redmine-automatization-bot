package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
)

type Hide struct{}

func init() {
	global.RegisterCommand(&Hide{}, "hide")
}

func (d *Hide) Handle(session *global.SessionData) (tgbotapi.Chattable, error) {
	return d.HandleCommandRow(session)
}

func (d *Hide) HandleCommandRow(session *global.SessionData) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(session.Message.Chat.ID, "Done")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	return msg, nil
}

func (_ *Hide) GetRequiredArgs() []string {
	return []string{}
}
