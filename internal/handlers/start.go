package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
)

type Start struct{}

func init() {
	var handler Start
	global.RegisterCommand(&handler, "start")
}

func (d *Start) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Choose command")
	buttons := make([]tgbotapi.KeyboardButton, len(global.TextHandlers))
	for command := range global.TextHandlers {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(command))
	}
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

	return msg, nil
}
