package handlers

import (
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/utils"
	"strings"
)

type handler struct{}

func (_ *handler) handlePlaceholders(h global.Handler, session *global.SessionData) (tgbotapi.Chattable, bool) {
	if strings.Index(session.Message.Text, "?") == -1 {
		return nil, false
	}

	args := strings.Split(session.Message.Text, " ")
	// собираем переменные аргументы, чтобы запросить у пользователя их значения
	requiredArgs := h.GetRequiredArgs()
	missedArguments := make([]string, len(args))
	for index, argument := range args {
		if argument == "?" {
			missedArguments = append(missedArguments, requiredArgs[index])
		}
	}

	global.Waiter.Set(session.Message.From.ID, func(message *tgbotapi.Message) tgbotapi.Chattable {
		args := strings.Split(message.Text, " ")
		formatString := strings.Replace(session.Message.Text, "?", "%s", -1)
		message.Text = fmt.Sprintf(formatString, utils.Iface(args)...)
		tmpSession := global.SessionData{
			Message: message,
			Api:     session.Api,
		}
		response, err := h.HandleCommandRow(&tmpSession)
		if err != nil {
			return tgbotapi.NewMessage(session.Message.Chat.ID, err.Error())
		}

		global.Waiter.Remove(session.Message.From.ID)

		return response
	})

	return tgbotapi.NewMessage(
		session.Message.Chat.ID,
		"Please, send variable values: "+strings.Join(missedArguments, " "),
	), true
}

// It will handle next user text message
func (_ *handler) handleNextTime(d global.Handler, session *global.SessionData) {
	global.Waiter.Set(session.Message.From.ID, func(message *tgbotapi.Message) tgbotapi.Chattable {
		msg, err := d.HandleCommandRow(session)
		if err == nil {
			global.Waiter.Remove(session.Message.From.ID)
		}

		return msg
	})
}

func (_ *handler) errorResponse(message *tgbotapi.Message, msg string) (tgbotapi.Chattable, error) {
	return tgbotapi.NewMessage(message.Chat.ID, msg), errors.New(msg)
}
