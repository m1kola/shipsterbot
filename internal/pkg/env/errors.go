package env

import "fmt"

type envVarNotFoundError struct {
	envVarName string
}

func (err envVarNotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find the \"%s\" env var", err.envVarName)
}
