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

func main() {
	bot, err := tgbotapi.NewBotAPI(getAPIToken())
	bot.Debug = isDebug()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorised on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook(fmt.Sprintf("/%s/webhook", bot.Token))
	go HandleUpdates(bot, updates)
	http.ListenAndServeTLS(
		":8443",
		os.Getenv("TLS_CERT_PATH"),
		os.Getenv("TLS_KEY_PATH"),
		nil)
}
