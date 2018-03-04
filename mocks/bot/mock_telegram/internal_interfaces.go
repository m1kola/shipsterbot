// Code generated by MockGen. DO NOT EDIT.
// Source: internal_interfaces.go

// Package mock_telegram is a generated GoMock package.
package mock_telegram

import (
	gomock "github.com/golang/mock/gomock"
	telegram_bot_api_v4 "gopkg.in/telegram-bot-api.v4"
	reflect "reflect"
)

// MockwebhookListener is a mock of webhookListener interface
type MockwebhookListener struct {
	ctrl     *gomock.Controller
	recorder *MockwebhookListenerMockRecorder
}

// MockwebhookListenerMockRecorder is the mock recorder for MockwebhookListener
type MockwebhookListenerMockRecorder struct {
	mock *MockwebhookListener
}

// NewMockwebhookListener creates a new mock instance
func NewMockwebhookListener(ctrl *gomock.Controller) *MockwebhookListener {
	mock := &MockwebhookListener{ctrl: ctrl}
	mock.recorder = &MockwebhookListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockwebhookListener) EXPECT() *MockwebhookListenerMockRecorder {
	return m.recorder
}

// ListenForWebhook mocks base method
func (m *MockwebhookListener) ListenForWebhook(pattern string) telegram_bot_api_v4.UpdatesChannel {
	ret := m.ctrl.Call(m, "ListenForWebhook", pattern)
	ret0, _ := ret[0].(telegram_bot_api_v4.UpdatesChannel)
	return ret0
}

// ListenForWebhook indicates an expected call of ListenForWebhook
func (mr *MockwebhookListenerMockRecorder) ListenForWebhook(pattern interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenForWebhook", reflect.TypeOf((*MockwebhookListener)(nil).ListenForWebhook), pattern)
}

// MocklistenerAndServer is a mock of listenerAndServer interface
type MocklistenerAndServer struct {
	ctrl     *gomock.Controller
	recorder *MocklistenerAndServerMockRecorder
}

// MocklistenerAndServerMockRecorder is the mock recorder for MocklistenerAndServer
type MocklistenerAndServerMockRecorder struct {
	mock *MocklistenerAndServer
}

// NewMocklistenerAndServer creates a new mock instance
func NewMocklistenerAndServer(ctrl *gomock.Controller) *MocklistenerAndServer {
	mock := &MocklistenerAndServer{ctrl: ctrl}
	mock.recorder = &MocklistenerAndServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MocklistenerAndServer) EXPECT() *MocklistenerAndServerMockRecorder {
	return m.recorder
}

// ListenAndServeTLS mocks base method
func (m *MocklistenerAndServer) ListenAndServeTLS(certFile, keyFile string) error {
	ret := m.ctrl.Call(m, "ListenAndServeTLS", certFile, keyFile)
	ret0, _ := ret[0].(error)
	return ret0
}

// ListenAndServeTLS indicates an expected call of ListenAndServeTLS
func (mr *MocklistenerAndServerMockRecorder) ListenAndServeTLS(certFile, keyFile interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAndServeTLS", reflect.TypeOf((*MocklistenerAndServer)(nil).ListenAndServeTLS), certFile, keyFile)
}

// ListenAndServe mocks base method
func (m *MocklistenerAndServer) ListenAndServe() error {
	ret := m.ctrl.Call(m, "ListenAndServe")
	ret0, _ := ret[0].(error)
	return ret0
}

// ListenAndServe indicates an expected call of ListenAndServe
func (mr *MocklistenerAndServerMockRecorder) ListenAndServe() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAndServe", reflect.TypeOf((*MocklistenerAndServer)(nil).ListenAndServe))
}

// Mocktokener is a mock of tokener interface
type Mocktokener struct {
	ctrl     *gomock.Controller
	recorder *MocktokenerMockRecorder
}

// MocktokenerMockRecorder is the mock recorder for Mocktokener
type MocktokenerMockRecorder struct {
	mock *Mocktokener
}

