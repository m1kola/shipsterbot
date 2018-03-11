package env

import (
	"os"
	"testing"
)

func TestGetDBConnectionString(t *testing.T) {
	envVarKey := "DATABASE_URL"

	if envVarValue := os.Getenv(envVarKey); envVarValue != "" {
		os.Unsetenv(envVarKey)
		// Restore env var at the end
		defer func() { os.Setenv(envVarKey, envVarValue) }()
	}

	t.Run("Var is unset", func(t *testing.T) {
		v := GetDBConnectionString()

		if v == "" {
			return
		}

		t.Error("Function must not return a value, if the env var is unset.")
	})

	t.Run("Var is set", func(t *testing.T) {
		envVarTestValue := "postgres://localhost/dbname?sslmode=disable"
		os.Setenv(envVarKey, envVarTestValue)

		v := GetDBConnectionString()
		if v == envVarTestValue {
			return
		}

		t.Errorf("\"%s\" is expected, got \"%s\"", envVarTestValue, v)
	})
}

func TestIsDebug(t *testing.T) {
	envVarKey := "DEBUG"

	if envVarValue := os.Getenv(envVarKey); envVarValue != "" {
		os.Unsetenv(envVarKey)
		// Restore env var at the end
		defer func() { os.Setenv(envVarKey, envVarValue) }()
	}

	t.Run("Var is unset", func(t *testing.T) {
		v := IsDebug()

		if v == false {
			return
		}

		t.Error("Function must return false, if the env var is unset.")
	})

	t.Run("Var is set", func(t *testing.T) {
		tests := []struct {
			value  string
			result bool
		}{
			{value: "true", result: true},
			{value: "false", result: false},
			{value: "something", result: false},
			{value: "", result: false},
		}

		for _, test := range tests {
			os.Setenv(envVarKey, test.value)

			v := IsDebug()
			if v == test.result {
				continue
			}

			t.Errorf("%v expected, got %v", test.result, v)
		}
	})
}
