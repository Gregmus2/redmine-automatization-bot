package global

import (
	"log"
	"redmine-automatization-bot/internal/bolt"
	"strings"
)

var CommandHandlers = make(map[string]Handler)
var Waiter *WaiterStorage
var TS *TemplateStorage
var RA *RedmineApis
var Users *UserStorage
var Stor Storage

func init() {
	Stor = bolt.NewStorage("bot")

	var err error
	Users, err = NewUserStorage(Stor)
	if err != nil {
		panic(err)
	}
	RA = NewRedmineApis(Users)

	Waiter = NewWaiters()
	TS, err = NewTemplateStorage(Stor)
	if err != nil {
		panic(err)
	}
}

func RegisterCommand(handler Handler, command string) {
	if _, exists := CommandHandlers[command]; exists {
		log.Fatalln(command, "Command already registered")
	}

	log.Println("Register", command, "command")
	CommandHandlers[command] = handler
}

func GetCommandsHelp() string {
	help := ""
	for command, handler := range CommandHandlers {
		help = help + command + "|" + strings.Join(handler.ArgsInOrder(), "|") + "\n"
	}

	return help
}
