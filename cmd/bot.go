package cmd

import (
	"database/sql"
	"log"

	"github.com/spf13/cobra"

	"github.com/m1kola/shipsterbot/bot/telegram"
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
		db, err := sql.Open("postgres", env.GetDBConnectionString())
		if err != nil {
			log.Fatal(err)
		}

		// Initialise a bot instance
		apiClient, err := telegram.NewAPIClient(env.GetTelegramAPIToken())
		if err != nil {
			log.Fatal(err)
		}
		apiClient.SetDebug(env.IsDebug())
		log.Printf("Authorised on account %s", apiClient.BotUserName())

		port, err := env.GetTelegramWebhookPort()
		if err != nil {
			log.Fatal(err)
		}
		TLSCertPath := env.GetTelegramTLSCertPath()
		TLSKeyPath := env.GetTelegramTLSKeyPath()

		// Create a app bot instance
		storage := storage.NewSQLStorage(db)
		botApp := telegram.NewBotApp(
			apiClient, storage,
			port, TLSCertPath, TLSKeyPath)
		err = telegram.StartBotApp(botApp)
		if err != nil {
			log.Fatal(err)
		}
	},
}