// NewMocktokener creates a new mock instance
func NewMocktokener(ctrl *gomock.Controller) *Mocktokener {
	mock := &Mocktokener{ctrl: ctrl}
	mock.recorder = &MocktokenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mocktokener) EXPECT() *MocktokenerMockRecorder {
	return m.recorder
}

// Token mocks base method
func (m *Mocktokener) Token() string {
	ret := m.ctrl.Call(m, "Token")
	ret0, _ := ret[0].(string)
	return ret0
}

// Token indicates an expected call of Token
func (mr *MocktokenerMockRecorder) Token() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Token", reflect.TypeOf((*Mocktokener)(nil).Token))
}

// Mocksender is a mock of sender interface
type Mocksender struct {
	ctrl     *gomock.Controller
	recorder *MocksenderMockRecorder
}

// MocksenderMockRecorder is the mock recorder for Mocksender
type MocksenderMockRecorder struct {
	mock *Mocksender
}

// NewMocksender creates a new mock instance
func NewMocksender(ctrl *gomock.Controller) *Mocksender {
	mock := &Mocksender{ctrl: ctrl}
	mock.recorder = &MocksenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mocksender) EXPECT() *MocksenderMockRecorder {
	return m.recorder
}

// Send mocks base method
func (m *Mocksender) Send(c telegram_bot_api_v4.Chattable) (telegram_bot_api_v4.Message, error) {
	ret := m.ctrl.Call(m, "Send", c)
	ret0, _ := ret[0].(telegram_bot_api_v4.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send
func (mr *MocksenderMockRecorder) Send(c interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*Mocksender)(nil).Send), c)
}

// MockcallbackQueryAnswerer is a mock of callbackQueryAnswerer interface
type MockcallbackQueryAnswerer struct {
	ctrl     *gomock.Controller
	recorder *MockcallbackQueryAnswererMockRecorder
}

// MockcallbackQueryAnswererMockRecorder is the mock recorder for MockcallbackQueryAnswerer
type MockcallbackQueryAnswererMockRecorder struct {
	mock *MockcallbackQueryAnswerer
}

// NewMockcallbackQueryAnswerer creates a new mock instance
func NewMockcallbackQueryAnswerer(ctrl *gomock.Controller) *MockcallbackQueryAnswerer {
	mock := &MockcallbackQueryAnswerer{ctrl: ctrl}
	mock.recorder = &MockcallbackQueryAnswererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockcallbackQueryAnswerer) EXPECT() *MockcallbackQueryAnswererMockRecorder {
	return m.recorder
}

// AnswerCallbackQuery mocks base method
func (m *MockcallbackQueryAnswerer) AnswerCallbackQuery(config telegram_bot_api_v4.CallbackConfig) (telegram_bot_api_v4.APIResponse, error) {
	ret := m.ctrl.Call(m, "AnswerCallbackQuery", config)
	ret0, _ := ret[0].(telegram_bot_api_v4.APIResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AnswerCallbackQuery indicates an expected call of AnswerCallbackQuery
func (mr *MockcallbackQueryAnswererMockRecorder) AnswerCallbackQuery(config interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AnswerCallbackQuery", reflect.TypeOf((*MockcallbackQueryAnswerer)(nil).AnswerCallbackQuery), config)
}

// MocktokenListenForWebhook is a mock of tokenListenForWebhook interface
type MocktokenListenForWebhook struct {
	ctrl     *gomock.Controller
	recorder *MocktokenListenForWebhookMockRecorder
}

// MocktokenListenForWebhookMockRecorder is the mock recorder for MocktokenListenForWebhook
type MocktokenListenForWebhookMockRecorder struct {
	mock *MocktokenListenForWebhook
}

// NewMocktokenListenForWebhook creates a new mock instance
func NewMocktokenListenForWebhook(ctrl *gomock.Controller) *MocktokenListenForWebhook {
	mock := &MocktokenListenForWebhook{ctrl: ctrl}
	mock.recorder = &MocktokenListenForWebhookMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MocktokenListenForWebhook) EXPECT() *MocktokenListenForWebhookMockRecorder {
	return m.recorder
}

// ListenForWebhook mocks base method
func (m *MocktokenListenForWebhook) ListenForWebhook(pattern string) telegram_bot_api_v4.UpdatesChannel {
	ret := m.ctrl.Call(m, "ListenForWebhook", pattern)
	ret0, _ := ret[0].(telegram_bot_api_v4.UpdatesChannel)
	return ret0
}

// ListenForWebhook indicates an expected call of ListenForWebhook
func (mr *MocktokenListenForWebhookMockRecorder) ListenForWebhook(pattern interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenForWebhook", reflect.TypeOf((*MocktokenListenForWebhook)(nil).ListenForWebhook), pattern)
}

// Token mocks base method
func (m *MocktokenListenForWebhook) Token() string {
	ret := m.ctrl.Call(m, "Token")
	ret0, _ := ret[0].(string)
	return ret0
}

// Token indicates an expected call of Token
func (mr *MocktokenListenForWebhookMockRecorder) Token() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Token", reflect.TypeOf((*MocktokenListenForWebhook)(nil).Token))
}

// MockbotClientInterface is a mock of botClientInterface interface
type MockbotClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockbotClientInterfaceMockRecorder
}

