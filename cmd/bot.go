package cmd

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/m1kola/shipsterbot/bot"
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

		// Initialise bot instance
		tgbot, err := tgbotapi.NewBotAPI(env.GetTelegramAPIToken())
		tgbot.Debug = env.IsDebug()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Authorised on account %s", tgbot.Self.UserName)

		botApp := bot.TelegramBotApp{
			Bot:     tgbot,
			Storage: storage.NewSQLStorage(db)}
		botApp.ListenForWebhook()
		log.Fatal(
			http.ListenAndServeTLS(
				":8443",
				env.GetTelegramTLSCertPath(),
				env.GetTelegramTLSKeyPath(),
				incommingRequstLogger(http.DefaultServeMux)))
	},
}

func incommingRequstLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
