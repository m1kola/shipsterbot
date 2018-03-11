package telegram

import (
	"errors"
	"testing"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func TestNewAPIClient(t *testing.T) {
	t.Run("Without error", func(t *testing.T) {
		expectedToken := "123"
		oldTgbotapiNewBotAPI := tgbotapiNewBotAPI
		defer func() { tgbotapiNewBotAPI = oldTgbotapiNewBotAPI }()
		tgbotapiNewBotAPI = func(token string) (*tgbotapi.BotAPI, error) {
			tgbotapiClient := &tgbotapi.BotAPI{
				Token: expectedToken,
			}

			return tgbotapiClient, nil
		}

		client, err := NewAPIClient(expectedToken)
		if err != nil {
			t.Errorf("Got error %v, expected nil", err)
		}

		actualToken := client.Token()
		if expectedToken != actualToken {
			t.Errorf("Expected token %s, got %s", expectedToken, actualToken)
		}
	})
	t.Run("Without error", func(t *testing.T) {
		expectedErr := errors.New("Fake error")
		oldTgbotapiNewBotAPI := tgbotapiNewBotAPI
		defer func() { tgbotapiNewBotAPI = oldTgbotapiNewBotAPI }()
		tgbotapiNewBotAPI = func(token string) (*tgbotapi.BotAPI, error) {
			return nil, expectedErr
		}

		_, err := NewAPIClient("123")
		if expectedErr != err {
			t.Errorf("Expected error %v, got %s", expectedErr, err)
		}
	})
}

func TestApiClientWrapper(t *testing.T) {
	t.Run("Token", func(t *testing.T) {
		expectedToken := "123"

		APIClientWrap := apiClientWrapper{
			&tgbotapi.BotAPI{
				Token: expectedToken,
			},
		}

		actualToken := APIClientWrap.Token()
		if expectedToken != actualToken {
			t.Errorf("Expected token %s, got %s", expectedToken, actualToken)
		}
	})
}

func TestAPIClient(t *testing.T) {
	expectedToken := "123"
	expectedBotUserName := "m1kola"
	client := APIClient{
		botClient: &apiClientWrapper{
			&tgbotapi.BotAPI{
				Token: expectedToken,
				Self: tgbotapi.User{
					UserName: expectedBotUserName,
				},
			},
		},
	}

	t.Run("Token", func(t *testing.T) {
		actualToken := client.Token()
		if expectedToken != actualToken {
			t.Errorf("Expected token %s, got %s",
				expectedToken, actualToken)
		}
	})

	t.Run("BotUserName", func(t *testing.T) {
		actualBotUserName := client.BotUserName()
		if expectedBotUserName != actualBotUserName {
			t.Errorf("Expected BotUserName %s, got %s",
				expectedBotUserName, actualBotUserName)
		}
	})

	t.Run("SetDebug", func(t *testing.T) {
		if client.botClient.Debug {
			t.Error("Debug mode should be turned off by default")
		}

		client.SetDebug(true)
		if !client.botClient.Debug {
			t.Error("Debug mode is expected to be turned on after a client.SetDebug(true) call")
		}

		client.SetDebug(false)
		if client.botClient.Debug {
			t.Error("Debug mode is expected to be turned off after a client.SetDebug(false) call")
		}

	})
}
