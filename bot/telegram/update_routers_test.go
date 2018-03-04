package telegram

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/m1kola/shipsterbot/mocks/bot/mock_telegram"
	"github.com/m1kola/shipsterbot/mocks/mock_storage"
	"github.com/m1kola/shipsterbot/storage"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func TestRouteUpdates(t *testing.T) {
	// Interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	// Function mocks
	routeUpdateIsCalled := make(chan bool)
	routeUpdateOld := routeUpdate
	defer func() { routeUpdate = routeUpdateOld }()
	routeUpdate = func(
		botClientInterface, storage.DataStorageInterface, tgbotapi.Update,
	) {
		routeUpdateIsCalled <- true
	}

	updates := make(chan tgbotapi.Update)
	defer close(updates)

	go routeUpdates(clientMock, stMock, updates)

	updates <- tgbotapi.Update{}
	if !<-routeUpdateIsCalled {
		t.Error("The routeUpdate func wasn't called")
	}
}

func TestRouteUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	t.Run("CallbackQuery update", func(t *testing.T) {
		// Data mocks
		callbackQueryMock := &tgbotapi.CallbackQuery{}
		updateMock := tgbotapi.Update{
			CallbackQuery: callbackQueryMock,
		}

		// Function mocks
		routeCallbackQueryIsCalled := false
		routeCallbackQueryOld := routeCallbackQuery
		defer func() { routeCallbackQuery = routeCallbackQueryOld }()
		routeCallbackQuery = func(
			_ botClientInterface,
			_ storage.DataStorageInterface,
			callbackQuery *tgbotapi.CallbackQuery,
		) error {
			routeCallbackQueryIsCalled = true
			if callbackQueryMock != callbackQuery {
				t.Error("Wrong CallbackQuery received")
			}

			return nil
		}

		routeUpdate(clientMock, stMock, updateMock)

		if !routeCallbackQueryIsCalled {
			t.Error("func handleCallbackQuery wasn't called")
		}
	})

	t.Run("Message update", func(t *testing.T) {
		// Data mocks
		messageMock := &tgbotapi.Message{}
		updateMock := tgbotapi.Update{
			Message: messageMock,
		}

		// Function mocks
		routeMessageIsCalled := false
		routeMessageOld := routeMessage
		defer func() { routeMessage = routeMessageOld }()
		routeMessage = func(
			_ sender,
			_ storage.DataStorageInterface,
			message *tgbotapi.Message,
		) error {
			routeMessageIsCalled = true
			if message != messageMock {
				t.Error("Wrong Message received")
			}

			return nil
		}

		routeUpdate(clientMock, stMock, updateMock)

		if !routeMessageIsCalled {
			t.Error("func handleCallbackQuery wasn't called")
		}
	})

	t.Run("CallbackQuery update error", func(t *testing.T) {
		// Data mocks
		errMock := errors.New("Fake error")
		messageMock := &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
		}
		updateMock := tgbotapi.Update{
			CallbackQuery: &tgbotapi.CallbackQuery{
				Message: messageMock,
			},
		}

		// Function mocks
		routeCallbackQueryOld := routeCallbackQuery
		defer func() { routeCallbackQuery = routeCallbackQueryOld }()
		routeCallbackQuery = func(_ botClientInterface, _ storage.DataStorageInterface, _ *tgbotapi.CallbackQuery) error {
			return errMock
		}

		routeErrorsIsCalled := false
		routeErrorsOld := routeErrors
		defer func() { routeErrors = routeErrorsOld }()
		routeErrors = func(
			_ botClientInterface,
			actualMessage *tgbotapi.Message,
			actualErr error,
		) {
			routeErrorsIsCalled = true

			if actualMessage != messageMock {
				t.Errorf("got %#v, expected %#v", actualMessage, messageMock)
			}

			if actualErr != errMock {
				t.Errorf("got %#v, expected %#v", actualErr, errMock)
			}
		}

		routeUpdate(clientMock, stMock, updateMock)

		if !routeErrorsIsCalled {
			t.Error("func routeErrors wasn't called")
		}
	})

	t.Run("Message update error", func(t *testing.T) {
		// Data mocks
		errMock := errors.New("Fake error")
		messageMock := &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
		}
		updateMock := tgbotapi.Update{
			Message: messageMock,
		}

		// Function mocks
		routeMessageOld := routeMessage
		defer func() { routeMessage = routeMessageOld }()
		routeMessage = func(_ sender, _ storage.DataStorageInterface, _ *tgbotapi.Message) error {
			return errMock
		}

		routeErrorsIsCalled := false
		routeErrorsOld := routeErrors
		defer func() { routeErrors = routeErrorsOld }()
		routeErrors = func(
			_ botClientInterface,
			actualMessage *tgbotapi.Message,
			actualErr error,
		) {
			routeErrorsIsCalled = true

			if actualMessage != messageMock {
				t.Errorf("got %#v, expected %#v", actualMessage, messageMock)
			}

			if actualErr != errMock {
				t.Errorf("got %#v, expected %#v", actualErr, errMock)
			}
		}

		routeUpdate(clientMock, stMock, updateMock)

		if !routeErrorsIsCalled {
			t.Error("func routeErrors wasn't called")
		}
	})
}

