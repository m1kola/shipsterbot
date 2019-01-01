package telegram

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// ValidateWebhookPort returns an error if an invalid port has been provided
func ValidateWebhookPort(port string) error {
	allowedPorts := []string{"443", "80", "88", "8443"}

	for _, allowedValue := range allowedPorts {
		if port == allowedValue {
			return nil
		}
	}

	return fmt.Errorf(
		"Wrong port. You can only use one of the following ports: %s",
		strings.Join(allowedPorts, ", "))
}

// newServerWithincomingRequestLogger creates a new server struct
// with an incoming request logger
func newServerWithincomingRequestLogger(port string, handler http.Handler) *http.Server {
	newHandler := incomingRequestLogger(handler)
	addr := fmt.Sprintf(":%s", port)

	return &http.Server{Addr: addr, Handler: newHandler}
}

// listenAndServe makes the server start handling requests.
// It serves TLS connections, if paths for TLS cert and key are provided,
// othervise it serves non-TLS connections
var listenAndServe = func(server listenerAndServer, TLSCertPath, TLSKeyPath string) error {
	if len(TLSCertPath) > 0 && len(TLSKeyPath) > 0 {
		return server.ListenAndServeTLS(TLSCertPath, TLSKeyPath)
	}

	return server.ListenAndServe()
}

func incomingRequestLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// getUpdatesChan regesters a webhook handler
// and return a channel for consuming updates
func getUpdatesChan(client *tgbotapi.BotAPI) <-chan tgbotapi.Update {
	return client.ListenForWebhook(
		fmt.Sprintf("/%s/webhook", client.Token))
}
