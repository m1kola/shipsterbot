package telegram

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/m1kola/shipsterbot/gomockhelpers"
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

	t.Run("CallbackQuery update: handlerCanNotHandleError", func(t *testing.T) {
		// Data mocks
		messageMock := &tgbotapi.Message{}
		callbackQueryMock := &tgbotapi.CallbackQuery{
			Message: messageMock,
		}
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

			return handlerCanNotHandleError{errors.New("Fake error")}
		}

		handleUnrecognisedMessageIsCalled := false
		handleUnrecognisedMessageOld := handleUnrecognisedMessage
		defer func() { handleUnrecognisedMessage = handleUnrecognisedMessageOld }()
		handleUnrecognisedMessage = func(_ sender, message *tgbotapi.Message) {
			handleUnrecognisedMessageIsCalled = true
			if message != messageMock {
				t.Error("Wrong Message received")
			}
		}

		routeUpdate(clientMock, stMock, updateMock)

		if !routeCallbackQueryIsCalled {
			t.Error("func handleCallbackQuery wasn't called")
		}
		if !handleUnrecognisedMessageIsCalled {
			t.Error("func handleUnrecognisedMessage wasn't called")
		}

	})

	t.Run("Message update: handlerCanNotHandleError", func(t *testing.T) {
		// Data mocks
		messageMock := &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
		}
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

			return handlerCanNotHandleError{errors.New("Fake error")}
		}

		handleUnrecognisedMessageIsCalled := false
		handleUnrecognisedMessageOld := handleUnrecognisedMessage
		defer func() { handleUnrecognisedMessage = handleUnrecognisedMessageOld }()
		handleUnrecognisedMessage = func(_ sender, message *tgbotapi.Message) {
			handleUnrecognisedMessageIsCalled = true
			if message != messageMock {
				t.Error("Wrong Message received")
			}
		}

		routeUpdate(clientMock, stMock, updateMock)

		if !routeMessageIsCalled {
			t.Error("func handleCallbackQuery wasn't called")
		}
		if !handleUnrecognisedMessageIsCalled {
			t.Error("func handleUnrecognisedMessage wasn't called")
		}
	})

	t.Run("CallbackQuery update error", func(t *testing.T) {
		// Data mocks
		messageMock := &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
		}
		callbackQueryMock := &tgbotapi.CallbackQuery{
			Message: messageMock,
		}
		updateMock := tgbotapi.Update{
			CallbackQuery: callbackQueryMock,
		}

		// Function mocks
		sendMsgMatcher := gomockhelpers.MatcherFunc(func(x interface{}) bool {
			if m, ok := x.(tgbotapi.MessageConfig); ok {
				return m.ChatID == messageMock.Chat.ID
			}

			return false
		})
		clientMock.EXPECT().Send(sendMsgMatcher)

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

			return errors.New("Fake error")
		}

		routeUpdate(clientMock, stMock, updateMock)

		if !routeCallbackQueryIsCalled {
			t.Error("func handleCallbackQuery wasn't called")
		}
	})

	t.Run("Message update error", func(t *testing.T) {
		// Data mocks
		messageMock := &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
		}
		updateMock := tgbotapi.Update{
			Message: messageMock,
		}

		// Function mocks
		sendMsgMatcher := gomockhelpers.MatcherFunc(func(x interface{}) bool {
			if m, ok := x.(tgbotapi.MessageConfig); ok {
				return m.ChatID == messageMock.Chat.ID
			}

			return false
		})
		clientMock.EXPECT().Send(sendMsgMatcher)

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

			return errors.New("Fake error")
		}

		routeUpdate(clientMock, stMock, updateMock)

		if !routeMessageIsCalled {
			t.Error("func handleCallbackQuery wasn't called")
		}
	})
}
