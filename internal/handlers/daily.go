package handlers

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
)

type Daily struct{}

func init() {
	global.RegisterText(&Daily{}, "daily")
}

func (d *Daily) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	_, err := api.CreateTimeEntry(35298, 0.25, 14, "Дейли")
	if err != nil {
		return nil, err
	}

	return tgbotapi.NewMessage(message.Chat.ID, "Done"), nil
}

func (_ *Daily) ValidateArgs(args []string) error {
	if len(args) == 0 {
		return nil
	}

	return errors.New("this command have no arguments")
}

func (_ *Daily) GetRequiredArgs() []string {
	return []string{}
}
