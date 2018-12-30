package telegram

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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
