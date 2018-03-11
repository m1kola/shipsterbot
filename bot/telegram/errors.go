package telegram

import "errors"

// updateRoutingError represents errors that can happenen during update routing
type updateRoutingError struct {
	error
}

var errCommandIsNotSupported = updateRoutingError{
	errors.New("Unable to find a handler for a command")}
