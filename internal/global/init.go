package global

import "log"

var CommandHandlers = make(map[string]Handler)
var TextHandlers = make(map[string]Handler)

func RegisterCommand(handler Handler, command string) {
	if _, exists := CommandHandlers[command]; exists {
		log.Fatalln(command, "Command already registered")
	}

	log.Println("Register", command, "command")
	CommandHandlers[command] = handler
}

func RegisterText(handler Handler, text string) {
	if _, exists := TextHandlers[text]; exists {
		log.Fatalln(text, "Text command already registered")
	}

	log.Println("Register", text, "text command")
	TextHandlers[text] = handler
}
