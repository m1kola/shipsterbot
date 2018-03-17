package env

import (
	"os"
)

// lookupEnv is shortcut for env var lookup
func lookupEnv(envVarName string) (string, error) {
	val, ok := os.LookupEnv(envVarName)

	if !ok {
		return "", EnvVarNotFoundError{envVarName: envVarName}
	}
	return val, nil
}
