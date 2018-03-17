package telegram

import (
	"testing"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
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
