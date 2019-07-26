package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
)

type RemoveTemplate struct{}

func init() {
	global.RegisterCommand(&RemoveTemplate{}, "remove_template")
}

func (d *RemoveTemplate) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	text := "Enter name of template which you want to remove"
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	global.Waiter.Set(message.From.ID, func(message *tgbotapi.Message) tgbotapi.Chattable {
		err := global.TS.RemoveTemplate(message.From.ID, message.Text)
		if err != nil {
			return tgbotapi.NewMessage(message.Chat.ID, err.Error())
		}

		return tgbotapi.NewMessage(
			message.Chat.ID,
			"Template was delete",
		)
	})

	return msg, nil
}

func (_ *RemoveTemplate) ValidateArgs(args []string) error {
	return nil
}

func (_ *RemoveTemplate) GetRequiredArgs() []string {
	return []string{"NAME"}
}
