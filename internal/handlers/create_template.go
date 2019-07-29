package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
	"strings"
)

type CreateTemplate struct{}

func init() {
	global.RegisterCommand(&CreateTemplate{}, "create_template")
}

func (d *CreateTemplate) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	text := "Enter data in format NAME COMMAND ARGS\nAvailable commands:\n" + global.GetCommandsHelp()
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	global.Waiter.Set(message.From.ID, func(message *tgbotapi.Message) tgbotapi.Chattable {
		return d.HandleCommandRow(message, api)
	})

	return msg, nil
}

func (d *CreateTemplate) HandleCommandRow(message *tgbotapi.Message, api *redmine.Api) tgbotapi.Chattable {
	args := strings.Split(message.Text, " ")
	name := args[0]

	err := global.TS.AddTemplate(message.From.ID, name, strings.Join(args[1:], " "))
	if err != nil {
		return tgbotapi.NewMessage(message.Chat.ID, err.Error())
	}

	return tgbotapi.NewMessage(
		message.Chat.ID,
		"New template created, you can use /show command to show all your templates "+
			"or send name of your template as a message to the bot",
	)
}

func (_ *CreateTemplate) GetRequiredArgs() []string {
	return []string{"NAME", "COMMAND", "ARGS"}
}
