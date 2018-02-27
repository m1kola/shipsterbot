package telegram

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/m1kola/shipsterbot/mocks/bot/mock_telegram"
	"github.com/m1kola/shipsterbot/mocks/mock_storage"
	"github.com/m1kola/shipsterbot/storage"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func TestHandleUpdates(t *testing.T) {
	// Interface mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mock_telegram.NewMockbotClientInterface(mockCtrl)
	stMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	// Function mocks
	isHandleUpdateCalled := make(chan bool)
	oldHandleUpdate := handleUpdate
	defer func() { handleUpdate = oldHandleUpdate }()
	handleUpdate = func(
		botClientInterface, storage.DataStorageInterface, tgbotapi.Update,
	) {
		isHandleUpdateCalled <- true
	}

	updates := make(chan tgbotapi.Update)
	defer close(updates)

	go handleUpdates(clientMock, stMock, updates)

	updates <- tgbotapi.Update{}
	if !<-isHandleUpdateCalled {
		t.Error("The handleUpdate func wasn't called")
	}
}
