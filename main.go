package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/m1kola/shipsterbot/bot"
	"github.com/m1kola/shipsterbot/storage"
)

func getDBConnectionString() string {
	return os.Getenv("DATABASE_URL")
}

func getAPIToken() string {
	return os.Getenv("TELEGRAM_API_TOKEN")
}

func isDebug() bool {
	return os.Getenv("DEBUG") == "true"
}

func incommingRequstLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	// Initialise DB connection pool
	db, err := sql.Open("postgres", getDBConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	// Initialise bot instance
	tgbot, err := tgbotapi.NewBotAPI(getAPIToken())
	tgbot.Debug = isDebug()
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
			os.Getenv("TLS_CERT_PATH"),
			os.Getenv("TLS_KEY_PATH"),
			incommingRequstLogger(http.DefaultServeMux)))
}
