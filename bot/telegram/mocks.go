package telegram

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type webhookListenerMock struct {
	webhookListenerFunc func(pattern string) tgbotapi.UpdatesChannel
}

func (mock *webhookListenerMock) ListenForWebhook(pattern string) tgbotapi.UpdatesChannel {
	return mock.webhookListenerFunc(pattern)
}

type tokenerMock struct {
	fakeToken string
}

func (mock *tokenerMock) Token() string {
	return mock.fakeToken
}

type tokenListenForWebhookMock struct {
	*webhookListenerMock
	*tokenerMock
}
