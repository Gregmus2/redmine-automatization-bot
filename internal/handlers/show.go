package handlers

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
)

type Show struct{}

func init() {
	global.RegisterCommand(&Show{}, "show")
}

func (d *Show) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	names := global.TS.GetTemplateNames(message.From.ID)
	commandsCount := len(names)
	if commandsCount == 0 {
		return tgbotapi.NewMessage(
			message.Chat.ID,
			"You haven't templates yet. Use /create_template command to create a new one",
		), nil
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Done")
	buttons := make([]tgbotapi.KeyboardButton, commandsCount)
	for _, name := range names {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(name))
	}
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

	return msg, nil
}

func (_ *Show) ValidateArgs(args []string) error {
	if len(args) == 0 {
		return nil
	}

	return errors.New("this command have no arguments")
}

func (_ *Show) GetRequiredArgs() []string {
	return []string{}
}
