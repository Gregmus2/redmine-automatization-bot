package handlers

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
	"strconv"
	"strings"
)

type TimeEntry struct{ handler }

func init() {
	global.RegisterCommand(&TimeEntry{}, "time_entry")
}

func (t *TimeEntry) Handle(session *global.SessionData) (tgbotapi.Chattable, error) {
	text := "Enter data in format " + strings.Join(t.ArgsInOrder(), "|") + "\nAvailable activities:\n" + session.Api.Activities.ToText()
	msg := tgbotapi.NewMessage(session.Message.Chat.ID, text)

	t.handleNextTime(t, session)

	return msg, nil
}

func (t *TimeEntry) HandleCommandRow(session *global.SessionData) (tgbotapi.Chattable, error) {
	args := strings.Split(session.Message.Text, "|")
	if len(args) < 3 {
		msg := "not enough arguments, try again"

		return tgbotapi.NewMessage(session.Message.Chat.ID, msg), errors.New(msg)
	}

	msg, exist := t.handlePlaceholders(t, session)
	if exist {
		return msg, nil
	}

	issueId, err := strconv.ParseUint(args[0], 10, 0)
	if err != nil {
		return t.errorResponse(session.Message, "wrong ISSUE_ID argument, try again")
	}

	hours, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return t.errorResponse(session.Message, "wrong HOURS argument, try again.")
	}

	activityId, err := strconv.ParseUint(args[2], 10, 8)
	if err != nil {
		return t.errorResponse(session.Message, "wrong ACTIVITY_ID argument, try agai")
	}

	var comment string
	if len(args) == 3 {
		comment = ""
	} else {
		comment = args[3]
	}

	timeEntryBody := redmine.TimeEntryBody{
		IssueId:    uint(issueId),
		Hours:      float32(hours),
		ActivityId: uint8(activityId),
		Comments:   comment,
	}
	_, err = session.Api.CreateTimeEntry(timeEntryBody)
	if err != nil {
		return tgbotapi.NewMessage(session.Message.Chat.ID, err.Error()), err
	}

	return tgbotapi.NewMessage(session.Message.Chat.ID, "Done"), nil
}

func (t *TimeEntry) ArgsInOrder() []string {
	return []string{"ISSUE_ID", "HOURS", "ACTIVITY_ID", "COMMENT"}
}
