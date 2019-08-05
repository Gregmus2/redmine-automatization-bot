package global

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/redmine"
)

type Handler interface {
	Handle(session *SessionData) (tgbotapi.Chattable, error)
	HandleCommandRow(session *SessionData) (tgbotapi.Chattable, error)
	GetRequiredArgs() []string
}

type Storage interface {
	Close()
	GetAll(collection string) (map[string]string, error)
	GetAllRaw(collection string) (map[string][]byte, error)
	Put(collection string, key string, value []byte) error
	CreateCollectionIfNotExist(collection string)
}

type SessionData struct {
	Message *tgbotapi.Message
	Api     *redmine.Api
}
