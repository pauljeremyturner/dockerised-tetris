package client

import "fmt"

type ClientSystemError struct {
	Msg   string
	Cause error
}

func NewClientSystemError(msg string, cause error) ClientSystemError {
	return ClientSystemError{
		Msg:   msg,
		Cause: cause,
	}
}

func (e ClientSystemError) Error() string {
	if e.Cause == nil {
		return fmt.Sprintf("invalid request, message: %s", e.Msg)
	}
	return fmt.Sprintf("invalid request, message: %s, cause: %s", e.Msg, e.Cause.Error())
}
