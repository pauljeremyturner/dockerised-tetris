package server

import "fmt"

type ErrorInvalidRequest struct {
	msg   string
	cause error
}

func (e ErrorInvalidRequest) Error() string {
	if e.cause == nil {
		return fmt.Sprintf("invalid request, message: %s", e.msg)
	}
	return fmt.Sprintf("invalid request, message: %s, cause: %s", e.msg, e.cause.Error())
}

func NewErrorInvalidRequest(msg string, cause error) ErrorInvalidRequest {
	return ErrorInvalidRequest{
		msg:   msg,
		cause: cause,
	}
}
