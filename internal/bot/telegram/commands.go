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
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/m1kola/shipsterbot/internal/pkg/storage"
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
	commandHandler           commandHandlerFunc
	unfinishedCommandHandler commandHandlerFunc
	callbackQueryHandler     callbackQueryHandlerFunc
}

var getBotCommandsMapping = func() map[string]botCommand {
	return map[string]botCommand{
		// The `/start` command is implicit: Telegram sends on user's behalf
		// when they start the bot.
		commandStart: {
			commandHandler: handleStart,
		},
		commandHelp: {
			description:       "Show the list of available commands and short descriptions",
			showInHelpMessage: true,
			commandHandler:    handleStart,
		},
		commandAdd: {
			description:              "Add an item into your shopping list",
			showInHelpMessage:        true,
			commandHandler:           handleAdd,
			unfinishedCommandHandler: handleAddSession,
		},
		commandList: {
			description:       "Display items the shopping list",
			showInHelpMessage: true,
			commandHandler:    handleList,
		},
		commandDel: {
			description:          "Delete an item from your shopping list",
			showInHelpMessage:    true,
			commandHandler:       handleDel,
			callbackQueryHandler: handleDelCallbackQuery,
		},
		commandClear: {
			description:          "Delete all items from the shopping list",
			showInHelpMessage:    true,
			commandHandler:       handleClear,
			callbackQueryHandler: handleClearCallbackQuery,
		},
	}
}

// commandHandlerFunc defines required signature for a command handler func
type commandHandlerFunc func(
	client sender, st storage.DataStorageInterface, message *tgbotapi.Message,
) error

// callbackQueryHandlerFunc defines required signature for a callback query func
type callbackQueryHandlerFunc func(
	client botClientInterface,
	st storage.DataStorageInterface,
	callbackQuery *tgbotapi.CallbackQuery,
	data string,
) error
