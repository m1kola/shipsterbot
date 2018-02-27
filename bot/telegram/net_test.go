package telegram

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/m1kola/shipsterbot/mocks/bot/mock_telegram"
)

func TestValidateWebhookPort(t *testing.T) {
	t.Run("Valid port", func(t *testing.T) {
		tests := []string{"443", "80", "88", "8443"}

		for _, port := range tests {
			err := ValidateWebhookPort(port)
			if err != nil {
				t.Errorf("Got an unexpected error: %s is a valid value", port)
			}
		}
	})

	t.Run("Invalid port", func(t *testing.T) {
		tests := []string{"", "-1", "0", "1234"}

		for _, port := range tests {
			err := ValidateWebhookPort(port)
			if err == nil {
				t.Errorf("Expected an error, got positive result: %s is invalid value",
					port)
			}
		}
	})
}

func TestNewServerWithIncommingRequstLogger(t *testing.T) {
	expectedAddr := ":8443"

	server := newServerWithIncommingRequstLogger("8443", http.DefaultServeMux)
	defer server.Close()

	if server.Addr != expectedAddr {
		t.Errorf("Expected addr is %s, got %s",
			expectedAddr, server.Addr)
	}
}

func TestListenAndServe(t *testing.T) {
	t.Run("TLS server", func(t *testing.T) {
		TLSCertPath, TLSKeyPath := "/test/cert.pem", "/test/cert.key"

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		fakeServer := mock_telegram.NewMocklistenerAndServer(mockCtrl)
		fakeServer.EXPECT().ListenAndServeTLS(TLSCertPath, TLSKeyPath)

		listenAndServe(fakeServer, TLSCertPath, TLSKeyPath)
	})
	t.Run("Non-TLS server", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		fakeServer := mock_telegram.NewMocklistenerAndServer(mockCtrl)
		fakeServer.EXPECT().ListenAndServe()

		listenAndServe(fakeServer, "", "")
	})
}

func TestGetUpdatesChan(t *testing.T) {
	expectedPattern := "/123/webhook"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock_telegram.NewMocktokenListenForWebhook(mockCtrl)
	tokenCall := mockObj.EXPECT().Token()
	tokenCall.Return("123")

	mockObj.EXPECT().ListenForWebhook(expectedPattern).After(tokenCall)

	getUpdatesChan(mockObj)

}

func TestIncommingRequstLogger(t *testing.T) {
	// Setup capturing buffer and restoer previous output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() { log.SetOutput(os.Stderr) }()

	originalIsCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalIsCalled = true
	})

	newHandler := incommingRequstLogger(mockHandler)

	mockRequest := httptest.NewRequest("GET", "/traget-url", nil)
	w := httptest.NewRecorder()
	newHandler.ServeHTTP(w, mockRequest)

	if !originalIsCalled {
		t.Error("Original handler expected to be called")
	}

	expectedLog := "192.0.2.1:1234 GET /traget-url"
	bufString := buf.String()
	if !strings.Contains(bufString, expectedLog) {
		t.Errorf("%s expected to contain %s", bufString, expectedLog)
	}
}
