// This file contains internal interfaces for vendor libraries
// that do not provide own interfaces

package telegram

// Generates mocks for tests
//go:generate mockgen -source=$GOFILE -destination=../../mocks/bot/mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type webhookListener interface {
	ListenForWebhook(pattern string) tgbotapi.UpdatesChannel
}

// listenAndServe replicates some http.Server signatures
type listenerAndServer interface {
	ListenAndServeTLS(certFile, keyFile string) error
	ListenAndServe() error
}

type sender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type callbackQueryAnswerer interface {
	AnswerCallbackQuery(config tgbotapi.CallbackConfig) (tgbotapi.APIResponse, error)
}

type botClientInterface interface {
	webhookListener
	sender
	callbackQueryAnswerer
}
