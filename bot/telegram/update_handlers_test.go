package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/golang/mock/gomock"
	"github.com/m1kola/shipsterbot/mocks/bot/mock_telegram"
	"github.com/m1kola/shipsterbot/mocks/mock_storage"
	"github.com/m1kola/shipsterbot/models"
	"github.com/m1kola/shipsterbot/internal/pkg/storage"
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

func TestHandleDel(t *testing.T) {
	// Common interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMocksender(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	// Common data mocks
	errMock := errors.New("fake error")
	messageMock := &tgbotapi.Message{
		Text: "Milk",
		Chat: &tgbotapi.Chat{ID: 123},
	}

	t.Run("Storage error", func(t *testing.T) {
		stMock.EXPECT().GetShoppingItems(gomock.Any()).Return(nil, errMock)

		err := handleDel(clientMock, stMock, messageMock)

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

			err := handleDel(clientMock, stMock, messageMock)
			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})

		t.Run("Shopping list with items", func(t *testing.T) {
			// Data mocks
			storageDataMock := []*models.ShoppingItem{
				&models.ShoppingItem{ID: 1, Name: "Milk"},
				&models.ShoppingItem{ID: 2, Name: "Молоко"},
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

				inlineKeyboardMarkup, ok := msgCfg.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
				if !ok {
					t.Fatal("Expected message to contain inline keybaord")
				}

				rowsNumber := len(inlineKeyboardMarkup.InlineKeyboard)
				expectedRowsNumber := len(storageDataMock)
				if rowsNumber != expectedRowsNumber {
					t.Fatalf(
						"Expected number of rows is %d, got %d",
						expectedRowsNumber, rowsNumber,
					)
				}

				for rowIndex, keyboardRow := range inlineKeyboardMarkup.InlineKeyboard {
					for _, keyboardButton := range keyboardRow {
						expectedCallbackData := fmt.Sprintf(
							"%s:%d", commandDel, storageDataMock[rowIndex].ID,
						)

						if *keyboardButton.CallbackData != expectedCallbackData {
							t.Errorf(
								"Expected callback data is %#v, got %#v",
								expectedCallbackData,
								*keyboardButton.CallbackData,
							)
						}
					}
				}
			})

			err := handleDel(clientMock, stMock, messageMock)
			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})
	})
}

