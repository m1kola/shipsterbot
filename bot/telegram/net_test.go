package telegram

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func TestGetUpdatesChan(t *testing.T) {
	var actualPattern string
	expectedPattern := "/123/webhook"

	mock := &tokenListenForWebhookMock{
		&webhookListenerMock{
			webhookListenerFunc: func(pattern string) tgbotapi.UpdatesChannel {
				actualPattern = pattern
				return make(tgbotapi.UpdatesChannel)
			},
		},
		&tokenerMock{
			fakeToken: "123",
		},
	}

	getUpdatesChan(mock)

	if actualPattern != expectedPattern {
		t.Errorf("%s expected, got %s", expectedPattern, actualPattern)
	}
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

	expectedSuffix := "192.0.2.1:1234 GET /traget-url"
	bufString := buf.String()
	if !strings.Contains(bufString, expectedSuffix) {
		t.Errorf("%s expected to contain %s suffix", bufString, expectedSuffix)
	}
}
