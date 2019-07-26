package handlers

import (
	"errors"
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
		args := strings.Split(message.Text, " ")
		name := args[0]
		command := args[1]
		commandArgs := args[2:]
		handler, exists := global.CommandHandlers[command]
		if !exists {
			return tgbotapi.NewMessage(message.Chat.ID, "Wrong command")
		}

		err := handler.ValidateArgs(commandArgs)
		if err != nil {
			return tgbotapi.NewMessage(message.Chat.ID, err.Error())
		}

		err = global.TS.AddTemplate(message.From.ID, name, strings.Join(args[1:], " "))
		if err != nil {
			return tgbotapi.NewMessage(message.Chat.ID, err.Error())
		}

		return tgbotapi.NewMessage(
			message.Chat.ID,
			"New template created, you can use /show command to show all your templates "+
				"or send name of your template as a message to the bot",
		)
	})

	return msg, nil
}

func (_ *CreateTemplate) ValidateArgs(args []string) error {
	if len(args) > 1 {
		return nil
	}

	return errors.New("please, enter data in format COMMAND_WITHOUT_SLASH ARGS_OF_COMMAND")
}

func (_ *CreateTemplate) GetRequiredArgs() []string {
	return []string{"NAME", "COMMAND", "ARGS"}
}
