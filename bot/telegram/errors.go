package telegram

import "errors"

// handlerCanNotHandleError represents errors that were
type handlerCanNotHandleError struct {
	error
}

var errCommandIsNotSupported = handlerCanNotHandleError{
	errors.New("Unable to find a handler for a command")}
