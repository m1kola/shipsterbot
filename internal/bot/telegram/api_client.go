package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// apiClientWrapper is required to statisfy internal interfaces
type apiClientWrapper struct {
	*tgbotapi.BotAPI
}

func (api *apiClientWrapper) Token() string {
	return api.BotAPI.Token
}
