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
