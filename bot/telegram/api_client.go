package telegram

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var tgbotapiNewBotAPI = tgbotapi.NewBotAPI

// NewAPIClient creates a new APIClientWrapper instance
func NewAPIClient(token string) (*APIClient, error) {
	client, err := tgbotapiNewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &APIClient{&apiClientWrapper{client}}, nil
}

// APIClient is a wrapper around the tgbotapi
// that incapsulates dependency from the code outside the telegram package.
// We shouldn't share tgbotapi's fields and methods directly
type APIClient struct {
	botClient *apiClientWrapper
}

// Token returns API token
func (api *APIClient) Token() string {
	return api.botClient.Token()
}

// SetDebug instructs the client to work in debug mode
func (api *APIClient) SetDebug(debug bool) {
	api.botClient.Debug = debug
}

// BotUserName returns bot's username
func (api *APIClient) BotUserName() string {
	return api.botClient.Self.UserName
}

// apiClientWrapper is required for internal porpuses,
// because inside the package we want to have access
// to tgbotapi's fields and methods.
// This also statisfies internal interfaces to allow us to
// write unit tests more easily.
type apiClientWrapper struct {
	*tgbotapi.BotAPI
}

func (api *apiClientWrapper) Token() string {
	return api.BotAPI.Token
}