// MockbotClientInterfaceMockRecorder is the mock recorder for MockbotClientInterface
type MockbotClientInterfaceMockRecorder struct {
	mock *MockbotClientInterface
}

// NewMockbotClientInterface creates a new mock instance
func NewMockbotClientInterface(ctrl *gomock.Controller) *MockbotClientInterface {
	mock := &MockbotClientInterface{ctrl: ctrl}
	mock.recorder = &MockbotClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockbotClientInterface) EXPECT() *MockbotClientInterfaceMockRecorder {
	return m.recorder
}

// ListenForWebhook mocks base method
func (m *MockbotClientInterface) ListenForWebhook(pattern string) telegram_bot_api_v4.UpdatesChannel {
	ret := m.ctrl.Call(m, "ListenForWebhook", pattern)
	ret0, _ := ret[0].(telegram_bot_api_v4.UpdatesChannel)
	return ret0
}

// ListenForWebhook indicates an expected call of ListenForWebhook
func (mr *MockbotClientInterfaceMockRecorder) ListenForWebhook(pattern interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenForWebhook", reflect.TypeOf((*MockbotClientInterface)(nil).ListenForWebhook), pattern)
}

// Token mocks base method
func (m *MockbotClientInterface) Token() string {
	ret := m.ctrl.Call(m, "Token")
	ret0, _ := ret[0].(string)
	return ret0
}

// Token indicates an expected call of Token
func (mr *MockbotClientInterfaceMockRecorder) Token() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Token", reflect.TypeOf((*MockbotClientInterface)(nil).Token))
}

// Send mocks base method
func (m *MockbotClientInterface) Send(c telegram_bot_api_v4.Chattable) (telegram_bot_api_v4.Message, error) {
	ret := m.ctrl.Call(m, "Send", c)
	ret0, _ := ret[0].(telegram_bot_api_v4.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send
func (mr *MockbotClientInterfaceMockRecorder) Send(c interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockbotClientInterface)(nil).Send), c)
}

// AnswerCallbackQuery mocks base method
func (m *MockbotClientInterface) AnswerCallbackQuery(config telegram_bot_api_v4.CallbackConfig) (telegram_bot_api_v4.APIResponse, error) {
	ret := m.ctrl.Call(m, "AnswerCallbackQuery", config)
	ret0, _ := ret[0].(telegram_bot_api_v4.APIResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AnswerCallbackQuery indicates an expected call of AnswerCallbackQuery
func (mr *MockbotClientInterfaceMockRecorder) AnswerCallbackQuery(config interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AnswerCallbackQuery", reflect.TypeOf((*MockbotClientInterface)(nil).AnswerCallbackQuery), config)
}
