package global

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"sync"
)

// it's handle some user requests, when we wait them for specific action
type WaiterStorage struct {
	mx sync.Mutex
	m  map[int]func(message *tgbotapi.Message) tgbotapi.Chattable
}

func NewWaiters() *WaiterStorage {
	return &WaiterStorage{m: make(map[int]func(message *tgbotapi.Message) tgbotapi.Chattable)}
}

func (w *WaiterStorage) Find(key int) (func(message *tgbotapi.Message) tgbotapi.Chattable, bool) {
	w.mx.Lock()
	defer w.mx.Unlock()
	val, ok := w.m[key]

	return val, ok
}

func (w *WaiterStorage) Set(key int, value func(message *tgbotapi.Message) tgbotapi.Chattable) {
	w.mx.Lock()
	defer w.mx.Unlock()
	w.m[key] = value
}

func (w *WaiterStorage) Remove(key int) {
	w.mx.Lock()
	defer w.mx.Unlock()
	delete(w.m, key)
}
