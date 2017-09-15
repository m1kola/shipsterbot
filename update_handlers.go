package main

import (
	"fmt"
	"log"

	"github.com/m1kola/telegram_shipsterbot/storage"
	"github.com/m1kola/telegram_shipsterbot/types"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TODO: Define interface for handlers

// HandleUpdates starts infinite loop that receives
// updates from Telegram.
func HandleUpdates(bot *tgbotapi.BotAPI, updates <-chan tgbotapi.Update) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("Message received: %s", update.Message.Text)

		isHandled := handleMessageEntities(bot, update.Message)
		if !isHandled {
			isHandled = handleMessage(bot, update.Message)

			if !isHandled {
				log.Print("No supported bot commands found")
				go handleUnrecognisedMessage(bot, update.Message)
			}
		}
	}
}

// handleMessageEntities returns true if the message is handled
func handleMessageEntities(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	if message.Entities == nil {
		return false
	}

	for _, entity := range *message.Entities {
		if entity.Type != "bot_command" {
			continue
		}

		// Get command name without the leading slash
		// TODO: handle out of range panic here
		botCommand := message.Text[entity.Offset+1 : entity.Offset+entity.Length]

		switch botCommand {
		case "help", "start":
			go handleStart(bot, message)
			return true
		case "add":
			go handleAdd(bot, message)
			return true
		case "list":
			go handleList(bot, message)
			return true
		case "del":
			go handleDel(bot, message)
			return true
		default:
			continue
		}
	}

	return false
}

// handleMessage handles unfinshed operations.
// Normally we listen to user's commands (`MessageEntities` of type `bot_command`)
// or using keyboard, but in some cases we need to handle message text.
// For example, when user asks us to add an item into the shopping list
func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	session, ok := storage.GetUnfinishedOperation(message.From.ID)

	if ok {
		switch session.Operation {
		case types.OperationAdd:
			storage.DeleteUnfinishedOperation(message.From.ID)
			handleAddSession(bot, message)
			return true
		}
	}

	return false
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
	/list - Displays items from your shopping list
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

func handleUnrecognisedMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	_handleHelpMessage(bot, message, false)
}

func handleAdd(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	storage.AddUnfinishedOperation(message.From.ID, types.OperationAdd)

	text := "Ok, what do you want to add into your shopping list?"
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msg)
}

func handleList(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var text string
	chatID := message.Chat.ID

	chatItems, ok := storage.GetShoppingItems(chatID)
	if !ok || chatItems == nil {
		text = "Your shopping list is empty. Who knows, maybe it's a good thing"
	} else {
		text = "Here is the list item in your shopping list:\n\n"

		// TODO: Add space offset for indexes
		for index, item := range *chatItems {
			text += fmt.Sprintf("%d. %s\n", index+1, item.Name)
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msg)
}

func handleAddSession(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// TODO: We should, probably cleanup text (remove @UserNameBot, etc)
	itemName := message.Text

	storage.AddShoppingItemIntoShoppingList(message.Chat.ID, &types.ShoppingItem{
		Name:      itemName,
		IsActive:  true,
		CreatedBy: message.From.ID})

	text := "Lovely! I've added \"%s\" into your shopping list. Anything else?"
	text = fmt.Sprintf(text, itemName)
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
