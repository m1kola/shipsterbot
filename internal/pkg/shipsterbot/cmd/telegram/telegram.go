package telegram

import (
	"database/sql"

	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/cobra"

	"github.com/m1kola/shipsterbot/internal/pkg/env"
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
	bot          *tgbotapi.BotAPI
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

// NewBotApp creates a new instance of a bot struct
func NewBotApp(
	storage storage.DataStorageInterface,
	client *tgbotapi.BotAPI,
	options ...(func(*BotApp) error),
) (*BotApp, error) {
	serverConfig := &webHookServerConfig{
		port: defaultServerPort,
	}

	botApp := &BotApp{
		bot:          client,
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

// NewStartTelegramBotCmd creates a new cobra.Command
func NewStartTelegramBotCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "telegram",
		Short: "Start a telegram bot",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialise DB connection pool
			dbConnectionStr, err := env.GetDBConnectionString()
			if err != nil {
				log.Fatal(err)
			}
			db, err := sql.Open("postgres", dbConnectionStr)
			if err != nil {
				log.Fatal(err)
			}

			// Get bot API token
			apiToken, err := env.GetTelegramAPIToken()
			if err != nil {
				log.Fatal(err)
			}
			client, err := tgbotapi.NewBotAPI(apiToken)
			if err != nil {
				log.Fatal(err)
			}

			// Create a app bot instance
			newBotAppOptions := []func(*BotApp) error{}

			TLSCertPath, TLSCertPathErr := env.GetTelegramTLSCertPath()
			TLSKeyPath, TLSKeyPathErr := env.GetTelegramTLSKeyPath()
			if TLSCertPathErr == nil && TLSKeyPathErr == nil {
				newBotAppOptions = append(
					newBotAppOptions,
					WebhookTLS(TLSCertPath, TLSKeyPath),
				)
			}

			storage := storage.NewSQLStorage(db)
			botApp, err := NewBotApp(
				storage,
				client,
				newBotAppOptions...,
			)
			if err != nil {
				log.Fatal(err)
			}
			if err := StartBotApp(botApp); err != nil {
				log.Fatal(err)
			}
		},
	}
}
