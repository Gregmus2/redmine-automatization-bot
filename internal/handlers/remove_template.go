package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
)

type RemoveTemplate struct{}

func init() {
	global.RegisterCommand(&RemoveTemplate{}, "remove_template")
}

func (d *RemoveTemplate) Handle(session *global.SessionData) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(session.Message.Chat.ID, "Enter name of template which you want to remove")

	global.Waiter.Set(session.Message.From.ID, func(message *tgbotapi.Message) tgbotapi.Chattable {
		msg, err := d.HandleCommandRow(session)
		if err == nil {
			global.Waiter.Remove(session.Message.From.ID)
		}

		return msg
	})

	return msg, nil
}

func (d *RemoveTemplate) HandleCommandRow(session *global.SessionData) (tgbotapi.Chattable, error) {
	err := global.TS.RemoveTemplate(session.Message.From.ID, session.Message.Text)
	if err != nil {
		return tgbotapi.NewMessage(session.Message.Chat.ID, err.Error()), err
	}

	return tgbotapi.NewMessage(
		session.Message.Chat.ID,
		"Template was delete",
	), nil
}

func (_ *RemoveTemplate) ArgsInOrder() []string {
	return []string{"NAME"}
}
