// Package env encapsulates all available env vars that the program can accept
package env

import "os"

// Common env vars
const (
	dbConnectionStringVarName = "DATABASE_URL"
	debugVarName              = "DEBUG"
)

// Telegram bot specific env vars
const (
	telegramAPITokenVarName    = "TELEGRAM_API_TOKEN"
	telegramTLSCertPathVarName = "TELEGRAM_TLS_CERT_PATH"
	telegramTLSKeyPathVarName  = "TELEGRAM_TLS_KEY_PATH"
	telegramWebhookPortVarName = "TELEGRAM_WEBHOOK_PORT"
)

// IsDebug returns a bool indicating current execution mode
func IsDebug() bool {
	return os.Getenv(debugVarName) == "true"
}

// GetDBConnectionString returns a string representing DB connection
func GetDBConnectionString() (string, error) {
	return lookupEnv(dbConnectionStringVarName)
}

// GetTelegramAPIToken returns telegram bot api token
func GetTelegramAPIToken() (string, error) {
	return lookupEnv(telegramAPITokenVarName)
}

// GetTelegramTLSCertPath returns path to TLS certificate
func GetTelegramTLSCertPath() (string, error) {
	return lookupEnv(telegramTLSCertPathVarName)
}

// GetTelegramTLSKeyPath returns path to TLS key
func GetTelegramTLSKeyPath() (string, error) {
	return lookupEnv(telegramTLSKeyPathVarName)
}

// GetTelegramWebhookPort returns port for Telegram webhook server
func GetTelegramWebhookPort() (string, error) {
	return lookupEnv(telegramWebhookPortVarName)
}
