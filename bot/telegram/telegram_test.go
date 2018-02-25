package telegram

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/m1kola/shipsterbot/mocks"
)

func TestNewBotApp(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	storageMock := mocks.NewMockDataStorageInterface(mockCtrl)
	APIClientMock := &APIClient{}

	expectedPort := "8443"
	expectedCert := "/fake/cert.pem"
	expectedKey := "/fake/key.key"

	app := NewBotApp(APIClientMock, storageMock,
		expectedPort, expectedCert, expectedKey)

	if app.serverConfig.port != expectedPort {
		t.Errorf("%s expected as a port, got %s",
			expectedPort, app.serverConfig.port)
	}

	if app.serverConfig.TLSCertPath != expectedCert {
		t.Errorf("%s expected as TLSCertPath, got %s",
			expectedCert, app.serverConfig.TLSCertPath)
	}

	if app.serverConfig.TLSKeyPath != expectedKey {
		t.Errorf("%s expected as TLSKeyPath, got %s",
			expectedKey, app.serverConfig.TLSKeyPath)
	}
}

func TestStartBotApp(t *testing.T) {
	mockBotApp := &BotApp{
		serverConfig: &webHookServerConfig{
			port: "8443",
		},
	}

	// Mock: getUpdatesChan
	oldGetUpdatesChan := getUpdatesChan
	defer func() { getUpdatesChan = oldGetUpdatesChan }()
	getUpdatesChan = func(actualBotAppClient tokenListenForWebhook) <-chan tgbotapi.Update {
		// No op mock
		return nil
	}

	t.Run("Without error", func(t *testing.T) {
		// Mock: listenAndServe
		oldListenAndServe := listenAndServe
		defer func() { listenAndServe = oldListenAndServe }()
		listenAndServe = func(server listenerAndServer, TLSCertPath, TLSKeyPath string) error {
			return nil
		}

		err := StartBotApp(mockBotApp)

		if err != nil {
			t.Errorf("Expected nil, got error %v", err)
		}
	})

	t.Run("With error", func(t *testing.T) {
		expectedErr := errors.New("Fake error")

		// Mock: listenAndServe
		oldListenAndServe := listenAndServe
		defer func() { listenAndServe = oldListenAndServe }()
		listenAndServe = func(server listenerAndServer, TLSCertPath, TLSKeyPath string) error {
			return expectedErr
		}

		err := StartBotApp(mockBotApp)

		if err != expectedErr {
			t.Errorf("Expected the %v error, got %v", expectedErr, err)
		}
	})
}
