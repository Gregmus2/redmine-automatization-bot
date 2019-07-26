package global

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/redmine"
)

type Handler interface {
	Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error)
	ValidateArgs(args []string) error
	GetRequiredArgs() []string
}

type Storage interface {
	Close()
	GetAll(collection string) (map[string]string, error)
	GetAllRaw(collection string) (map[string][]byte, error)
	Put(collection string, key string, value []byte) error
	CreateCollectionIfNotExist(collection string)
}
