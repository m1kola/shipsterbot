package telegram

import (
	"testing"
)

// Stupid command mapping config test
func TestGetBotCommandsMapping(t *testing.T) {
	allExpectedCommands := []string{
		commandStart, commandHelp, commandAdd, commandList,
		commandDel, commandClear,
	}
	commandsWithUnfinishedCommandHandler := []string{commandAdd}
	commandsWithCallbackQueryHandler := []string{commandDel, commandClear}

	mapping := getBotCommandsMapping()

	t.Run("Test basic config", func(t *testing.T) {
		actualLen := len(mapping)
		expectedLen := len(allExpectedCommands)
		if expectedLen < actualLen {
			t.Errorf("There are more commands then expected. Expected %d, got %d",
				expectedLen, actualLen)
		}

		for _, expectedCommand := range allExpectedCommands {
			item, ok := mapping[expectedCommand]
			if !ok {
				t.Errorf("Command %#v wasn't found in the mapping", expectedCommand)
			}

			if item.commandHandler == nil {
				t.Errorf("Command %#v has no commandHandler", expectedCommand)
			}
		}
	})

	t.Run("Commands with unfinishedCommandHandler", func(t *testing.T) {
		for _, expectedCommand := range commandsWithUnfinishedCommandHandler {
			item := mapping[expectedCommand]

			if item.unfinishedCommandHandler == nil {
				t.Errorf("Command %#v has no unfinishedCommandHandler", expectedCommand)
			}
		}
	})

	t.Run("Commands with callbackQueryHandler", func(t *testing.T) {
		for _, expectedCommand := range commandsWithCallbackQueryHandler {
			item := mapping[expectedCommand]

			if item.callbackQueryHandler == nil {
				t.Errorf("Command %#v has no callbackQueryHandler", expectedCommand)
			}
		}
	})
}
