package telegram

import "errors"

// updateRoutingError represents errors that can happen during update routing
type updateRoutingError struct {
	error
}

var errCommandIsNotSupported = updateRoutingError{
	errors.New("Unable to find a handler for a command")}
