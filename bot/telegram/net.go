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

// newServerWithIncommingRequstLogger creates a new server struct
// with an incomming request logger
func newServerWithIncommingRequstLogger(port string, handler http.Handler) *http.Server {
	newHandler := incommingRequstLogger(handler)
	addr := fmt.Sprintf(":%s", port)

	server := &http.Server{Addr: addr, Handler: newHandler}
	return server
}

// listenAndServe makes the server start handling requests.
// It serves TLS connectons, if paths for TLS cert and key are provided,
// othervise it serves non-TLS connections
func listenAndServe(server listenerAndServer, TLSCertPath, TLSKeyPath string) error {
	if len(TLSCertPath) > 0 && len(TLSKeyPath) > 0 {
		return server.ListenAndServeTLS(TLSCertPath, TLSKeyPath)
	}

	return server.ListenAndServe()
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
