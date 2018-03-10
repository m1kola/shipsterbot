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
	"github.com/m1kola/shipsterbot/models"
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

	t.Run("Success", func(t *testing.T) {
		t.Run("Private chat", func(t *testing.T) {
			messageMock := &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID:   123,
					Type: "private",
				},
				From: &tgbotapi.User{
					ID:        321,
					FirstName: "m1kola",
				},
			}

			stMock.EXPECT().AddUnfinishedCommand(gomock.Any()).Return(nil)
			clientMock.EXPECT().Send(gomockhelpers.MatcherFunc(func(x interface{}) bool {
				msgCfg, ok := x.(tgbotapi.MessageConfig)
				if !ok {
					return false
				}

				// Reply to the same chat
				if msgCfg.ChatID != messageMock.Chat.ID {
					return false
				}

				// If replyMarkup is present, check that
				// bot is NOT forcing a client to reply
				if replyMarkup, ok := msgCfg.ReplyMarkup.(tgbotapi.ForceReply); ok {
					if !replyMarkup.ForceReply || !replyMarkup.Selective {
						return false
					}
				}

				return true
			}))

			err := handleAdd(clientMock, stMock, messageMock)

			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})

		t.Run("Group chat", func(t *testing.T) {
			messageMock := &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID:   123,
					Type: "group",
				},
				From: &tgbotapi.User{
					ID:        321,
					FirstName: "m1kola",
				},
			}

			stMock.EXPECT().AddUnfinishedCommand(gomock.Any()).Return(nil)
			clientMock.EXPECT().Send(gomockhelpers.MatcherFunc(func(x interface{}) bool {
				msgCfg, ok := x.(tgbotapi.MessageConfig)
				if !ok {
					return false
				}

				// Reply to the same chat
				if msgCfg.ChatID != messageMock.Chat.ID {
					return false
				}

				// Check that bot is forcing a client to reply
				replyMarkup, ok := msgCfg.ReplyMarkup.(tgbotapi.ForceReply)
				if !ok {
					return false
				}
				if !replyMarkup.ForceReply || !replyMarkup.Selective {
					return false
				}

				return true
			}))

			err := handleAdd(clientMock, stMock, messageMock)

			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})
	})
}

func TestHandleList(t *testing.T) {
	// Common interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMocksender(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	// Common data mocks
	errMock := errors.New("fake error")
	messageMock := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 123,
		},
		From: &tgbotapi.User{
			ID: 321,
		},
	}

	t.Run("Storage error", func(t *testing.T) {
		stMock.EXPECT().GetShoppingItems(messageMock.Chat.ID).Return(nil, errMock)

		err := handleList(clientMock, stMock, messageMock)

		if !strings.Contains(err.Error(), errMock.Error()) {
			t.Errorf("Expected err %#v, got %#v", errMock, err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		t.Run("Empty shopping list", func(t *testing.T) {
			// Data mocks
			storageDataMock := []*models.ShoppingItem{}

			stMock.EXPECT().GetShoppingItems(gomock.Any()).Return(storageDataMock, nil)
			clientMock.EXPECT().Send(gomockhelpers.MatcherFunc(func(x interface{}) bool {
				msgCfg, ok := x.(tgbotapi.MessageConfig)
				if !ok {
					return false
				}

				// Reply to the same chat
				if msgCfg.ChatID != messageMock.Chat.ID {
					return false
				}

				if !strings.Contains(msgCfg.Text, "list is empty") {
					return false
				}

				return true
			}))

			err := handleList(clientMock, stMock, messageMock)

			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})

		t.Run("Shopping list with items", func(t *testing.T) {
			// Data mocks
			storageDataMock := []*models.ShoppingItem{
				&models.ShoppingItem{
					Name: "Milk",
				},
				&models.ShoppingItem{
					Name: "Молоко",
				},
			}

			stMock.EXPECT().GetShoppingItems(gomock.Any()).Return(storageDataMock, nil)
			clientMock.EXPECT().Send(gomockhelpers.MatcherFunc(func(x interface{}) bool {
				msgCfg, ok := x.(tgbotapi.MessageConfig)
				if !ok {
					return false
				}

				// Reply to the same chat
				if msgCfg.ChatID != messageMock.Chat.ID {
					return false
				}

				for _, dataItem := range storageDataMock {
					if !strings.Contains(msgCfg.Text, dataItem.Name) {
						return false
					}
				}

				return true
			}))

			err := handleList(clientMock, stMock, messageMock)

			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})
	})
}
