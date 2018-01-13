package env

import (
	"fmt"
	"os"
	"strings"
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
	allowedPorts := []string{"443", "80", "88", "8443"}
	value, ok := os.LookupEnv("TELEGRAM_WEBHOOK_PORT")

	if !ok {
		return "8443", nil
	}

	for _, allowedValue := range allowedPorts {
		if value == allowedValue {
			return value, nil
		}
	}

	err := fmt.Errorf(
		"Wrong port. You can only use one of the following ports: %s",
		strings.Join(allowedPorts, ", "))
	return "", err
}
