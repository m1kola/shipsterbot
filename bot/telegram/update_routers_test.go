package telegram

import (
	"errors"
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
