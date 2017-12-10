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
			bot_app.handleCallbackQuery(&update)
		}

		if update.Message != nil {
			bot_app.handleMessage(update.Message)
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
		go bot_app.handleDelCallbackQuery(update.CallbackQuery, dataPieces[1])
	case "clear":
		go bot_app.handleClearCallbackQuery(update.CallbackQuery, dataPieces[1])
	}
}

func (bot_app TelegramBotApp) handleMessage(message *tgbotapi.Message) {
	log.Printf("Message received: %s", message.Text)

	isHandled := bot_app.handleMessageEntities(message)
	if !isHandled {
		isHandled = bot_app.handleMessageText(message)

		if !isHandled {
			log.Print("No supported bot commands found")
			go bot_app.handleUnrecognisedMessage(message)
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
		go bot_app.handleStart(message)
		return true
	case "add":
		go bot_app.handleAdd(message)
		return true
	case "list":
		go bot_app.handleList(message)
		return true
	case "del":
		go bot_app.handleDel(message)
		return true
	case "clear":
		go bot_app.handleClear(message)
		return true
	}

	return false
}

// handleMessageText handles unfinshed operations.
// Normally we listen to user's text commands or inline keyboard,
// but in some cases we need to handle message text.
// For example, when user asks us to add an item into the shopping list.
func (bot_app TelegramBotApp) handleMessageText(message *tgbotapi.Message) bool {
	session, ok := bot_app.Storage.GetUnfinishedCommand(message.Chat.ID,
		message.From.ID)

	if ok {
		switch session.Command {
		case models.CommandAddShoppingItem:
			bot_app.Storage.DeleteUnfinishedCommand(message.Chat.ID,
				message.From.ID)
			go bot_app.handleAddSession(message)
			return true
		}
	}

	return false
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
	bot_app._handleHelpMessage(message, false)
}

func (bot_app TelegramBotApp) handleAdd(message *tgbotapi.Message) {
	itemName := message.CommandArguments()

	if itemName == "" && message.Chat.IsPrivate() {
		// If item is not provided and we are in a private chat,
		// allow the user to add an item in two steps
		bot_app.Storage.AddUnfinishedCommand(models.UnfinishedCommand{
			Command:   models.CommandAddShoppingItem,
			ChatID:    message.Chat.ID,
			CreatedBy: message.From.ID,
		})

		text := "Ok, what do you want to add into your shopping list?"
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		bot_app.Bot.Send(msg)
	} else {
		// If we are in a group chat, the user must provide
		// an item name in the as an argument, because of group security policy
		// See: https://core.telegram.org/bots/#privacy-mode
		// TODO: Discover https://core.telegram.org/bots/api#forcereply
		//       Probably it's possible to improve UX
		if itemName == "" {
			// If item name is not supplied, give the user a clue
			format := "I'm so sorry, but in a group chat you have to"
			format += "specify an item you want to add using an argument. "
			format += "Try `/%s milk`"
			text := fmt.Sprintf(format, message.CommandWithAt())

			msg := tgbotapi.NewMessage(message.Chat.ID, text)
			msg.ParseMode = tgbotapi.ModeMarkdown
			bot_app.Bot.Send(msg)
		} else {
			// If item name is supplied - just add it
			bot_app.handleAddSession(message)
		}
	}
}

func (bot_app TelegramBotApp) handleList(message *tgbotapi.Message) {
	var text string
	chatID := message.Chat.ID

	chatItems, ok := bot_app.Storage.GetShoppingItems(chatID)
	if !ok || len(chatItems) == 0 {
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

	bot_app.Storage.AddShoppingItemIntoShoppingList(models.ShoppingItem{
		Name:      itemName,
		ChatID:    message.Chat.ID,
		CreatedBy: message.From.ID})

	text := "Lovely! I've added \"%s\" into your shopping list. Anything else?"
	text = fmt.Sprintf(text, itemName)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	bot_app.Bot.Send(msg)
}

func (bot_app TelegramBotApp) handleDel(message *tgbotapi.Message) {
	var text string
	var itemButtons []tgbotapi.InlineKeyboardButton

	chatID := message.Chat.ID
	chatItems, ok := bot_app.Storage.GetShoppingItems(chatID)
	isEmpty := !ok || len(chatItems) == 0
	if isEmpty {
		text = "Your shopping list is empty. No need to delete items 🙂"
	} else {
		for _, item := range chatItems {
			callbackData := fmt.Sprintf("del:%s", strconv.FormatInt(item.ID, 10))
			itemButton := tgbotapi.NewInlineKeyboardButtonData(item.Name, callbackData)
			itemButtons = append(itemButtons, itemButton)
		}

		text = "Ok, what item do you want to delete from your shopping list?"
	}

	msg := tgbotapi.NewMessage(chatID, text)
	if !isEmpty {
		msg.BaseChat.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(itemButtons)
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
	item, ok := bot_app.Storage.GetShoppingItem(itemID)
	if ok {
		bot_app.Storage.DeleteShoppingItem(itemID)

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

	chatItems, ok := bot_app.Storage.GetShoppingItems(chatID)
	isEmpty := !ok || len(chatItems) == 0
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

		bot_app.Storage.DeleteAllShoppingItems(chatID)
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
