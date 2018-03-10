package telegram

import (
	"fmt"
	"strconv"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/m1kola/shipsterbot/models"
	"github.com/m1kola/shipsterbot/storage"
)

var sendHelpMessage = func(client sender, message *tgbotapi.Message, isStart bool) {
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
	client.Send(msg)
}

func handleStart(
	client sender,
	st storage.DataStorageInterface,
	message *tgbotapi.Message,
) error {
	sendHelpMessage(client, message, true)
	return nil
}

// handleUnrecoverableError sends the "something went wrong" message to a chat
//
// TODO: We should actually have notifications about errors at some point
// 		 See: https://github.com/m1kola/shipsterbot/issues/27
var handleUnrecoverableError = func(
	client botClientInterface,
	chatID int64,
	_ error,
) {
	text := "Sorry, but something went wrong. I'll inform developers about this issue. Please, try again a bit later."
	msg := tgbotapi.NewMessage(chatID, text)
	client.Send(msg)
}

func handleAdd(
	client sender,
	st storage.DataStorageInterface,
	message *tgbotapi.Message,
) error {
	itemName := message.CommandArguments()
	if itemName != "" {
		// If item name is supplied - just add it
		return handleAddSession(client, st, message)
	}

	// If an item name is not provided in arguments,
	// allow the user to add an item following the two-step process
	err := st.AddUnfinishedCommand(models.UnfinishedCommand{
		Command:   commandAdd,
		ChatID:    message.Chat.ID,
		CreatedBy: message.From.ID,
	})

	if err != nil {
		return fmt.Errorf(
			"Unable to create an unfinished comamnd (%v): %v",
			commandAdd, err)
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

	client.Send(msg)

	return nil
}

func handleList(
	client sender,
	st storage.DataStorageInterface,
	message *tgbotapi.Message,
) error {
	var text string
	chatID := message.Chat.ID

	chatItems, err := st.GetShoppingItems(chatID)
	if err != nil {
		return fmt.Errorf(
			"Unable to get all shopping items (ChatID=%d): %v",
			chatID, err)
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
	client.Send(msg)

	return nil
}

var handleAddSession = func(
	client sender,
	st storage.DataStorageInterface,
	message *tgbotapi.Message,
) error {
	itemName := message.CommandArguments()
	if itemName == "" {
		itemName = message.Text
	}

	err := st.AddShoppingItemIntoShoppingList(models.ShoppingItem{
		Name:      itemName,
		ChatID:    message.Chat.ID,
		CreatedBy: message.From.ID})
	if err != nil {
		return fmt.Errorf(
			"Unable to add a new shopping item (ItemName=%s, ChatID=%d, UserId=%d): %v",
			itemName, message.Chat.ID, message.From.ID, err)
	}

	text := "Lovely! I've added \"%s\" into your shopping list. Anything else?"
	text = fmt.Sprintf(text, itemName)
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	client.Send(msg)
	return nil
}

func handleDel(
	client sender,
	st storage.DataStorageInterface,
	message *tgbotapi.Message,
) error {
	var text string
	var itemButtonRows [][]tgbotapi.InlineKeyboardButton

	chatID := message.Chat.ID
	chatItems, err := st.GetShoppingItems(chatID)
	if err != nil {
		return fmt.Errorf(
			"Unable to get all shopping items (ChatID=%d): %v",
			chatID, err)
	}

	isEmpty := len(chatItems) == 0
	if isEmpty {
		text = "Your shopping list is empty. No need to delete items 🙂"
	} else {
		for _, item := range chatItems {
			callbackData := joinCallbackQueryData(
				commandDel, strconv.FormatInt(item.ID, 10),
			)
			itemButtonRow := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(item.Name, callbackData),
			)
			itemButtonRows = append(itemButtonRows, itemButtonRow)
		}

		text = "Ok, what item do you want to delete from your shopping list?"
	}

	msg := tgbotapi.NewMessage(chatID, text)
	if !isEmpty {
		msg.BaseChat.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			itemButtonRows...)
	}
	client.Send(msg)

	return nil
}

