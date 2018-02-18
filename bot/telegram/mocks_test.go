package telegram

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TODO: Consider using gomock
// In case of negative decision, mocks from the package "mocks"
// need to be converted into simple stuct mocks or something else.

type webhookListenerMock struct {
	webhookListenerFunc func(pattern string) tgbotapi.UpdatesChannel
}

func (mock *webhookListenerMock) ListenForWebhook(pattern string) tgbotapi.UpdatesChannel {
	return mock.webhookListenerFunc(pattern)
}

type listenerAndServerMock struct {
	listenAndServeTLSFunc func(certFile, keyFile string) error
	listenAndServeFunc    func() error
}

func (mock *listenerAndServerMock) ListenAndServeTLS(certFile, keyFile string) error {
	return mock.listenAndServeTLSFunc(certFile, keyFile)
}

func (mock *listenerAndServerMock) ListenAndServe() error {
	return mock.listenAndServeFunc()
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