func TestHandleDelCallbackQuery(t *testing.T) {
	// Common interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	// Common data mocks
	errMock := errors.New("Fake error")
	callbackQueryMock := &tgbotapi.CallbackQuery{
		ID: "some-callback-id",
		Message: &tgbotapi.Message{
			MessageID: 123,
			Chat:      &tgbotapi.Chat{ID: 123},
			Text:      "Milk",
		},
	}

	generateCallbackQueryIDChecker := func(t *testing.T) interface{} {
		return func(config tgbotapi.CallbackConfig) {
			if config.CallbackQueryID != callbackQueryMock.ID {
				t.Fatalf(
					"Expected callaback query ID %s, got %s",
					callbackQueryMock.ID, config.CallbackQueryID,
				)
			}
		}
	}

	t.Run("Callback data parsing error", func(t *testing.T) {
		// Data mocks
		invalidData := "not int"

		// Interface mocks
		clientMock.EXPECT().AnswerCallbackQuery(
			gomock.Any(),
		).Do(generateCallbackQueryIDChecker(t))

		err := handleDelCallbackQuery(clientMock, stMock, callbackQueryMock, invalidData)
		if !strings.Contains(err.Error(), strconv.ErrSyntax.Error()) {
			t.Errorf(
				"Expected error to contain %#v, got %#v",
				strconv.ErrSyntax.Error(), err.Error(),
			)
		}
	})

	t.Run("Storage error", func(t *testing.T) {
		t.Run("GetShoppingItem", func(t *testing.T) {
			// Data mocks
			expectedItemID := int64(123)
			dataMock := strconv.FormatInt(expectedItemID, 10)

			// Interface mocks
			clientMock.EXPECT().AnswerCallbackQuery(
				gomock.Any(),
			).Do(generateCallbackQueryIDChecker(t))

			stMock.EXPECT().GetShoppingItem(expectedItemID).Return(nil, errMock)

			err := handleDelCallbackQuery(clientMock, stMock, callbackQueryMock, dataMock)
			if !strings.Contains(err.Error(), errMock.Error()) {
				t.Errorf(
					"Expected error to contain %#v, got %#v",
					errMock.Error(), err.Error(),
				)
			}
		})
		t.Run("DeleteShoppingItem", func(t *testing.T) {
			// Data mocks
			expectedItemID := int64(123)
			dataMock := strconv.FormatInt(expectedItemID, 10)
			item := &models.ShoppingItem{Name: "Milk"}

			// Interface mocks
			clientMock.EXPECT().AnswerCallbackQuery(
				gomock.Any(),
			).Do(generateCallbackQueryIDChecker(t))

			stMock.EXPECT().GetShoppingItem(expectedItemID).Return(item, nil)
			stMock.EXPECT().DeleteShoppingItem(expectedItemID).Return(errMock)

			err := handleDelCallbackQuery(clientMock, stMock, callbackQueryMock, dataMock)
			if !strings.Contains(err.Error(), errMock.Error()) {
				t.Errorf(
					"Expected error to contain %#v, got %#v",
					errMock.Error(), err.Error(),
				)
			}
		})
	})

	t.Run("Success", func(t *testing.T) {
		// Common data mocks
		expectedItemID := int64(123)
		dataMock := strconv.FormatInt(expectedItemID, 10)
		item := &models.ShoppingItem{Name: "Milk"}

		t.Run("Item wasn't found", func(t *testing.T) {
			// Interface mocks
			clientMock.EXPECT().AnswerCallbackQuery(
				gomock.Any(),
			).Do(generateCallbackQueryIDChecker(t))
			stMock.EXPECT().GetShoppingItem(expectedItemID).Return(nil, nil)

			sendHideKeybaordCall := clientMock.EXPECT().Send(gomock.Any())
			sendHideKeybaordCall.Do(
				generateSendHideKeybaordCallChecker(t, callbackQueryMock),
			)
			sendTextCall := clientMock.EXPECT().Send(gomock.Any())
			sendTextCall.Do(func(msgCfg tgbotapi.MessageConfig) {
				if msgCfg.ChatID != callbackQueryMock.Message.Chat.ID {
					t.Errorf(
						"Expected to reply to the chat with ID %d, but reply sent to %d",
						msgCfg.ChatID,
						callbackQueryMock.Message.Chat.ID,
					)
				}

				expectedText := "Can't find an item"
				if !strings.Contains(msgCfg.Text, expectedText) {
					t.Fatalf(
						"Expected message to contain %#v. Got: %#v",
						expectedText, msgCfg.Text,
					)
				}
			})

			err := handleDelCallbackQuery(clientMock, stMock, callbackQueryMock, dataMock)
			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
			}
		})

		t.Run("Item found", func(t *testing.T) {
			// Interface mocks
			clientMock.EXPECT().AnswerCallbackQuery(
				gomock.Any(),
			).Do(generateCallbackQueryIDChecker(t))
			stMock.EXPECT().GetShoppingItem(expectedItemID).Return(item, nil)
			stMock.EXPECT().DeleteShoppingItem(expectedItemID).Return(nil)

			sendHideKeybaordCall := clientMock.EXPECT().Send(gomock.Any())
			sendHideKeybaordCall.Do(
				generateSendHideKeybaordCallChecker(t, callbackQueryMock),
			)
			sendTextCall := clientMock.EXPECT().Send(gomock.Any())
			sendTextCall.Do(func(msgCfg tgbotapi.MessageConfig) {
				if msgCfg.ChatID != callbackQueryMock.Message.Chat.ID {
					t.Errorf(
						"Expected to reply to the chat with ID %d, but reply sent to %d",
						msgCfg.ChatID,
						callbackQueryMock.Message.Chat.ID,
					)
				}

				if !strings.Contains(msgCfg.Text, item.Name) {
					t.Fatalf(
						"Expected message to contain the item name: %#v, got: %#v",
						item.Name, msgCfg.Text,
					)
				}
			})

			err := handleDelCallbackQuery(clientMock, stMock, callbackQueryMock, dataMock)
			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
			}
		})

	})
}

