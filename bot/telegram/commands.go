// This module defines command representaton and handler interfaces.
//
// There are three types of handlers:
//  * Command handlers - To handle bot commands like `/add`, `/del`, `/help`.
//    In the majority of cases commands just start interactions with users;
//
//  * Unfinished command handlers - handle multistep commands like
//    the `/add` command.
//
//    Command handler for the `/add` command sends a message
//    asking an user about items they wants to add into their shopping list.
//    User replies in a separate message that that is
//    handled by an unfinished command handler;
//
//  * Callback query handlers - To handle commands that use an inline keyboard.
//
//    For example, the `/del` command handler sends a message with an inline
//    keybaord that asks an users to confirm what item they want to delete
//    from their shopping list.

package telegram

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/m1kola/shipsterbot/storage"
)

// All possible bot commands
const (
	commandStart = "start"
	commandHelp  = "help"
	commandAdd   = "add"
	commandList  = "list"
	commandDel   = "del"
	commandClear = "clear"
)

// botCommand defines command and it's handlers
type botCommand struct {
	description              string
	showInHelpMessage        bool
	commandHandler           commandHandler
	unfinishedCommandHandler commandHandler
	callbackQueryHandler     callbackQueryHandler
}

var getBotCommandsMapping = func() map[string]botCommand {
	return map[string]botCommand{
		// The `/start` command is implicit: Telegram sends on user's behalf
		// when they start the bot.
		commandStart: botCommand{
			commandHandler: commandHandlerFunc(handleStart),
		},
		commandHelp: botCommand{
			description:       "Show the list of available commands and short descriptions",
			showInHelpMessage: true,
			commandHandler:    commandHandlerFunc(handleStart),
		},
		commandAdd: botCommand{
			description:              "Add an item into your shopping list",
			showInHelpMessage:        true,
			commandHandler:           commandHandlerFunc(handleAdd),
			unfinishedCommandHandler: commandHandlerFunc(handleAddSession),
		},
		commandList: botCommand{
			description:       "Display items the shopping list",
			showInHelpMessage: true,
			commandHandler:    commandHandlerFunc(handleList),
		},
		commandDel: botCommand{
			description:          "Delete an item from your shopping list",
			showInHelpMessage:    true,
			commandHandler:       commandHandlerFunc(handleDel),
			callbackQueryHandler: callbackQueryHandlerFunc(handleDelCallbackQuery),
		},
		commandClear: botCommand{
			description:          "Delete all items from the shopping list",
			showInHelpMessage:    true,
			commandHandler:       commandHandlerFunc(handleClear),
			callbackQueryHandler: callbackQueryHandlerFunc(handleClearCallbackQuery),
		},
	}
}

type commandHandler interface {
	HandleCommand(
		client sender,
		st storage.DataStorageInterface,
		message *tgbotapi.Message,
	) error
}

type unfinishedCommandHandler interface {
	HandleUnfinishedCommand(
		client sender,
		st storage.DataStorageInterface,
		message *tgbotapi.Message,
	) error
}

type callbackQueryHandler interface {
	HandleCallbackQuery(
		client botClientInterface,
		st storage.DataStorageInterface,
		callbackQuery *tgbotapi.CallbackQuery,
		data string,
	) error
}

// commandHandlerFunc makes a func to implement commandHandler
type commandHandlerFunc func(
	client sender, st storage.DataStorageInterface, message *tgbotapi.Message,
) error

func (f commandHandlerFunc) HandleCommand(
	client sender, st storage.DataStorageInterface, message *tgbotapi.Message,
) error {
	return f(client, st, message)
}

// callbackQueryHandlerFunc makes a func to implement callbackQueryHandler
type callbackQueryHandlerFunc func(
	client botClientInterface,
	st storage.DataStorageInterface,
	callbackQuery *tgbotapi.CallbackQuery,
	data string,
) error

func (f callbackQueryHandlerFunc) HandleCallbackQuery(
	client botClientInterface,
	st storage.DataStorageInterface,
	callbackQuery *tgbotapi.CallbackQuery,
	data string,
) error {
	return f(client, st, callbackQuery, data)
}
