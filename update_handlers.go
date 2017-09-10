package main

import (
	"fmt"
	"log"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// HandleUpdates starts infinite loop that receives
// updates from telegram
func HandleUpdates(bot *tgbotapi.BotAPI, updates <-chan tgbotapi.Update) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Entities == nil {
			log.Printf(
				"Message without commands received: %s",
				update.Message.Text)
			go handleUnknownCommand(bot, update.Message)
			continue
		}

		log.Printf("Message received: %s", update.Message.Text)

		// sendHelpMessage indicates that we should
		// send a help message after processing.
		var sendHelpMessage bool
		for _, entity := range *update.Message.Entities {
			if entity.Type != "bot_command" {
				sendHelpMessage = true
				continue
			}
			sendHelpMessage = false

			// Get command name without the leading slash
			// TODO: handle out of range panic here
			botCommand := update.Message.Text[entity.Offset+1 : entity.Offset+entity.Length]

			switch {
			case botCommand == "help" || botCommand == "start":
				go handleStart(bot, update.Message)
			case botCommand == "add":
				go handleAdd(bot, update.Message)
			case botCommand == "del":
				go handleDel(bot, update.Message)
			default:
				sendHelpMessage = true
			}

			// Support only one command per message, for now
			break
		}

		if sendHelpMessage {
			log.Print("No supported bot commands found")
			go handleUnknownCommand(bot, update.Message)
		}
	}
}

func _handleHelpMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, isStart bool) {
	var greeting string
	if isStart {
		greeting = "Hi %s,"
	} else {
		greeting = "%s, I'm very sorry, but I don't understand you."
	}
	textTemplate := greeting + `

	I can help you to manage your shopping list.

	You can control me by sending these commands:

	*Shopping list*

	/add - Adds an item into your shopping list
	/del - Removesan item from your shopping list
	`
	text := fmt.Sprintf(textTemplate, message.From.FirstName)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msg)
}

func handleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	_handleHelpMessage(bot, message, true)
}

func handleUnknownCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	_handleHelpMessage(bot, message, false)
}

func handleAdd(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	text := "Not implemented yet, sorry"

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msg)
}

func handleDel(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	text := "Not implemented yet, sorry"

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msg)
}
