package telegram

import (
	"testing"
)

func TestSplitCallbackQueryData(t *testing.T) {
	t.Run("No error", func(t *testing.T) {
		expectedCommand := "command_name"
		expectedPayload := "123"
		command, payload, err := splitCallbackQueryData("command_name:123")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if command != expectedCommand {
			t.Fatalf("Expected %#v, got %#v", expectedCommand, command)
		}

		if payload != "123" {
			t.Fatalf("Expected %#v, got %#v", expectedPayload, payload)
		}
	})

	t.Run("Error", func(t *testing.T) {
		testCases := []string{
			"",
			":",
			"The",
			"quick:",
			":brown",
			"fox:jumps:",
			":over:the",
			"lazy:dog:lorem:ipsum",
		}

		for _, test := range testCases {
			_, _, err := splitCallbackQueryData(test)

			if err == nil {
				t.Errorf("Expected error for the test case %#v", test)
			}
		}
	})
}

func TestJoinCallbackQueryData(t *testing.T) {
	expectedData := "one:two"
	data := joinCallbackQueryData("one", "two")

	if expectedData != data {
		t.Fatalf("Expected %#v, got %#v", expectedData, data)
	}
}