func handleDelCallbackQuery(
	client botClientInterface,
	st storage.DataStorageInterface,
	callbackQuery *tgbotapi.CallbackQuery,
	data string,
) error {
	client.AnswerCallbackQuery(tgbotapi.NewCallback(
		callbackQuery.ID, ""))

	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID
	itemID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		// User can't amend CallBackData, so most likely it's our fault
		return fmt.Errorf(
			"Unable to parse ItemID from the CallbackQuery data %s: %v",
			data, err)
	}

	item, err := st.GetShoppingItem(itemID)
	if err != nil {
		return fmt.Errorf(
			"Unable to get a shopping item (ItemID=%d): %v",
			itemID, err)
	}

	var text string
	if item != nil {
		err := st.DeleteShoppingItem(itemID)
		if err != nil {
			return fmt.Errorf(
				"Unable to delete a shopping item (ItemID=%d): %v",
				itemID, err)
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
		// It's important to send replyMarkup exactly like this,
		// because otherwise telegram clients will not hide the keybaord
		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{},
		)
		msg := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			messageID,
			replyMarkup,
		)
		client.Send(msg)
	}

	// Send deletion confimration text
	{
		msg := tgbotapi.NewMessage(chatID, text)
		client.Send(msg)
	}

	return nil
}

const clearCallbackDataConfim = "1"
const clearCallbackDataCancel = "0"

func handleClear(
	client sender,
	st storage.DataStorageInterface,
	message *tgbotapi.Message,
) error {
	var text string

	chatID := message.Chat.ID

	chatItems, err := st.GetShoppingItems(chatID)
	if err != nil {
		// User can't amend CallBackData, so most likely it's our fault
		return fmt.Errorf(
			"Unable to get all shopping items (ChatID=%d): %v",
			chatID, err)
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
		yesCallbackData := joinCallbackQueryData(commandClear, clearCallbackDataConfim)
		cancelCallbackData := joinCallbackQueryData(commandClear, clearCallbackDataCancel)

		msg.BaseChat.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				// Order of the buttons in keyboard is important
				tgbotapi.NewInlineKeyboardButtonData("Yes", yesCallbackData),
				tgbotapi.NewInlineKeyboardButtonData("Cancel", cancelCallbackData)})
	}
	client.Send(msg)
	return nil
}

func handleClearCallbackQuery(
	client botClientInterface,
	st storage.DataStorageInterface,
	callbackQuery *tgbotapi.CallbackQuery,
	data string,
) error {

	client.AnswerCallbackQuery(tgbotapi.NewCallback(
		callbackQuery.ID, ""))

	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID
	dataIsValid := data == clearCallbackDataConfim || data == clearCallbackDataCancel
	if !dataIsValid {
		return fmt.Errorf(
			"Unable to parse confirmation from the CallbackQuery data %#v",
			data,
		)
	}

	var text string
	confirmed := data == clearCallbackDataConfim
	if confirmed {
		text = "Ok, I've deleted all items from you shopping list.\n\nNow you can start from scratch, if you wish."

		err := st.DeleteAllShoppingItems(chatID)
		if err != nil {
			return fmt.Errorf(
				"Unable to delete all shopping items (ChatID=%d): %v",
				chatID, err)
		}
	} else {
		text = "Canceling. Your items are still in your list."
	}

	// Edit previous message to hide the keyboard
	{
		// It's important to send replyMarkup exactly like this,
		// because otherwise telegram clients will not hide the keybaord
		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{},
		)
		msg := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			messageID,
			replyMarkup,
		)
		client.Send(msg)
	}

	// Send deletion confimration text
	{
		msg := tgbotapi.NewMessage(chatID, text)
		client.Send(msg)
	}

	return nil
}
