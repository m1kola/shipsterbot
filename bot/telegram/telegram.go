package telegram

import (
	"log"
	"net/http"

	"github.com/m1kola/shipsterbot/storage"
)

type webHookServerConfig struct {
	port        string
	TLSCertPath string
	TLSKeyPath  string
}

// BotApp is a struct for handeling interactions
// with the Telegram API
type BotApp struct {
	bot          botClientInterface
	storage      storage.DataStorageInterface
	serverConfig *webHookServerConfig
}

// NewBotApp creates a new instance of a bot struct
func NewBotApp(
	tgbot *APIClient,
	storage storage.DataStorageInterface,
	port, TLSCertPath, TLSKeyPath string,
) *BotApp {
	serverConfig := &webHookServerConfig{
		port:        port,
		TLSCertPath: TLSCertPath,
		TLSKeyPath:  TLSKeyPath,
	}

	return &BotApp{
		bot:          tgbot.botClient,
		storage:      storage,
		serverConfig: serverConfig,
	}
}

// StartBotApp starts the  bot
func StartBotApp(bapp *BotApp) error {
	updates := getUpdatesChan(bapp.bot)
	go routeUpdates(bapp.bot, bapp.storage, updates)

	server := newServerWithIncommingRequstLogger(
		bapp.serverConfig.port, http.DefaultServeMux)

	log.Printf("Start listening on %s", server.Addr)
	err := listenAndServe(
		server,
		bapp.serverConfig.TLSCertPath,
		bapp.serverConfig.TLSKeyPath,
	)
	return err
}