func TestRouteErrors(t *testing.T) {
	// Common Interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)

	// Common data mocks
	messageMock := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 123,
		},
	}

	t.Run("error is nil", func(t *testing.T) {
		// Function mocks
		handleUnrecoverableErrorIsCalled := false
		handleUnrecoverableErrorOld := handleUnrecoverableError
		defer func() { handleUnrecoverableError = handleUnrecoverableErrorOld }()
		handleUnrecoverableError = func(_ botClientInterface, _ int64, _ error) {
			handleUnrecoverableErrorIsCalled = true
		}

		routeErrors(clientMock, messageMock, nil)

		if handleUnrecoverableErrorIsCalled {
			t.Error("routeErrors must not continue routing, when it receives err == nil")
		}
	})

	t.Run("error is not nil", func(t *testing.T) {
		// Data mocks
		errMock := errors.New("fake error")

		// Function mocks
		handleUnrecoverableErrorIsCalled := false
		handleUnrecoverableErrorOld := handleUnrecoverableError
		defer func() { handleUnrecoverableError = handleUnrecoverableErrorOld }()
		handleUnrecoverableError = func(
			_ botClientInterface, actualChatID int64, actualErr error,
		) {
			handleUnrecoverableErrorIsCalled = true

			if actualErr != errMock {
				t.Errorf("got %#v, expected %#v", actualErr, errMock)
			}

			if actualChatID != messageMock.Chat.ID {
				t.Errorf("got chat ID == %d, expected %d",
					actualChatID, messageMock.Chat.ID)
			}
		}

		routeErrors(clientMock, messageMock, errMock)

		if !handleUnrecoverableErrorIsCalled {
			t.Error("handleUnrecoverableError wasn't called")
		}
	})

	t.Run("error is handlerCanNotHandleError", func(t *testing.T) {
		// Data mocks
		errMock := handlerCanNotHandleError{
			errors.New("fake error")}

		// Function mocks
		handleUnrecognisedMessageIsCalled := false
		handleUnrecognisedMessageOld := handleUnrecognisedMessage
		defer func() { handleUnrecognisedMessage = handleUnrecognisedMessageOld }()
		handleUnrecognisedMessage = func(
			_ sender, actualMessage *tgbotapi.Message,
		) {
			handleUnrecognisedMessageIsCalled = true

			if actualMessage != messageMock {
				t.Errorf("got %#v, expected %#v", actualMessage, messageMock)
			}
		}

		routeErrors(clientMock, messageMock, errMock)

		if !handleUnrecognisedMessageIsCalled {
			t.Error("handleUnrecoverableError wasn't called")
		}
	})
}

// TODO: Simplify tests: not it's easier to mock clear and del commands
func TestRouteCallbackQuery(t *testing.T) {
	// Interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	t.Run("Commands", func(t *testing.T) {

		t.Run("Del command", func(t *testing.T) {
			// Data mocks
			expectedPayload := "123"
			callbackQueryMock := &tgbotapi.CallbackQuery{
				Data: "del:123",
			}
			errMock := errors.New("Fake error")

			// Function mocks
			handleDelCallbackQueryOld := handleDelCallbackQuery
			defer func() { handleDelCallbackQuery = handleDelCallbackQueryOld }()
			handleDelCallbackQuery = func(
				_ botClientInterface,
				_ storage.DataStorageInterface,
				callbackQuery *tgbotapi.CallbackQuery,
				payload string,
			) error {
				if callbackQueryMock != callbackQuery {
					t.Error("Unexpected callbackQuery")
				}

				if expectedPayload != payload {
					t.Errorf("Expected paylod %#v, got %#v",
						expectedPayload, payload)
				}
				return errMock
			}

			err := routeCallbackQuery(clientMock, stMock, callbackQueryMock)
			if errMock != err {
				t.Fatalf("Expected the %#v error, got %#v", errMock, err)
			}
		})

		t.Run("Clear command", func(t *testing.T) {
			// Data mocks
			expectedPayload := "123"
			callbackQueryMock := &tgbotapi.CallbackQuery{
				Data: "clear:123",
			}
			errMock := errors.New("Fake error")

			// Function mocks
			handleClearCallbackQueryOld := handleClearCallbackQuery
			defer func() { handleClearCallbackQuery = handleClearCallbackQueryOld }()
			handleClearCallbackQuery = func(
				_ botClientInterface,
				_ storage.DataStorageInterface,
				callbackQuery *tgbotapi.CallbackQuery,
				payload string,
			) error {
				if callbackQueryMock != callbackQuery {
					t.Error("Unexpected callbackQuery")
				}

				if expectedPayload != payload {
					t.Errorf("Expected paylod %#v, got %#v",
						expectedPayload, payload)
				}
				return errMock
			}

			err := routeCallbackQuery(clientMock, stMock, callbackQueryMock)
			if errMock != err {
				t.Fatalf("Expected the %#v error, got %#v", errMock, err)
			}
		})

		t.Run("Unknown command", func(t *testing.T) {
			// Data mocks
			callbackQueryMock := &tgbotapi.CallbackQuery{
				Data: "valid_but_unknown_command_name:123",
			}

			err := routeCallbackQuery(clientMock, stMock, callbackQueryMock)
			if _, ok := err.(handlerCanNotHandleError); !ok {
				t.Fatalf("expected %T got %T", handlerCanNotHandleError{}, err)
			}
		})
	})

	t.Run("Callback data error", func(t *testing.T) {
		// Data mocks
		callbackQueryMock := &tgbotapi.CallbackQuery{
			Data: "invalid_data",
		}

		err := routeCallbackQuery(clientMock, stMock, callbackQueryMock)
		if _, ok := err.(handlerCanNotHandleError); !ok {
			t.Fatalf("expected %T got %T", handlerCanNotHandleError{}, err)
		}
	})
}

