package global

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/redmine"
)

type Handler interface {
	Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error)
}
