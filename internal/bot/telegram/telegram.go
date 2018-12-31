package telegram

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/m1kola/shipsterbot/internal/pkg/storage"
)

const defaultServerPort = "8443"

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

// WebhookTLS allows webhook webserver to operate in secure transport mode
func WebhookTLS(TLSCertPath, TLSKeyPath string) func(*BotApp) error {
	return func(app *BotApp) error {
		app.serverConfig.TLSCertPath = TLSCertPath
		app.serverConfig.TLSKeyPath = TLSKeyPath

		return nil
	}
}

// WebhookPort sets a custom webhook port
func WebhookPort(port string) func(*BotApp) error {
	return func(app *BotApp) error {
		err := ValidateWebhookPort(port)
		if err != nil {
			return err
		}

		app.serverConfig.port = port
		return nil
	}
}

var tgbotapiNewBotAPI = tgbotapi.NewBotAPI

// NewBotApp creates a new instance of a bot struct
func NewBotApp(
	storage storage.DataStorageInterface,
	apiToken string,
	options ...(func(*BotApp) error),
) (*BotApp, error) {
	client, err := tgbotapiNewBotAPI(apiToken)
	if err != nil {
		return nil, err
	}

	serverConfig := &webHookServerConfig{
		port: defaultServerPort,
	}

	botApp := &BotApp{
		bot:          &apiClientWrapper{client},
		storage:      storage,
		serverConfig: serverConfig,
	}

	for _, option := range options {
		err := option(botApp)
		if err != nil {
			return nil, err
		}
	}

	return botApp, nil
}

// StartBotApp starts the  bot
func StartBotApp(bapp *BotApp) error {
	updates := getUpdatesChan(bapp.bot)
	go routeUpdates(bapp.bot, bapp.storage, updates)

	server := newServerWithincomingRequestLogger(
		bapp.serverConfig.port, http.DefaultServeMux)

	log.Printf("Start listening on %s", server.Addr)
	err := listenAndServe(
		server,
		bapp.serverConfig.TLSCertPath,
		bapp.serverConfig.TLSKeyPath,
	)
	return err
}
