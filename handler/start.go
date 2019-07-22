package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/redmine"
)

type Start struct {}

func init() {
	var handler Start
	RegisterCommand(&handler, "start")
}

func (d *Start) Handle(message *tgbotapi.Message, bot *tgbotapi.BotAPI, api *redmine.Api) error {
	_, err := api.CreateTimeEntry(35300, 0.25, 14, "Кофе")
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Choose command")
	buttons := make([]tgbotapi.KeyboardButton, len(textHandlers))
	for command := range textHandlers {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(command))
	}
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
