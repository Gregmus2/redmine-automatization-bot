package handlers

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
	"strconv"
	"strings"
)

type TimeEntry struct{}

func init() {
	global.RegisterCommand(&TimeEntry{}, "time_entry")
}

func (d *TimeEntry) Handle(message *tgbotapi.Message, api *redmine.Api) (tgbotapi.Chattable, error) {
	text := "Enter data in format ISSUE_ID HOURS ACTIVITY_ID COMMENT\nAvailable activities:\n" + api.Activities.ToText()
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	global.Waiter.Set(message.From.ID, func(message *tgbotapi.Message) tgbotapi.Chattable {
		args := strings.Split(message.Text, " ")
		if len(args) < 3 {
			return tgbotapi.NewMessage(message.Chat.ID, "Not enough arguments, try again")
		}

		issueId, err := strconv.ParseUint(args[0], 10, 0)
		if err != nil {
			return tgbotapi.NewMessage(message.Chat.ID, "Wrong ISSUE_ID argument, try again")
		}

		hours, err := strconv.ParseFloat(args[1], 32)
		if err != nil {
			return tgbotapi.NewMessage(message.Chat.ID, "Wrong HOURS argument, try again")
		}

		activityId, err := strconv.ParseUint(args[2], 10, 8)
		if err != nil {
			return tgbotapi.NewMessage(message.Chat.ID, "Wrong ACTIVITY_ID argument, try again")
		}

		var comment string
		if len(args) == 3 {
			comment = ""
		} else {
			comment = args[3]
		}

		_, err = api.CreateTimeEntry(uint(issueId), float32(hours), uint8(activityId), comment)
		if err != nil {
			return tgbotapi.NewMessage(message.Chat.ID, err.Error())
		}

		return tgbotapi.NewMessage(message.Chat.ID, "Done")
	})

	return msg, nil
}

// todo need to move validate logic to some general method
func (_ *TimeEntry) ValidateArgs(args []string) error {
	if len(args) < 3 {
		return errors.New("not enough arguments, try again")
	}

	_, err := strconv.ParseUint(args[0], 10, 0)
	if err != nil {
		return errors.New("wrong ISSUE_ID argument, try again")
	}

	_, err = strconv.ParseFloat(args[1], 32)
	if err != nil {
		return errors.New("wrong HOURS argument, try again")
	}

	_, err = strconv.ParseUint(args[2], 10, 8)
	if err != nil {
		return errors.New("wrong ACTIVITY_ID argument, try again")
	}

	return nil
}

func (_ *TimeEntry) GetRequiredArgs() []string {
	return []string{"ISSUE_ID", "HOURS", "ACTIVITY_ID", "COMMENT"}
}
