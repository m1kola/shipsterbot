package telegram

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// apiClientWrapper is required to statisfy internal interfaces
type apiClientWrapper struct {
	*tgbotapi.BotAPI
}

func (api *apiClientWrapper) Token() string {
	return api.BotAPI.Token
}