func TestHandleClear(t *testing.T) {
	// Common interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMocksender(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	// Common data mocks
	errMock := errors.New("fake error")
	messageMock := &tgbotapi.Message{
		Text: "Milk",
		Chat: &tgbotapi.Chat{ID: 123},
	}

	t.Run("Storage error", func(t *testing.T) {
		stMock.EXPECT().GetShoppingItems(gomock.Any()).Return(nil, errMock)

		err := handleClear(clientMock, stMock, messageMock)

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

			err := handleClear(clientMock, stMock, messageMock)
			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})

		t.Run("Shopping list with items", func(t *testing.T) {
			// Data mocks
			storageDataMock := []*models.ShoppingItem{
				&models.ShoppingItem{ID: 1, Name: "Milk"},
				&models.ShoppingItem{ID: 2, Name: "Молоко"},
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

				inlineKeyboardMarkup, ok := msgCfg.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
				if !ok {
					t.Fatal("Expected message to contain inline keybaord")
				}

				expectedRowsNumber := 1
				rowsNumber := len(inlineKeyboardMarkup.InlineKeyboard)
				if rowsNumber != expectedRowsNumber {
					t.Fatalf(
						"Expected number of rows is %d, got %d",
						expectedRowsNumber, rowsNumber,
					)
				}

				expectedCallbackData := []string{
					fmt.Sprintf("%s:%s", commandClear, clearCallbackDataConfim),
					fmt.Sprintf("%s:%s", commandClear, clearCallbackDataCancel),
				}

				// Order of the buttons in keyboard is important
				for buttonOrder, callbackData := range expectedCallbackData {
					keyboardButton := inlineKeyboardMarkup.InlineKeyboard[0][buttonOrder]

					if *keyboardButton.CallbackData != callbackData {
						t.Errorf(
							"Expected callback data is %#v, got %#v",
							expectedCallbackData,
							*keyboardButton.CallbackData,
						)
					}
				}
			})

			err := handleClear(clientMock, stMock, messageMock)
			if err != nil {
				t.Errorf("Unexpected err: got %#v", err)
			}
		})
	})
}

