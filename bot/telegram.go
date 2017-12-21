package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/m1kola/shipsterbot/models"
	"github.com/m1kola/shipsterbot/storage"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TelegramBotApp is a struct for handeling iteractions
// with the Telegram API
type TelegramBotApp struct {
	AppInterface
	Bot     *tgbotapi.BotAPI
	Storage storage.DataStorageInterface
}

// ListenForWebhook starts infinite loop that receives
// updates from Telegram.
func (bot_app TelegramBotApp) ListenForWebhook() {
	updates := bot_app.Bot.ListenForWebhook(
		fmt.Sprintf("/%s/webhook", bot_app.Bot.Token))

	go bot_app.handleUpdates(updates)
}

func (bot_app TelegramBotApp) handleUpdates(updates <-chan tgbotapi.Update) {
	for update := range updates {
		if update.CallbackQuery != nil {
			go bot_app.handleCallbackQuery(&update)
		}

		if update.Message != nil {
			go bot_app.handleMessage(update.Message)
		}
	}
}

func (bot_app TelegramBotApp) handleCallbackQuery(update *tgbotapi.Update) {
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
		bot_app.handleDelCallbackQuery(update.CallbackQuery, dataPieces[1])
	case "clear":
		bot_app.handleClearCallbackQuery(update.CallbackQuery, dataPieces[1])
	}
}

func (bot_app TelegramBotApp) handleMessage(message *tgbotapi.Message) {
	log.Printf("Message received: \"%s\"", message.Text)

	isHandled := bot_app.handleMessageEntities(message)
	if !isHandled {
		isHandled = bot_app.handleMessageText(message)

		if !isHandled {
			bot_app.handleUnrecognisedMessage(message)
		}
	}
}

// handleMessageEntities returns true if the message is handled
func (bot_app TelegramBotApp) handleMessageEntities(message *tgbotapi.Message) bool {
	if message.Entities == nil {
		return false
	}

	botCommand := message.Command()
	switch botCommand {
	case "help", "start":
		bot_app.handleStart(message)
		return true
	case "add":
		bot_app.handleAdd(message)
		return true
	case "list":
		bot_app.handleList(message)
		return true
	case "del":
		bot_app.handleDel(message)
		return true
	case "clear":
		bot_app.handleClear(message)
		return true
	}

	return false
}

// handleMessageText handles unfinshed operations.
// Normally we listen to user's text commands or inline keyboard,
// but in some cases we need to handle message text.
// For example, when user asks us to add an item into the shopping list.
func (bot_app TelegramBotApp) handleMessageText(message *tgbotapi.Message) bool {
	session, err := bot_app.Storage.GetUnfinishedCommand(message.Chat.ID,
		message.From.ID)

	if err != nil {
		log.Printf(
			"Unable to get an unfinished comamnd (ChatID=%d and UserId=%d): %q",
			message.Chat.ID, message.From.ID, err)
		return false
	}

	switch session.Command {
	case models.CommandAddShoppingItem:
		err := bot_app.Storage.DeleteUnfinishedCommand(message.Chat.ID,
			message.From.ID)

		if err != nil {
			log.Printf(
				"Unable to delete an unfinished comamnd (ChatID=%d and UserId=%d): %q",
				message.Chat.ID, message.From.ID, err)
		}

		bot_app.handleAddSession(message)
		return true
	default:
		return false
	}
}

