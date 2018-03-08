package telegram

import (
	"errors"
	"strings"
	"testing"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/golang/mock/gomock"
	"github.com/m1kola/shipsterbot/gomockhelpers"
	"github.com/m1kola/shipsterbot/mocks/bot/mock_telegram"
	"github.com/m1kola/shipsterbot/mocks/mock_storage"
	"github.com/m1kola/shipsterbot/storage"
)

func TestHelpMessages(t *testing.T) {
	// Common interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMocksender(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	messageMock := &tgbotapi.Message{
		From: &tgbotapi.User{
			FirstName: "m1kola",
		},
		Chat: &tgbotapi.Chat{
			ID: 123,
		},
		Text: "Some text",
	}

	greetingMatcher := gomockhelpers.MatcherFunc(func(x interface{}) bool {
		msgCfg, ok := x.(tgbotapi.MessageConfig)
		if !ok {
			return false
		}

		// Reply to the same chat
		if msgCfg.ChatID != messageMock.Chat.ID {
			return false
		}

		if !strings.Contains(msgCfg.Text, "Hi m1kola") {
			return false
		}

		return true
	})

	unknownCommandMatcher := gomockhelpers.MatcherFunc(func(x interface{}) bool {
		msgCfg, ok := x.(tgbotapi.MessageConfig)
		if !ok {
			return false
		}

		// Reply to the same chat
		if msgCfg.ChatID != messageMock.Chat.ID {
			return false
		}

		if !strings.Contains(msgCfg.Text, "m1kola, I'm very sorry") {
			return false
		}

		return true
	})

	t.Run("sendHelpMessage", func(t *testing.T) {
		t.Run("Greeting", func(t *testing.T) {
			clientMock.EXPECT().Send(greetingMatcher)

			sendHelpMessage(clientMock, messageMock, true)
		})

		t.Run("Unknown command", func(t *testing.T) {
			clientMock.EXPECT().Send(unknownCommandMatcher)

			sendHelpMessage(clientMock, messageMock, false)
		})
	})

	t.Run("handleStart", func(t *testing.T) {
		t.Run("Greeting", func(t *testing.T) {
			clientMock.EXPECT().Send(greetingMatcher)

			handleStart(clientMock, stMock, messageMock)
		})
	})
}

func TestHandleUnrecoverableError(t *testing.T) {
	var expectedChatID int64 = 123

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)

	clientMock.EXPECT().Send(gomockhelpers.MatcherFunc(func(x interface{}) bool {
		msgCfg, ok := x.(tgbotapi.MessageConfig)
		if !ok {
			return false
		}

		// Reply to the same chat
		if msgCfg.ChatID != expectedChatID {
			return false
		}

		if !strings.Contains(msgCfg.Text, "Please, try again a bit later") {
			return false
		}

		return true
	}))

	handleUnrecoverableError(clientMock, expectedChatID, errors.New("fake err"))
}

func TestHandleAdd(t *testing.T) {
	// Common interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMocksender(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	// Common data mocks
	errMock := errors.New("fake error")

	t.Run("Handle command arguments", func(t *testing.T) {
		messageMock := mock_telegram.MessageCommandMockSetup(commandAdd, "some item")

		handleAddSessionOld := handleAddSession
		defer func() { handleAddSession = handleAddSessionOld }()
		handleAddSession = func(
			_ sender,
			_ storage.DataStorageInterface,
			message *tgbotapi.Message,
		) error {
			if message != messageMock {
				t.Error("Unexpected message received")
			}

			return errMock
		}

		err := handleAdd(clientMock, stMock, messageMock)

		if errMock != err {
			t.Errorf("Expected err %#v, got %#v", errMock, err)
		}
	})

	t.Run("Storage error", func(t *testing.T) {
		messageMock := &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			From: &tgbotapi.User{
				ID: 321,
			},
		}

		stMock.EXPECT().AddUnfinishedCommand(gomock.Any()).Return(errMock)

		err := handleAdd(clientMock, stMock, messageMock)

		if !strings.Contains(err.Error(), errMock.Error()) {
			t.Errorf("Expected err %#v, got %#v", errMock, err)
		}
	})
}
