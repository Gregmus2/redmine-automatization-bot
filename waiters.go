package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"sync"
)

// it's handle some user requests, when we wait them for specific action
type Waiters struct {
	mx sync.Mutex
	m  map[int]func(message *tgbotapi.Message, bot *tgbotapi.BotAPI)
}

func NewWaiters() *Waiters {
	return &Waiters{m: make(map[int]func(message *tgbotapi.Message, bot *tgbotapi.BotAPI))}
}

func (w *Waiters) Find(key int) (func(message *tgbotapi.Message, bot *tgbotapi.BotAPI), bool) {
	w.mx.Lock()
	defer w.mx.Unlock()
	val, ok := w.m[key]

	return val, ok
}

func (w *Waiters) Set(key int, value func(message *tgbotapi.Message, bot *tgbotapi.BotAPI)) {
	w.mx.Lock()
	defer w.mx.Unlock()
	w.m[key] = value
}

func (w *Waiters) Remove(key int) {
	w.mx.Lock()
	defer w.mx.Unlock()
	delete(w.m, key)
}
