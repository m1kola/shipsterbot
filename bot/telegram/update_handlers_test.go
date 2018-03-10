package telegram

import (
	"errors"
	"strings"
	"testing"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/golang/mock/gomock"
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
		From: &tgbotapi.User{FirstName: "m1kola"},
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "Some text",
	}

	generateGreetingChecker := func(t *testing.T) interface{} {
		return func(msgCfg tgbotapi.MessageConfig) {
			if msgCfg.ChatID != messageMock.Chat.ID {
				t.Errorf(
					"Expected to reply to the chat with ID %d, but reply sent to %d",
					msgCfg.ChatID,
					messageMock.Chat.ID,
				)
			}

			expectedText := "Hi m1kola"
			if !strings.Contains(msgCfg.Text, expectedText) {
				t.Fatalf("Expected message to contain %#v", expectedText)
			}
		}
	}

	generateUnknownCommandChecker := func(t *testing.T) interface{} {
		return func(msgCfg tgbotapi.MessageConfig) {
			if msgCfg.ChatID != messageMock.Chat.ID {
				t.Errorf(
					"Expected to reply to the chat with ID %d, but reply sent to %d",
					msgCfg.ChatID,
					messageMock.Chat.ID,
				)
			}

			expectedText := "m1kola, I'm very sorry"
			if !strings.Contains(msgCfg.Text, expectedText) {
				t.Fatalf("Expected message to contain %#v", expectedText)
			}
		}
	}

	t.Run("sendHelpMessage", func(t *testing.T) {
		t.Run("Greeting", func(t *testing.T) {
			clientMock.EXPECT().Send(gomock.Any()).Do(
				generateGreetingChecker(t),
			)

			sendHelpMessage(clientMock, messageMock, true)
		})

		t.Run("Unknown command", func(t *testing.T) {
			clientMock.EXPECT().Send(gomock.Any()).Do(
				generateUnknownCommandChecker(t),
			)

			sendHelpMessage(clientMock, messageMock, false)
		})
	})

	t.Run("handleStart", func(t *testing.T) {
		t.Run("Greeting", func(t *testing.T) {
			clientMock.EXPECT().Send(gomock.Any()).Do(
				generateGreetingChecker(t),
			)

			handleStart(clientMock, stMock, messageMock)
		})
	})
}

func TestHandleUnrecoverableError(t *testing.T) {
	var expectedChatID int64 = 123

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)

	clientMock.EXPECT().Send(gomock.Any()).Do(func(msgCfg tgbotapi.MessageConfig) {
		if msgCfg.ChatID != expectedChatID {
			t.Errorf(
				"Expected to reply to the chat with ID %d, but reply sent to %d",
				msgCfg.ChatID,
				expectedChatID,
			)
		}

		expectedText := "Please, try again a bit later"
		if !strings.Contains(msgCfg.Text, expectedText) {
			t.Fatalf("Expected message to contain %#v", expectedText)
		}
	})

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
			Chat: &tgbotapi.Chat{ID: 123},
			From: &tgbotapi.User{ID: 321},
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
				Chat: &tgbotapi.Chat{ID: 123, Type: "private"},
				From: &tgbotapi.User{ID: 321, FirstName: "m1kola"},
			}

			stMock.EXPECT().AddUnfinishedCommand(gomock.Any()).Return(nil)
			clientMock.EXPECT().Send(gomock.Any()).Do(func(msgCfg tgbotapi.MessageConfig) {
				if msgCfg.ChatID != messageMock.Chat.ID {
					t.Errorf(
						"Expected to reply to the chat with ID %d, but reply sent to %d",
						msgCfg.ChatID,
						messageMock.Chat.ID,
					)
				}

				replyMarkup, ok := msgCfg.ReplyMarkup.(tgbotapi.ForceReply)
				if ok && replyMarkup.ForceReply && replyMarkup.Selective {
					t.Error("Expected bot to not force clients to reply in a private chat")
				}
			})

			err := handleAdd(clientMock, stMock, messageMock)

			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})

		t.Run("Group chat", func(t *testing.T) {
			messageMock := &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123, Type: "group"},
				From: &tgbotapi.User{ID: 321, FirstName: "m1kola"},
			}

			stMock.EXPECT().AddUnfinishedCommand(gomock.Any()).Return(nil)
			clientMock.EXPECT().Send(gomock.Any()).Do(func(msgCfg tgbotapi.MessageConfig) {
				if msgCfg.ChatID != messageMock.Chat.ID {
					t.Errorf(
						"Expected to reply to the chat with ID %d, but reply sent to %d",
						msgCfg.ChatID,
						messageMock.Chat.ID,
					)
				}

				replyMarkup, ok := msgCfg.ReplyMarkup.(tgbotapi.ForceReply)
				if !ok || !replyMarkup.ForceReply || !replyMarkup.Selective {
					t.Error("Expected bot to force clients to reply in a group chat")
				}
			})

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
		Chat: &tgbotapi.Chat{ID: 123},
		From: &tgbotapi.User{ID: 321},
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
			clientMock.EXPECT().Send(gomock.Any()).Do(func(msgCfg tgbotapi.MessageConfig) {
				if msgCfg.ChatID != messageMock.Chat.ID {
					t.Errorf(
						"Expected to reply to the chat with ID %d, but reply sent to %d",
						msgCfg.ChatID,
						messageMock.Chat.ID,
					)
				}

				expectedText := "list is empty"
				if !strings.Contains(msgCfg.Text, expectedText) {
					t.Fatalf("Expected message to contain %#v", expectedText)
				}
			})

			err := handleList(clientMock, stMock, messageMock)

			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})

		t.Run("Shopping list with items", func(t *testing.T) {
			// Data mocks
			storageDataMock := []*models.ShoppingItem{
				&models.ShoppingItem{Name: "Milk"},
				&models.ShoppingItem{Name: "Молоко"},
			}

			stMock.EXPECT().GetShoppingItems(gomock.Any()).Return(storageDataMock, nil)
			clientMock.EXPECT().Send(gomock.Any()).Do(func(msgCfg tgbotapi.MessageConfig) {
				if msgCfg.ChatID != messageMock.Chat.ID {
					t.Errorf(
						"Expected to reply to the chat with ID %d, but reply sent to %d",
						msgCfg.ChatID,
						messageMock.Chat.ID,
					)
				}

				// Chec if all items are present in the message
				for _, dataItem := range storageDataMock {
					if !strings.Contains(msgCfg.Text, dataItem.Name) {
						t.Errorf("Expected message to contain %#v", dataItem.Name)
					}
				}
			})

			err := handleList(clientMock, stMock, messageMock)

			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})
	})
}

