package handlers

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"redmine-automatization-bot/internal/global"
	"redmine-automatization-bot/internal/redmine"
	"redmine-automatization-bot/internal/utils"
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
	// todo Waiter.Remove
	global.Waiter.Set(message.From.ID, func(message *tgbotapi.Message) tgbotapi.Chattable {
		return d.HandleCommandRow(message, api)
	})

	return msg, nil
}

func (d *TimeEntry) HandleCommandRow(message *tgbotapi.Message, api *redmine.Api) tgbotapi.Chattable {
	commandText := message.Text
	args := strings.Split(commandText, " ")
	if len(args) < 3 {
		return tgbotapi.NewMessage(message.Chat.ID, "Not enough arguments, try again")
	}

	if strings.Index(commandText, "?") != -1 {
		// собираем переменные аргументы, чтобы запросить у пользователя их значения
		requiredArgs := d.GetRequiredArgs()
		missedArguments := make([]string, len(args))
		for index, argument := range args {
			if argument == "?" {
				missedArguments = append(missedArguments, requiredArgs[index])
			}
		}

		// todo Waiter.Remove
		global.Waiter.Set(message.From.ID, func(message *tgbotapi.Message) tgbotapi.Chattable {
			args := strings.Split(message.Text, " ")
			formatString := strings.Replace(commandText, "?", "%s", -1)
			message.Text = fmt.Sprintf(formatString, utils.Iface(args)...)

			return d.HandleCommandRow(message, api)
		})

		return tgbotapi.NewMessage(
			message.Chat.ID,
			"Please, send variable values: "+strings.Join(missedArguments, " "),
		)
	}

	issueId, err := strconv.ParseUint(args[0], 10, 0)
	if err != nil {
		return tgbotapi.NewMessage(message.Chat.ID, "Wrong ISSUE_ID argument, try again")
	}

	hours, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return tgbotapi.NewMessage(message.Chat.ID, "Wrong HOURS argument, try again.")
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
}

// todo need to move validate logic to some general method
func (_ *TimeEntry) GetRequiredArgs() []string {
	return []string{"ISSUE_ID", "HOURS", "ACTIVITY_ID", "COMMENT"}
}
