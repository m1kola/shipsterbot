package env

import (
	"os"

	"github.com/m1kola/shipsterbot/internal/bot/telegram"
)

// GetTelegramAPIToken returns telegram bot api token
func GetTelegramAPIToken() string {
	return os.Getenv("TELEGRAM_API_TOKEN")
}

// GetTelegramTLSCertPath returns path to TLS certificate
func GetTelegramTLSCertPath() string {
	return os.Getenv("TELEGRAM_TLS_CERT_PATH")
}

// GetTelegramTLSKeyPath returns path to TLS key
func GetTelegramTLSKeyPath() string {
	return os.Getenv("TELEGRAM_TLS_KEY_PATH")
}

// GetTelegramWebhookPort returns port for Telegram webhook server
// Default port is 8443
func GetTelegramWebhookPort() (string, error) {
	value, ok := os.LookupEnv("TELEGRAM_WEBHOOK_PORT")
	if !ok {
		return "8443", nil
	}

	err := telegram.ValidateWebhookPort(value)
	if err != nil {
		return "", err
	}
	return value, nil
}
