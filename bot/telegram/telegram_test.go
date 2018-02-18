package telegram

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/m1kola/shipsterbot/mocks"
)

func TestNewBotApp(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	storageMock := mocks.NewMockDataStorageInterface(mockCtrl)
	APIClientMock := &APIClient{}

	expectedPort := "8443"
	expectedCert := "/fake/cert.pem"
	expectedKey := "/fake/key.key"

	app := NewBotApp(APIClientMock, storageMock,
		expectedPort, expectedCert, expectedKey)

	if app.serverConfig.port != expectedPort {
		t.Errorf("%s expected as a port, got %s",
			expectedPort, app.serverConfig.port)
	}

	if app.serverConfig.TLSCertPath != expectedCert {
		t.Errorf("%s expected as TLSCertPath, got %s",
			expectedCert, app.serverConfig.TLSCertPath)
	}

	if app.serverConfig.TLSKeyPath != expectedKey {
		t.Errorf("%s expected as TLSKeyPath, got %s",
			expectedKey, app.serverConfig.TLSKeyPath)
	}
}
