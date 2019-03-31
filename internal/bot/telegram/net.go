package telegram

import (
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// newServerWithIncommingRequstLogger creates a new server struct
// with an incomming request logger
func newServerWithIncommingRequstLogger(port string, handler http.Handler) *http.Server {
	newHandler := incommingRequstLogger(handler)
	addr := fmt.Sprintf(":%s", port)

	return &http.Server{Addr: addr, Handler: newHandler}
}

// listenAndServe makes the server start handling requests.
// It serves TLS connectons, if paths for TLS cert and key are provided,
// othervise it serves non-TLS connections
var listenAndServe = func(server listenerAndServer, TLSCertPath, TLSKeyPath string) error {
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
var getUpdatesChan = func(bot tokenListenForWebhook) <-chan tgbotapi.Update {
	return bot.ListenForWebhook(
		fmt.Sprintf("/%s/webhook", bot.Token()))
}