func TestHandleClearCallbackQuery(t *testing.T) {
	// Common interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	// Common data mocks
	errMock := errors.New("Fake error")
	callbackQueryMock := &tgbotapi.CallbackQuery{
		ID: "some-callback-id",
		Message: &tgbotapi.Message{
			MessageID: 123,
			Chat:      &tgbotapi.Chat{ID: 123},
			Text:      "Milk",
		},
	}

	generateCallbackQueryIDChecker := func(t *testing.T) interface{} {
		return func(config tgbotapi.CallbackConfig) {
			if config.CallbackQueryID != callbackQueryMock.ID {
				t.Fatalf(
					"Expected callaback query ID %s, got %s",
					callbackQueryMock.ID, config.CallbackQueryID,
				)
			}
		}
	}

	t.Run("Callback data parsing error", func(t *testing.T) {
		// Data mocks
		invalidData := "not bool"

		// Interface mocks
		clientMock.EXPECT().AnswerCallbackQuery(
			gomock.Any(),
		).Do(generateCallbackQueryIDChecker(t))

		err := handleClearCallbackQuery(clientMock, stMock, callbackQueryMock, invalidData)

		expectedErrorText := "Unable to parse confirmation"
		if !strings.Contains(err.Error(), expectedErrorText) {
			t.Errorf(
				"Expected error to contain %#v, got %#v",
				expectedErrorText, err.Error(),
			)
		}
	})

	t.Run("Storage error", func(t *testing.T) {
		t.Run("DeleteAllShoppingItems", func(t *testing.T) {
			// Interface mocks
			clientMock.EXPECT().AnswerCallbackQuery(
				gomock.Any(),
			).Do(generateCallbackQueryIDChecker(t))

			stMock.EXPECT().DeleteAllShoppingItems(
				callbackQueryMock.Message.Chat.ID,
			).Return(errMock)

			err := handleClearCallbackQuery(
				clientMock, stMock, callbackQueryMock, clearCallbackDataConfim,
			)
			if !strings.Contains(err.Error(), errMock.Error()) {
				t.Errorf(
					"Expected error to contain %#v, got %#v",
					errMock.Error(), err.Error(),
				)
			}
		})
	})

	t.Run("Success", func(t *testing.T) {
		t.Run("User confirms deletion", func(t *testing.T) {
			// Interface mocks
			clientMock.EXPECT().AnswerCallbackQuery(
				gomock.Any(),
			).Do(generateCallbackQueryIDChecker(t))
			stMock.EXPECT().DeleteAllShoppingItems(
				callbackQueryMock.Message.Chat.ID,
			).Return(nil)

			sendHideKeybaordCall := clientMock.EXPECT().Send(gomock.Any())
			sendHideKeybaordCall.Do(
				generateSendHideKeybaordCallChecker(t, callbackQueryMock),
			)
			sendTextCall := clientMock.EXPECT().Send(gomock.Any())
			sendTextCall.Do(func(msgCfg tgbotapi.MessageConfig) {
				if msgCfg.ChatID != callbackQueryMock.Message.Chat.ID {
					t.Errorf(
						"Expected to reply to the chat with ID %d, but reply sent to %d",
						msgCfg.ChatID,
						callbackQueryMock.Message.Chat.ID,
					)
				}

				expectedText := "deleted all items"
				if !strings.Contains(msgCfg.Text, expectedText) {
					t.Fatalf(
						"Expected message to contain %#v. Got: %#v",
						expectedText, msgCfg.Text,
					)
				}
			})

			err := handleClearCallbackQuery(
				clientMock, stMock, callbackQueryMock, clearCallbackDataConfim,
			)
			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
			}
		})

		t.Run("User cancels deletion", func(t *testing.T) {
			// Interface mocks
			clientMock.EXPECT().AnswerCallbackQuery(
				gomock.Any(),
			).Do(generateCallbackQueryIDChecker(t))

			sendHideKeybaordCall := clientMock.EXPECT().Send(gomock.Any())
			sendHideKeybaordCall.Do(
				generateSendHideKeybaordCallChecker(t, callbackQueryMock),
			)
			sendTextCall := clientMock.EXPECT().Send(gomock.Any())
			sendTextCall.Do(func(msgCfg tgbotapi.MessageConfig) {
				if msgCfg.ChatID != callbackQueryMock.Message.Chat.ID {
					t.Errorf(
						"Expected to reply to the chat with ID %d, but reply sent to %d",
						msgCfg.ChatID,
						callbackQueryMock.Message.Chat.ID,
					)
				}

				expectedText := "Canceling"
				if !strings.Contains(msgCfg.Text, expectedText) {
					t.Fatalf(
						"Expected message to contain %#v. Got: %#v",
						expectedText, msgCfg.Text,
					)
				}
			})

			err := handleClearCallbackQuery(
				clientMock, stMock, callbackQueryMock, clearCallbackDataCancel,
			)
			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
			}
		})

	})
}

// --- Utils

func generateSendHideKeybaordCallChecker(
	t *testing.T,
	callbackQueryMock *tgbotapi.CallbackQuery,
) interface{} {
	return func(msgCfg tgbotapi.EditMessageReplyMarkupConfig) {
		if msgCfg.ChatID != callbackQueryMock.Message.Chat.ID {
			t.Errorf(
				"Expected to reply to the chat with ID %d, but reply sent to %d",
				msgCfg.ChatID,
				callbackQueryMock.Message.Chat.ID,
			)
		}

		hasOneRow := len(msgCfg.ReplyMarkup.InlineKeyboard) == 1
		firstRowIsEmpty := len(msgCfg.ReplyMarkup.InlineKeyboard[0]) == 1
		if hasOneRow && firstRowIsEmpty {
			t.Error(
				"Expected the message update to contain empty inline",
				"keyboard layout to hide the keybaord",
			)
		}
	}
}
