package cli

import (
	"database/sql"
	"log"

	"github.com/spf13/cobra"

	"github.com/m1kola/shipsterbot/internal/bot/telegram"
	"github.com/m1kola/shipsterbot/internal/pkg/env"
	"github.com/m1kola/shipsterbot/internal/pkg/storage"
)

func init() {
	rootCmd.AddCommand(startBotCmd)
	startBotCmd.AddCommand(startTelegramBotCmd)
}

var startBotCmd = &cobra.Command{
	Use:   "startbot",
	Short: "Start a bot",
}

var startTelegramBotCmd = &cobra.Command{
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

		// Create a app bot instance
		newBotAppOptions := []func(*telegram.BotApp) error{}

		TLSCertPath, TLSCertPathErr := env.GetTelegramTLSCertPath()
		TLSKeyPath, TLSKeyPathErr := env.GetTelegramTLSKeyPath()
		if TLSCertPathErr == nil && TLSKeyPathErr == nil {
			newBotAppOptions = append(
				newBotAppOptions,
				telegram.WebhookTLS(TLSCertPath, TLSKeyPath),
			)
		}

		storage := storage.NewSQLStorage(db)
		botApp, err := telegram.NewBotApp(
			storage,
			apiToken,
			newBotAppOptions...,
		)
		if err != nil {
			log.Fatal(err)
		}
		if err := telegram.StartBotApp(botApp); err != nil {
			log.Fatal(err)
		}
	},
}
