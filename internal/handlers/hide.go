package handlers

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
)

type Hide struct{}

func init() {
	global.RegisterCommand(&Hide{}, "hide")
}

func (d *Hide) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Done")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	return msg, nil
}

func (_ *Hide) ValidateArgs(args []string) error {
	if len(args) == 0 {
		return nil
	}

	return errors.New("this command have no arguments")
}

func (_ *Hide) GetRequiredArgs() []string {
	return []string{}
}
