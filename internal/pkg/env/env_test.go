package env

import (
	"os"
	"testing"
)

func TestSimpleEnvLookupFunctions(t *testing.T) {
	testCases := []struct {
		funcToTest func() (string, error)
		testName   string
		envVarKey  string
	}{
		{
			funcToTest: GetDBConnectionString,
			testName:   "GetDBConnectionString",
			envVarKey:  dbConnectionStringVarName,
		},
		{
			funcToTest: GetTelegramAPIToken,
			testName:   "GetTelegramAPIToken",
			envVarKey:  telegramAPITokenVarName,
		},
		{
			funcToTest: GetTelegramTLSCertPath,
			testName:   "GetTelegramTLSCertPath",
			envVarKey:  telegramTLSCertPathVarName,
		},
		{
			funcToTest: GetTelegramTLSKeyPath,
			testName:   "GetTelegramTLSKeyPath",
			envVarKey:  telegramTLSKeyPathVarName,
		},
		{
			funcToTest: GetTelegramWebhookPort,
			testName:   "GetTelegramWebhookPort",
			envVarKey:  telegramWebhookPortVarName,
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
				_, err := testCase.funcToTest()

				if err == nil {
					t.Error("Function must return an error, when the env var is unset.")
				}

				if _, ok := err.(envVarNotFoundError); !ok {
					t.Errorf(
						"Function must return an error of type %T, when the env var is unset. Got %T",
						envVarNotFoundError{}, err,
					)
				}

			})

			t.Run("Var is set", func(t *testing.T) {
				envVarTestValue := "some_val"
				os.Setenv(testCase.envVarKey, envVarTestValue)

				v, err := testCase.funcToTest()

				if err != nil {
					t.Error("Function must not return an error, when the env var is set.")
				}

				if v != envVarTestValue {
					t.Errorf("%#v is expected, got %#v", envVarTestValue, v)
				}

			})
		})
	}
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
