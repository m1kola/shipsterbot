package telegram

import (
	"errors"
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/mock/gomock"

	"github.com/m1kola/shipsterbot/internal/pkg/mocks/mock_storage"
)

func TestNewBotApp(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	storageMock := mock_storage.NewMockDataStorageInterface(mockCtrl)

	// Mocks 3rd party call (tgbotapi.NewBotAPI)
	// which actually calls the telegram bot API for some weird reason.
	tgbotapiNewBotAPIOld := tgbotapiNewBotAPI
	defer func() { tgbotapiNewBotAPI = tgbotapiNewBotAPIOld }()
	tgbotapiNewBotAPI = func(token string) (*tgbotapi.BotAPI, error) {
		return &tgbotapi.BotAPI{}, nil
	}

	t.Run("TLS", func(t *testing.T) {
		expectedCert := "/fake/cert.pem"
		expectedKey := "/fake/key.key"

		app, err := NewBotApp(
			storageMock,
			"fake_token",
			WebhookTLS(expectedCert, expectedKey),
		)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if app.serverConfig.TLSCertPath != expectedCert {
			t.Errorf("%s expected as TLSCertPath, got %s",
				expectedCert, app.serverConfig.TLSCertPath)
		}

		if app.serverConfig.TLSKeyPath != expectedKey {
			t.Errorf("%s expected as TLSKeyPath, got %s",
				expectedKey, app.serverConfig.TLSKeyPath)
		}
	})

	t.Run("Port", func(t *testing.T) {
		t.Run("Default port", func(t *testing.T) {
			expectedPort := "8443"

			app, err := NewBotApp(
				storageMock,
				"fake_token",
			)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if app.serverConfig.port != expectedPort {
				t.Errorf("%s expected as a port, got %s",
					expectedPort, app.serverConfig.port)
			}
		})

		t.Run("Custom port", func(t *testing.T) {
			tests := []string{"443", "80", "88", "8443", "6000"}

			for _, port := range tests {
				t.Run(fmt.Sprintf("port %#v", port), func(t *testing.T) {
					app, err := NewBotApp(
						storageMock,
						"fake_token",
						WebhookPort(port),
					)

					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}

					if app.serverConfig.port != port {
						t.Fatalf("%s expected as a port, got %s",
							port, app.serverConfig.port)
					}
				})
			}
		})
	})

	t.Run("Option error", func(t *testing.T) {
		expectedErr := errors.New("Fake error")

		_, err := NewBotApp(
			storageMock,
			"fake_token",
			func(*BotApp) error { return expectedErr },
		)
		if err != expectedErr {
			t.Fatalf("expected the %v error, got %v", err, expectedErr)
		}
	})

	t.Run("API client error", func(t *testing.T) {
		expectedErr := errors.New("Fake error")

		tgbotapiNewBotAPIOld := tgbotapiNewBotAPI
		defer func() { tgbotapiNewBotAPI = tgbotapiNewBotAPIOld }()
		tgbotapiNewBotAPI = func(token string) (*tgbotapi.BotAPI, error) {
			return nil, expectedErr
		}

		_, err := NewBotApp(
			storageMock,
			"fake_token",
			func(*BotApp) error { return expectedErr },
		)
		if err != expectedErr {
			t.Fatalf("expected the %v error, got %v", err, expectedErr)
		}
	})
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
