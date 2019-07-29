package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
)

type Show struct{}

func init() {
	global.RegisterCommand(&Show{}, "show")
}

func (d *Show) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	return d.HandleCommandRow(message, api), nil
}

func (d *Show) HandleCommandRow(message *tgbotapi.Message, api *redmine.Api) tgbotapi.Chattable {
	names := global.TS.GetTemplateNames(message.From.ID)
	commandsCount := len(names)
	if commandsCount == 0 {
		return tgbotapi.NewMessage(
			message.Chat.ID,
			"You haven't templates yet. Use /create_template command to create a new one",
		)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Done")
	buttons := make([]tgbotapi.KeyboardButton, commandsCount)
	for _, name := range names {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(name))
	}
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

	return msg
}

func (_ *Show) GetRequiredArgs() []string {
	return []string{}
}