func TestRouteMessage(t *testing.T) {
	// Mock setup funcs
	var routeMessageEntitiesMockSetup = func(errMock error) func() {
		routeMessageEntitiesOld := routeMessageEntities
		tearDownFunc := func() { routeMessageEntities = routeMessageEntitiesOld }

		routeMessageEntities = func(
			_ sender,
			_ storage.DataStorageInterface,
			_ *tgbotapi.Message,
		) error {
			return errMock
		}

		return tearDownFunc
	}

	// Common interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	// Comon data mocks
	messageMock := &tgbotapi.Message{}

	t.Run("Proxy errors from routeMessageEntities", func(t *testing.T) {

		t.Run("handlerCanNotHandleError", func(t *testing.T) {
			// Sould expect an error from routeMessageText

			t.Run("errCommandIsNotSupported", func(t *testing.T) {
				// Function mocks
				tearDownFunc := routeMessageEntitiesMockSetup(
					errCommandIsNotSupported,
				)
				defer tearDownFunc()

				err := routeMessage(clientMock, stMock, messageMock)

				if errCommandIsNotSupported != err {
					t.Errorf("Expected %#v, got %#v",
						errCommandIsNotSupported, err)
				}
			})
		})

		t.Run("Non-handlerCanNotHandleError", func(t *testing.T) {
			// Data mocks
			errMock := errors.New("Fake error")

			// Function mocks
			tearDownFunc := routeMessageEntitiesMockSetup(errMock)
			defer tearDownFunc()

			err := routeMessage(clientMock, stMock, messageMock)

			if errMock != err {
				t.Errorf("Expected %#v, got %#v", errMock, err)
			}
		})
	})

	t.Run("Proxy errors from routeMessageText", func(t *testing.T) {
		// Data mocks
		errFromrouteMessageTextMock := handlerCanNotHandleError{
			errors.New("Expected fake error")}
		errFromRouteMessageEntities := handlerCanNotHandleError{
			errors.New("Fake error")}

		// Function mocks
		tearDownFunc := routeMessageEntitiesMockSetup(errFromRouteMessageEntities)
		defer tearDownFunc()

		routeMessageTextOld := routeMessageText
		defer func() { routeMessageText = routeMessageTextOld }()
		routeMessageText = func(
			client sender,
			st storage.DataStorageInterface,
			message *tgbotapi.Message,
		) error {
			return errFromrouteMessageTextMock
		}

		err := routeMessage(clientMock, stMock, messageMock)

		if errFromRouteMessageEntities == err {
			t.Fatalf("Expected %#v, got %#v", errFromrouteMessageTextMock, err)
		}

		if errFromrouteMessageTextMock != err {
			t.Fatalf("Expected %#v, got %#v", errFromrouteMessageTextMock, err)
		}
	})
}

