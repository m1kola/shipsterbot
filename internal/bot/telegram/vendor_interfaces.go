// This file contains internal interfaces for vendor libraries
// that do not provide own interfaces

package telegram

// Generates mocks for tests
//go:generate mockgen -source=$GOFILE -destination=../../mocks/bot/mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type webhookListener interface {
	ListenForWebhook(pattern string) tgbotapi.UpdatesChannel
}

// listenAndServe replicates some http.Server signatures
type listenerAndServer interface {
	ListenAndServeTLS(certFile, keyFile string) error
	ListenAndServe() error
}

type tokener interface {
	Token() string
}

type sender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type callbackQueryAnswerer interface {
	AnswerCallbackQuery(config tgbotapi.CallbackConfig) (tgbotapi.APIResponse, error)
}

type tokenListenForWebhook interface {
	webhookListener
	tokener
}

type botClientInterface interface {
	webhookListener
	tokener
	sender
	callbackQueryAnswerer
}
