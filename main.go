package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

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
	bot, err := tgbotapi.NewBotAPI(getAPIToken())
	bot.Debug = isDebug()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorised on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook(fmt.Sprintf("/%s/webhook", bot.Token))
	go HandleUpdates(bot, updates)
	log.Fatal(
		http.ListenAndServeTLS(
			":8443",
			os.Getenv("TLS_CERT_PATH"),
			os.Getenv("TLS_KEY_PATH"),
			incommingRequstLogger(http.DefaultServeMux)))
}
