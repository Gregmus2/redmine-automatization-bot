package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"strings"
)

type CreateTemplate struct{ handler }

func init() {
	global.RegisterCommand(&CreateTemplate{}, "create_template")
}

func (d *CreateTemplate) Handle(session *global.SessionData) (tgbotapi.Chattable, error) {
	text := "Enter data in format " + strings.Join(d.ArgsInOrder(), "|") + "\nAvailable commands:\n" + global.GetCommandsHelp()
	msg := tgbotapi.NewMessage(session.Message.Chat.ID, text)

	d.handleNextTime(d, session)

	return msg, nil
}

func (d *CreateTemplate) HandleCommandRow(session *global.SessionData) (tgbotapi.Chattable, error) {
	args := strings.Split(session.Message.Text, "|")
	name := args[0]

	err := global.TS.AddTemplate(session.Message.From.ID, name, strings.Join(args[1:], "|"))
	if err != nil {
		return tgbotapi.NewMessage(session.Message.Chat.ID, err.Error()), err
	}

	return tgbotapi.NewMessage(
		session.Message.Chat.ID,
		"New template created, you can use /show command to show all your templates "+
			"or send name of your template as a message to the bot",
	), nil
}

func (_ *CreateTemplate) ArgsInOrder() []string {
	return []string{"NAME", "COMMAND", "ARGS"}
}
