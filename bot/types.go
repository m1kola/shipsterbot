package bot

// handlerCanNotHandleError represents errors that were
type handlerCanNotHandleError struct {
	error
}

func (e handlerCanNotHandleError) Error() string {
	return e.error.Error()
}