// TODO: simplify tests: use interface mocks. See TODO in the "Commands" run
func TestRouteMessageEntities(t *testing.T) {
	// Mock setup funcs
	var messageMockSetup = func(command string) *tgbotapi.Message {
		commandWithSlash := fmt.Sprintf("/%s", command)
		message := &tgbotapi.Message{
			Entities: &[]tgbotapi.MessageEntity{
				tgbotapi.MessageEntity{
					Type:   "bot_command",
					Offset: 0,
					Length: len(commandWithSlash),
				},
			},
			Text: commandWithSlash,
		}

		return message
	}

	// Common interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	t.Run("Empty message Entities", func(t *testing.T) {
		// Data mocks
		messageMock := &tgbotapi.Message{}

		err := routeMessageEntities(clientMock, stMock, messageMock)
		if _, ok := err.(handlerCanNotHandleError); !ok {
			t.Fatalf("Expected error of type %T, got %T",
				handlerCanNotHandleError{}, err)
		}
	})

	t.Run("Commands", func(t *testing.T) {

		// TODO: Check if there a better way to if there a better way to mock funcs
		//       I tried to define a type for *func(sender, storage.DataStorageInterface, *tgbotapi.Message) error
		//       and override hander funcs dynamically, but it doesn't work.
		//       This definitely can be ahived by an interface, but it seems
		//       unnecessary to define sturct for each andler...
		//       But something like HandlerFunc from the net package should work

		t.Run("Supported command", func(t *testing.T) {
			t.Run("start and help", func(t *testing.T) {
				// Data mocks
				errMock := errors.New("Error from start or help")

				// Function mocks
				handlerFuncOld := handleStart
				defer func() { handleStart = handlerFuncOld }()
				handleStart = commandHandlerFunc(func(
					_ sender,
					_ storage.DataStorageInterface,
					_ *tgbotapi.Message,
				) error {
					return errMock
				})

				t.Run("help", func(t *testing.T) {
					messageMock := messageMockSetup("help")
					err := routeMessageEntities(clientMock, stMock, messageMock)

					if errMock != err {
						t.Errorf("Expected %#v, got %#v", errMock, err)
					}
				})

				t.Run("start", func(t *testing.T) {
					messageMock := messageMockSetup("start")
					err := routeMessageEntities(clientMock, stMock, messageMock)

					if errMock != err {
						t.Errorf("Expected %#v, got %#v", errMock, err)
					}
				})
			})

			t.Run("add", func(t *testing.T) {
				// Data mocks
				errMock := errors.New("Error from add")
				messageMock := messageMockSetup("add")

				// Function mocks
				handlerFuncOld := handleAdd
				defer func() { handleAdd = handlerFuncOld }()
				handleAdd = func(_ sender, _ storage.DataStorageInterface, _ *tgbotapi.Message) error {
					return errMock
				}

				err := routeMessageEntities(clientMock, stMock, messageMock)

				if errMock != err {
					t.Errorf("Expected %#v, got %#v", errMock, err)
				}
			})

			t.Run("list", func(t *testing.T) {
				// Data mocks
				errMock := errors.New("Error from list")
				messageMock := messageMockSetup("list")

				// Function mocks
				handlerFuncOld := handleList
				defer func() { handleList = handlerFuncOld }()
				handleList = func(_ sender, _ storage.DataStorageInterface, _ *tgbotapi.Message) error {
					return errMock
				}

				err := routeMessageEntities(clientMock, stMock, messageMock)

				if errMock != err {
					t.Errorf("Expected %#v, got %#v", errMock, err)
				}
			})

			t.Run("del", func(t *testing.T) {
				// Data mocks
				errMock := errors.New("Error from del")
				messageMock := messageMockSetup("del")

				// Function mocks
				handlerFuncOld := handleDel
				defer func() { handleDel = handlerFuncOld }()
				handleDel = func(_ sender, _ storage.DataStorageInterface, _ *tgbotapi.Message) error {
					return errMock
				}

				err := routeMessageEntities(clientMock, stMock, messageMock)

				if errMock != err {
					t.Errorf("Expected %#v, got %#v", errMock, err)
				}
			})

			t.Run("clear", func(t *testing.T) {
				// Data mocks
				errMock := errors.New("Error from clear")
				messageMock := messageMockSetup("clear")

				// Function mocks
				handlerFuncOld := handleClear
				defer func() { handleClear = handlerFuncOld }()
				handleClear = func(_ sender, _ storage.DataStorageInterface, _ *tgbotapi.Message) error {
					return errMock
				}

				err := routeMessageEntities(clientMock, stMock, messageMock)

				if errMock != err {
					t.Errorf("Expected %#v, got %#v", errMock, err)
				}
			})
		})

		t.Run("Not supported command", func(t *testing.T) {
			command := "/invalid_command"

			// Data mocks
			messageMock := messageMockSetup(command)

			err := routeMessageEntities(clientMock, stMock, messageMock)
			if errCommandIsNotSupported != err {
				t.Fatalf("Expected %#v, got %#v", errCommandIsNotSupported, err)
			}
		})
	})
}
