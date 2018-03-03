package telegram

import (
	"errors"
	"fmt"
	"log"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/m1kola/shipsterbot/models"
	"github.com/m1kola/shipsterbot/storage"
)

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
