package env

import "os"

func GetTelegramAPIToken() string {
	return os.Getenv("TELEGRAM_API_TOKEN")
}

func GetTelegramTLSCertPath() string {
	return os.Getenv("TELEGRAM_TLS_CERT_PATH")
}

func GetTelegramTLSKeyPath() string {
	return os.Getenv("TELEGRAM_TLS_KEY_PATH")
}
