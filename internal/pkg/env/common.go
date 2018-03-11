package env

import "os"

// GetDBConnectionString returns a string representing DB connection
func GetDBConnectionString() string {
	return os.Getenv("DATABASE_URL")
}

// IsDebug returns a bool indicating current execution mode
func IsDebug() bool {
	return os.Getenv("DEBUG") == "true"
}
