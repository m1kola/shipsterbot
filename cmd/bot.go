package cmd

import (
	"database/sql"
	"log"

	"github.com/spf13/cobra"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/m1kola/shipsterbot/bot/telegram"
	"github.com/m1kola/shipsterbot/env"
	"github.com/m1kola/shipsterbot/storage"
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
		tgbot, err := tgbotapi.NewBotAPI(env.GetTelegramAPIToken())
		if err != nil {
			log.Fatal(err)
		}
		tgbot.Debug = env.IsDebug()
		log.Printf("Authorised on account %s", tgbot.Self.UserName)

		port, err := env.GetTelegramWebhookPort()
		if err != nil {
			log.Fatal(err)
		}
		TLSCertPath := env.GetTelegramTLSCertPath()
		TLSKeyPath := env.GetTelegramTLSKeyPath()

		// Create a app bot instance
		storage := storage.NewSQLStorage(db)
		botApp := telegram.NewBotApp(tgbot, storage, port, TLSCertPath, TLSKeyPath)
		err = botApp.Start()
		if err != nil {
			log.Fatal(err)
		}
	},
}
