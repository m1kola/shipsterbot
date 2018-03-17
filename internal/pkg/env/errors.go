package env

import "fmt"

type EnvVarNotFoundError struct {
	envVarName string
}

func (err EnvVarNotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find the \"%s\" env var", err.envVarName)
}