func (bot_app TelegramBotApp) _handleHelpMessage(message *tgbotapi.Message, isStart bool) {
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
/del - Removes an item from your shopping list
/clear - Removes all items from the shopping list`
	text := fmt.Sprintf(textTemplate, message.From.FirstName)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot_app.Bot.Send(msg)
}

func (bot_app TelegramBotApp) handleStart(message *tgbotapi.Message) {
	bot_app._handleHelpMessage(message, true)
}

func (bot_app TelegramBotApp) handleUnrecognisedMessage(message *tgbotapi.Message) {
	log.Print("No supported bot commands found")

	if len(message.Text) > 0 {
		// Display help text only if we received
		// a message text: we don't want to reply
		// to service messages (people added or removed from the group, etc.)
		bot_app._handleHelpMessage(message, false)
	}
}

func (bot_app TelegramBotApp) handleAdd(message *tgbotapi.Message) {
	itemName := message.CommandArguments()

	if itemName != "" {
		// If item name is supplied - just add it
		bot_app.handleAddSession(message)
		return
	}

	// If an item name is not provided in arguments,
	// allow the user to add an item following the two-step process
	command := models.CommandAddShoppingItem
	err := bot_app.Storage.AddUnfinishedCommand(models.UnfinishedCommand{
		Command:   command,
		ChatID:    message.Chat.ID,
		CreatedBy: message.From.ID,
	})

	if err != nil {
		log.Printf(
			"Unable to create an unfinished comamnd (%q): %q",
			command, err)

		// TODO: send "Something went wrong" to an user
		return
	}

	format := "Ok [%s](tg://user?id=%d), what do you want to add into your shopping list?"
	text := fmt.Sprintf(format, message.From.FirstName, message.From.ID)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	// If we are not in a private chat,
	// force user to reply because we can't listen all messages in a group.
	// See: https://core.telegram.org/bots/#privacy-mode
	if !message.Chat.IsPrivate() {
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
	}

	bot_app.Bot.Send(msg)
}

func (bot_app TelegramBotApp) handleList(message *tgbotapi.Message) {
	var text string
	chatID := message.Chat.ID

	chatItems, err := bot_app.Storage.GetShoppingItems(chatID)
	if err != nil {
		if err != nil {
			log.Printf(
				"Unable to get all shopping items (ChatID=%d): %q",
				chatID, err)

			// TODO: send "Something went wrong" to an user
			return
		}
	}

	if len(chatItems) == 0 {
		text = "Your shopping list is empty. Who knows, maybe it's a good thing"
	} else {
		offset := len(strconv.Itoa(len(chatItems)))
		listItemFormat := fmt.Sprintf("%%%dd. %%s\n", offset)

		listNumber := 1
		for _, item := range chatItems {
			text += fmt.Sprintf(listItemFormat, listNumber, item.Name)
			listNumber++
		}

		text = fmt.Sprintf(
			"%s\n\n```\n%s```",
			"Here is the list item in your shopping list:",
			text)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot_app.Bot.Send(msg)
}

func (bot_app TelegramBotApp) handleAddSession(message *tgbotapi.Message) {
	itemName := message.CommandArguments()
	if itemName == "" {
		itemName = message.Text
	}

	err := bot_app.Storage.AddShoppingItemIntoShoppingList(models.ShoppingItem{
		Name:      itemName,
		ChatID:    message.Chat.ID,
		CreatedBy: message.From.ID})
	if err != nil {
		log.Printf("Unable to add a new shopping item (ItemName=%s, ChatID=%d, UserId=%d): %q",
			itemName, message.Chat.ID, message.From.ID, err)

		// TODO: send "Something went wrong" to an user
		return
	}

	text := "Lovely! I've added \"%s\" into your shopping list. Anything else?"
	text = fmt.Sprintf(text, itemName)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	bot_app.Bot.Send(msg)
}

func (bot_app TelegramBotApp) handleDel(message *tgbotapi.Message) {
	var text string
	var itemButtonRows [][]tgbotapi.InlineKeyboardButton

	chatID := message.Chat.ID
	chatItems, err := bot_app.Storage.GetShoppingItems(chatID)
	if err != nil {
		log.Printf(
			"Unable to get all shopping items (ChatID=%d): %q",
			chatID, err)

		// TODO: send "Something went wrong" to an user
		return
	}

	isEmpty := len(chatItems) == 0
	if isEmpty {
		text = "Your shopping list is empty. No need to delete items 🙂"
	} else {
		for _, item := range chatItems {
			callbackData := fmt.Sprintf("del:%s", strconv.FormatInt(item.ID, 10))
			itemButtonRow := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item.Name, callbackData))
			itemButtonRows = append(itemButtonRows, itemButtonRow)
		}

		text = "Ok, what item do you want to delete from your shopping list?"
	}

	msg := tgbotapi.NewMessage(chatID, text)
	if !isEmpty {
		msg.BaseChat.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			itemButtonRows...)
	}
	bot_app.Bot.Send(msg)
}

func (bot_app TelegramBotApp) handleDelCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, data string) {
	bot_app.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(
		callbackQuery.ID, ""))

	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID
	itemID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return
	}

	var text string
	item, err := bot_app.Storage.GetShoppingItem(itemID)
	if err != nil {
		log.Printf("Unable to get a shopping item (ItemID=%d): %q",
			itemID, err)

		// TODO: send "Something went wrong" to an user
		return
	}

	if item != nil {
		err := bot_app.Storage.DeleteShoppingItem(itemID)
		if err != nil {
			log.Printf("Unable to delete a shopping item (ItemID=%d): %q",
				itemID, err)

			// TODO: send "Something went wrong" to an user
			return
		}

		text = "It's nice to see that you think that you don't "
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
		bot_app.Bot.Send(msg)
	}

	// Send deletion confimration text
	{
		msg := tgbotapi.NewMessage(chatID, text)
		bot_app.Bot.Send(msg)
	}
}

func (bot_app TelegramBotApp) handleClear(message *tgbotapi.Message) {
	var text string

	chatID := message.Chat.ID

	chatItems, err := bot_app.Storage.GetShoppingItems(chatID)
	if err != nil {
		log.Printf(
			"Unable to get all shopping items (ChatID=%d): %q",
			chatID, err)

		// TODO: send "Something went wrong" to an user
		return
	}

	isEmpty := len(chatItems) == 0
	if isEmpty {
		text = "Your shopping list is empty. No need to delete items 🙂"
	} else {
		text = "Are you sure that you want to *remove all items* from you shopping list?"
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if !isEmpty {
		msg.BaseChat.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Yes", "clear:1"),
				tgbotapi.NewInlineKeyboardButtonData("Cancel", "clear:0")})
	}
	bot_app.Bot.Send(msg)
}

func (bot_app TelegramBotApp) handleClearCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, data string) {
	bot_app.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(
		callbackQuery.ID, ""))

	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID
	confirmed, err := strconv.ParseBool(data)
	if err != nil {
		return
	}

	var text string
	if confirmed {
		text = "Ok, I've deleted all items from you shopping list.\n\nNow you can start from scratch, if you wish."

		err := bot_app.Storage.DeleteAllShoppingItems(chatID)
		if err != nil {
			log.Printf(
				"Unable to delete all shopping items (ChatID=%d): %q",
				chatID, err)

			// TODO: send "Something went wrong" to an user
			return
		}
	} else {
		text = "Canceling. Your items are still in your list."
	}

	// Edit previous message to hide the keyboard
	{
		msg := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			messageID,
			tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{}))
		bot_app.Bot.Send(msg)
	}

	// Send deletion confimration text
	{
		msg := tgbotapi.NewMessage(chatID, text)
		bot_app.Bot.Send(msg)
	}
}
