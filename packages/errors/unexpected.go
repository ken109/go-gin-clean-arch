package errors

import (
	"encoding/json"
	"fmt"
)

type Unexpected struct {
	message  string
	original string
	stack    []Frame
}

func NewUnexpected(err error) *Error {
	return &Error{
		kind: KindUnexpected,
		unexpected: &Unexpected{
			message:  err.Error(),
			original: fmt.Sprintf("%+v", err),
			stack:    callers(),
		},
	}
}

func (e Unexpected) Error() string {
	return fmt.Sprintf("message: %s, stack: %s", e.message, e.stack)
}

func (e Unexpected) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message  string  `json:"message"`
		Original string  `json:"original"`
		Stack    []Frame `json:"stack"`
	}{
		Message:  e.message,
		Original: e.original,
		Stack:    e.stack,
	})
}
