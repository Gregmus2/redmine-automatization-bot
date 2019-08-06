package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
)

type Show struct{ handler }

func init() {
	global.RegisterCommand(&Show{}, "show")
}

func (d *Show) Handle(session *global.SessionData) (tgbotapi.Chattable, error) {
	return d.HandleCommandRow(session)
}

func (d *Show) HandleCommandRow(session *global.SessionData) (tgbotapi.Chattable, error) {
	names := global.TS.GetTemplateNames(session.Message.From.ID)
	commandsCount := len(names)
	if commandsCount == 0 {
		return d.errorResponse(session.Message, "you haven't templates yet. Use /create_template command to create a new one")
	}

	msg := tgbotapi.NewMessage(session.Message.Chat.ID, "Done")
	buttons := make([]tgbotapi.KeyboardButton, commandsCount)
	for _, name := range names {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(name))
	}
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

	return msg, nil
}

func (_ *Show) ArgsInOrder() []string {
	return []string{}
}
