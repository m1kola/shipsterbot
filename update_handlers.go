package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/m1kola/telegram_shipsterbot/storage"
	"github.com/m1kola/telegram_shipsterbot/types"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TODO: Define interface for handlers

// HandleUpdates starts infinite loop that receives
// updates from Telegram.
func HandleUpdates(bot *tgbotapi.BotAPI, updates <-chan tgbotapi.Update) {
	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallbackQuery(bot, &update)
		}

		if update.Message != nil {
			handleMessage(bot, update.Message)
		}
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if update.CallbackQuery.Message == nil {
		return
	}

	dataPieces := strings.SplitN(update.CallbackQuery.Data, ":", 2)
	if len(dataPieces) != 2 {
		return
	}

	botCommand := dataPieces[0]

	switch botCommand {
	case "del":
		go handleDelCallbackQuery(bot, update.CallbackQuery, dataPieces[1])
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("Message received: %s", message.Text)

	isHandled := handleMessageEntities(bot, message)
	if !isHandled {
		isHandled = handleMessageText(bot, message)

		if !isHandled {
			log.Print("No supported bot commands found")
			go handleUnrecognisedMessage(bot, message)
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
		commandStartPos := entity.Offset + 1
		commandEndPos := entity.Offset + entity.Length
		if commandStartPos < 0 || commandEndPos > len(message.Text) {
			return false
		}
		botCommand := message.Text[commandStartPos:commandEndPos]

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

// handleMessageText handles unfinshed operations.
// Normally we listen to user's commands (`MessageEntities` of type `bot_command`)
// or using keyboard, but in some cases we need to handle message text.
// For example, when user asks us to add an item into the shopping list
func handleMessageText(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	session, ok := storage.GetUnfinishedCommand(message.From.ID)

	if ok {
		switch session.Command {
		case types.CommandAddShoppingItem:
			storage.DeleteUnfinishedCommand(message.From.ID)
			go handleAddSession(bot, message)
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
	storage.AddUnfinishedCommand(message.From.ID,
		types.CommandAddShoppingItem)

	text := "Ok, what do you want to add into your shopping list?"
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	bot.Send(msg)
}

func handleList(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var text string
	chatID := message.Chat.ID

	chatItems, ok := storage.GetShoppingItems(chatID)
	if !ok || len(chatItems) == 0 {
		text = "Your shopping list is empty. Who knows, maybe it's a good thing"
	} else {
		offset := len(strconv.Itoa(len(chatItems)))
		listItemFormat := fmt.Sprintf("%%%dd. %%s\n", offset)
		// TODO: replace index with id/pk and hide storage logic in the storage module
		for index, item := range chatItems {
			text += fmt.Sprintf(listItemFormat, index+1, item.Name)
		}

		text = fmt.Sprintf(
			"%s\n\n```\n%s```",
			"Here is the list item in your shopping list:",
			text)
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
	bot.Send(msg)
}

func handleDel(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var text string
	var itemButtons []tgbotapi.InlineKeyboardButton

	chatID := message.Chat.ID
	chatItems, ok := storage.GetShoppingItems(chatID)
	if !ok || len(chatItems) == 0 {
		text = "Your shopping list is empty. No need to delete items :)"
	} else {
		for index, item := range chatItems {
			callbackData := fmt.Sprintf("del:%s", strconv.Itoa(index+1))
			itemButton := tgbotapi.NewInlineKeyboardButtonData(item.Name, callbackData)
			itemButtons = append(itemButtons, itemButton)
		}

		text = "Ok, what item do you want to delete from your shopping list?"
	}

	msg := tgbotapi.NewMessage(chatID, text)
	if len(itemButtons) > 0 {
		msg.BaseChat.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(itemButtons)
	}
	bot.Send(msg)
}

func handleDelCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, data string) {
	bot.AnswerCallbackQuery(tgbotapi.NewCallback(
		callbackQuery.ID, ""))

	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID
	itemID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return
	}

	var text string
	item, ok := storage.GetShoppingItem(chatID, itemID)
	if ok {
		storage.DeleteShoppingItem(chatID, itemID)

		text = "It's nice to see that you think that you don't"
		text += "need this \"%s\" thing. "
		text += "I've removed it from your shopping list.\n\n"
		text += "Can I do anything else for you?"
		text = fmt.Sprintf(text, item.Name)
	} else {
		text = "Can't find an item, sorry."
	}

	// Edit previous message to hide the keyboard
	{
		msg := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			messageID,
			tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{}))
		bot.Send(msg)
	}

	// Send deletion confimration text
	{
		msg := tgbotapi.NewMessage(chatID, text)
		bot.Send(msg)
	}
}
