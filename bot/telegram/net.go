package telegram

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// ValidateWebhookPort returns an error if an invalid port has been provided
func ValidateWebhookPort(port string) error {
	allowedPorts := []string{"443", "80", "88", "8443"}

	for _, allowedValue := range allowedPorts {
		if port == allowedValue {
			return nil
		}
	}

	err := fmt.Errorf(
		"Wrong port. You can only use one of the following ports: %s",
		strings.Join(allowedPorts, ", "))
	return err
}

// startWebhookServer starts a new http server for handling Telegram webhooks
func startWebhookServer(port, TLSCertPath, TLSKeyPath string) error {
	handler := incommingRequstLogger(http.DefaultServeMux)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Start listening on %s", addr)

	var err error
	if len(TLSCertPath) > 0 && len(TLSKeyPath) > 0 {
		err = http.ListenAndServeTLS(addr, TLSCertPath, TLSKeyPath, handler)
	} else {
		err = http.ListenAndServe(addr, handler)
	}

	return err
}

func incommingRequstLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// getUpdatesChan regesters a webhook handler
// and return a channel for consuming updates
func getUpdatesChan(bot tokenListenForWebhook) <-chan tgbotapi.Update {
	return bot.ListenForWebhook(
		fmt.Sprintf("/%s/webhook", bot.Token()))
}
