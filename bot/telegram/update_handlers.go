package telegram

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/m1kola/shipsterbot/models"
	"github.com/m1kola/shipsterbot/storage"
)

var errCommandIsNotSupported = handlerCanNotHandleError{
	errors.New("Unable to find a handler for a command")}

// handleUpdates receives updates and starts goroutines to handle them
func handleUpdates(
	client botClientInterface,
	st storage.DataStorageInterface,
	updates <-chan tgbotapi.Update,
) {
	for update := range updates {
		go handleUpdate(client, st, update)
	}
}

var handleUpdate = func(
	client botClientInterface,
	st storage.DataStorageInterface,
	update tgbotapi.Update,
) {
	var err error
	var message *tgbotapi.Message

	if update.CallbackQuery != nil {
		message = update.CallbackQuery.Message

		err = handleCallbackQuery(client, st, update.CallbackQuery)
	} else if update.Message != nil {
		message = update.Message

		err = handleMessage(client, st, message)
	}

	if err != nil {
		log.Print(err)

		if _, ok := err.(handlerCanNotHandleError); ok {
			// It's ok if we can't handle a message,
			// because an user can send nonsense.
			// Let's send a message saying that
			// we don't understand the input.
			handleUnrecognisedMessage(client, message)
		} else {
			// Other types of error mean that we are in trouble
			// and we need to do something with it
			text := "Sorry, but something went wrong. I'll inform developers about this issue. Please, try again a bit later."
			msg := tgbotapi.NewMessage(message.Chat.ID, text)
			client.Send(msg)
		}
	}
}

// handleCallbackQuery handles user's interactions with the client's UI
// User can interaction with a bot using an inline keyboard, for example
var handleCallbackQuery = func(
	client botClientInterface,
	st storage.DataStorageInterface,
	callbackQuery *tgbotapi.CallbackQuery,
) error {
	if callbackQuery.Data == "" {
		return handlerCanNotHandleError{
			errors.New("Empty data in the CallbackQuery")}
	}

	dataPieces := strings.SplitN(callbackQuery.Data, ":", 2)
	if len(dataPieces) != 2 {
		return handlerCanNotHandleError{
			fmt.Errorf("Wrong data format in the CallbackQuery: %v",
				callbackQuery.Data)}
	}

	botCommand := dataPieces[0]
	switch botCommand {
	case "del":
		return handleDelCallbackQuery(client, st, callbackQuery, dataPieces[1])
	case "clear":
		return handleClearCallbackQuery(client, st, callbackQuery, dataPieces[1])
	}

	return handlerCanNotHandleError{
		fmt.Errorf("Unable to find a handler for CallbackQuery: %v",
			callbackQuery.Data)}
}

// handleMessage handles messages.
// Messages can contain entities in some cases (commands, mentions, etc),
// but can also be plain text messages.
var handleMessage = func(
	client sender,
	st storage.DataStorageInterface,
	message *tgbotapi.Message,
) error {
	log.Printf("Message received: \"%s\"", message.Text)

	err := handleMessageEntities(client, st, message)
	// We should only try to continue processing an message,
	// if we receive an handlerCanNotHandleError error.
	if _, ok := err.(handlerCanNotHandleError); ok {
		// But it doesn't make sense to continue if it's
		// an errCommandIsNotSupported error
		if err == errCommandIsNotSupported {
			return err
		}

		return handleMessageText(client, st, message)
	}

	return err
}

// handleMessageEntities handles entities form a message
func handleMessageEntities(
	client sender,
	st storage.DataStorageInterface,
	message *tgbotapi.Message,
) error {
	if message.Entities == nil {
		return handlerCanNotHandleError{
			errors.New("Message doesn't have entities to handle")}
	}

	botCommand := message.Command()
	switch botCommand {
	case "help", "start":
		return handleStart(client, message)
	case "add":
		return handleAdd(client, st, message)
	case "list":
		return handleList(client, st, message)
	case "del":
		return handleDel(client, st, message)
	case "clear":
		return handleClear(client, st, message)
	}

	return errCommandIsNotSupported
}

// handleMessageText handles text from a message.
// Normally we listen to user's text commands or inline keyboard,
// but in some cases we need to handle message text.
// For example, when user asks us to add an item into the shopping list.
func handleMessageText(
	client sender,
	st storage.DataStorageInterface,
	message *tgbotapi.Message,
) error {
	session, err := st.GetUnfinishedCommand(message.Chat.ID,
		message.From.ID)

	if err != nil {
		return fmt.Errorf(
			"Unable to get an unfinished comamnd (ChatID=%d and UserId=%d): %v",
			message.Chat.ID, message.From.ID, err)
	}

	if session == nil {
		// Unfinished command doesn't exist. It's ok,
		// but we need to return an error just to indicate that
		// we didn't manage to handele this message
		return handlerCanNotHandleError{
			fmt.Errorf(
				"Can't find unfinished commands (ChatID=%d and UserId=%d)",
				message.Chat.ID, message.From.ID)}
	}

	switch session.Command {
	case models.CommandAddShoppingItem:
		err := st.DeleteUnfinishedCommand(message.Chat.ID,
			message.From.ID)

		if err != nil {
			return fmt.Errorf(
				"Unable to delete an unfinished comamnd (ChatID=%d and UserId=%d): %v",
				message.Chat.ID, message.From.ID, err)
		}

		return handleAddSession(client, st, message)
	}

	return handlerCanNotHandleError{
		errors.New("Unable to find a handler for the message")}
}

func handleHelpMessage(client sender, message *tgbotapi.Message, isStart bool) {
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

func handleStart(client sender, message *tgbotapi.Message) error {
	handleHelpMessage(client, message, true)
	return nil
}

var handleUnrecognisedMessage = func(client sender, message *tgbotapi.Message) {
	if len(message.Text) > 0 {
		log.Print("No supported bot commands found")

		// Display help text only if we received
		// a message text: we don't want to reply
		// to service messages (people added or removed from the group, etc.)
		handleHelpMessage(client, message, false)
	}
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
	command := models.CommandAddShoppingItem
	err := st.AddUnfinishedCommand(models.UnfinishedCommand{
		Command:   command,
		ChatID:    message.Chat.ID,
		CreatedBy: message.From.ID,
	})

	if err != nil {
		return fmt.Errorf(
			"Unable to create an unfinished comamnd (%v): %v",
			command, err)
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

func handleAddSession(
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

	var text string
	item, err := st.GetShoppingItem(itemID)
	if err != nil {
		return fmt.Errorf(
			"Unable to get a shopping item (ItemID=%d): %v",
			itemID, err)
	}

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
		msg := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			messageID,
			tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{}))
		client.Send(msg)
	}

	// Send deletion confimration text
	{
		msg := tgbotapi.NewMessage(chatID, text)
		client.Send(msg)
	}

	return nil
}

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
		msg.BaseChat.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Yes", "clear:1"),
				tgbotapi.NewInlineKeyboardButtonData("Cancel", "clear:0")})
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
	confirmed, err := strconv.ParseBool(data)
	if err != nil {
		return fmt.Errorf(
			"Unable to parse confirmation from the CallbackQuery data %s: %v",
			data, err)
	}

	var text string
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
		msg := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			messageID,
			tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{}))
		client.Send(msg)
	}

	// Send deletion confimration text
	{
		msg := tgbotapi.NewMessage(chatID, text)
		client.Send(msg)
	}

	return nil
}
