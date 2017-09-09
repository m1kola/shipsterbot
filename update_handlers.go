package main

import (
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
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
