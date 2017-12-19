package cmd

import (
	"database/sql"
	"fmt"
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

		// Initialise a bot instance
		tgbot, err := tgbotapi.NewBotAPI(env.GetTelegramAPIToken())
		if err != nil {
			log.Fatal(err)
		}
		tgbot.Debug = env.IsDebug()
		log.Printf("Authorised on account %s", tgbot.Self.UserName)

		// Create a app bot instance
		botApp := bot.TelegramBotApp{
			Bot:     tgbot,
			Storage: storage.NewSQLStorage(db)}
		botApp.ListenForWebhook()

		startWebServer()
	},
}

func incommingRequstLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func startWebServer() {
	var err error
	handler := incommingRequstLogger(http.DefaultServeMux)

	port, err := env.GetTelegramWebhookPort()
	if err != nil {
		log.Fatal(err)
	}
	TLSCertPath := env.GetTelegramTLSCertPath()
	TLSKeyPath := env.GetTelegramTLSKeyPath()
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Listening on %s", addr)
	if len(TLSCertPath) > 0 && len(TLSKeyPath) > 0 {
		err = http.ListenAndServeTLS(addr, TLSCertPath, TLSKeyPath, handler)
	} else {
		err = http.ListenAndServe(addr, handler)
	}

	log.Fatal(err)
}
