package mock_telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// MessageCommandMockSetup creates fake message with a command
func MessageCommandMockSetup(command, commandArgs string) *tgbotapi.Message {
	commandWithSlash := fmt.Sprintf("/%s", command)

	text := commandWithSlash
	if commandArgs != "" {
		text = fmt.Sprintf("%s %s", commandWithSlash, commandArgs)
	}

	message := &tgbotapi.Message{
		Entities: &[]tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: len(commandWithSlash),
			},
		},
		Text: text,
	}

	return message
}