func TestHandleAddSession(t *testing.T) {
	// Common interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMocksender(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	t.Run("Storage error", func(t *testing.T) {
		// Data mocks
		errMock := errors.New("fake error")
		messageMock := &tgbotapi.Message{
			Text: "Milk",
			Chat: &tgbotapi.Chat{ID: 123},
			From: &tgbotapi.User{ID: 321},
		}

		stMock.EXPECT().AddShoppingItemIntoShoppingList(gomock.Any()).Return(errMock)

		err := handleAddSession(clientMock, stMock, messageMock)

		if !strings.Contains(err.Error(), errMock.Error()) {
			t.Errorf("Expected err %#v, got %#v", errMock, err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		// Test cases
		expectedItemName := "milk"
		testCases := []*tgbotapi.Message{
			// The commandAdd with an item name as a command argument
			mock_telegram.MessageCommandMockSetup(commandAdd, expectedItemName),

			// Plain text message
			&tgbotapi.Message{Text: expectedItemName},
		}

		for _, messageMock := range testCases {
			// Set up common message mock fields
			messageMock.Chat = &tgbotapi.Chat{ID: 123}
			messageMock.From = &tgbotapi.User{ID: 321}

			// Set up interface mocks
			stMock.EXPECT().AddShoppingItemIntoShoppingList(
				gomock.Any(),
			).Do(func(item models.ShoppingItem) {
				if item.Name != expectedItemName {
					t.Errorf(
						"Expected item with name %#v, got %#v",
						expectedItemName, item.Name,
					)
				}
				if item.ChatID != messageMock.Chat.ID {
					t.Errorf(
						"Expected ChatID to be %d, got %d",
						messageMock.Chat.ID, item.ChatID,
					)
				}
			}).Return(nil)

			clientMock.EXPECT().Send(
				gomock.Any(),
			).Do(func(msgCfg tgbotapi.MessageConfig) {
				if msgCfg.ChatID != messageMock.Chat.ID {
					t.Errorf(
						"Expected to reply to the chat with ID %d, but reply sent to %d",
						msgCfg.ChatID,
						messageMock.Chat.ID,
					)
				}

				expectedText := expectedItemName
				if !strings.Contains(msgCfg.Text, expectedText) {
					t.Fatalf("Expected message to contain %#v", expectedText)
				}
			})

			err := handleAddSession(clientMock, stMock, messageMock)
			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		}

	})
}
