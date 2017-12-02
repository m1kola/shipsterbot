package env

import "os"

func GetDBConnectionString() string {
	return os.Getenv("DATABASE_URL")
}

func IsDebug() bool {
	return os.Getenv("DEBUG") == "true"
}
