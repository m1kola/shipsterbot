package env

import (
	"os"
	"testing"
)

func TestSimpleEnvGetFunctions(t *testing.T) {
	testCases := []struct {
		funcToTest func() string
		testName   string
		envVarKey  string
	}{
		{
			funcToTest: GetTelegramAPIToken,
			testName:   "GetTelegramAPIToken",
			envVarKey:  "TELEGRAM_API_TOKEN",
		},
		{
			funcToTest: GetTelegramTLSCertPath,
			testName:   "GetTelegramTLSCertPath",
			envVarKey:  "TELEGRAM_TLS_CERT_PATH",
		},
		{
			funcToTest: GetTelegramTLSKeyPath,
			testName:   "GetTelegramTLSKeyPath",
			envVarKey:  "TELEGRAM_TLS_KEY_PATH",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			if envVarValue := os.Getenv(testCase.envVarKey); envVarValue != "" {
				os.Unsetenv(testCase.envVarKey)
				// Restore env var at the end
				defer func() { os.Setenv(testCase.envVarKey, envVarValue) }()
			}

			t.Run("Var is unset", func(t *testing.T) {
				v := testCase.funcToTest()

				if v == "" {
					return
				}

				t.Error("Function must not return a value, if the env var is unset.")
			})

			t.Run("Var is set", func(t *testing.T) {
				envVarTestValue := "some_val"
				os.Setenv(testCase.envVarKey, envVarTestValue)

				v := testCase.funcToTest()
				if v == envVarTestValue {
					return
				}

				t.Errorf("\"%s\" is expected, got \"%s\"", envVarTestValue, v)
			})
		})
	}
}

func TestGetTelegramWebhookPort(t *testing.T) {
	envVarKey := "TELEGRAM_WEBHOOK_PORT"

	if envVarValue := os.Getenv(envVarKey); envVarValue != "" {
		os.Unsetenv(envVarKey)
		// Restore env var at the end
		defer func() { os.Setenv(envVarKey, envVarValue) }()
	}

	t.Run("Var is unset", func(t *testing.T) {
		v, err := GetTelegramWebhookPort()

		if err != nil {
			t.Error("Got an unexpected error")
		}
		if v != "8443" {
			t.Error("Function must not return a value, if the env var is unset.")
		}
	})

	t.Run("Var is set", func(t *testing.T) {
		tests := []struct {
			value  string
			result string
		}{
			{value: "443", result: "443"},
			{value: "80", result: "80"},
			{value: "88", result: "88"},
			{value: "8443", result: "8443"},
		}

		for _, test := range tests {
			os.Setenv(envVarKey, test.value)

			v, err := GetTelegramWebhookPort()
			if err != nil {
				t.Error("Got an unexpected error")
			}
			if v != test.result {
				t.Errorf("%v expected, got %v", test.result, v)
			}
		}
	})

	t.Run("Var not allowed", func(t *testing.T) {
		os.Setenv(envVarKey, "values_is_not_allowed")
		_, err := GetTelegramWebhookPort()

		if err == nil {
			t.Error("Expected an error")
		}
	})
}
