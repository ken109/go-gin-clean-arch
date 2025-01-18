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

type UnexpectedOption interface {
	Kind() UnexpectedOptionKind
}

type UnexpectedOptionKind = string

const UnexpectedOptionKindPanic = "panic"

type WithUnexpectedPanic struct{}

func (w WithUnexpectedPanic) Kind() UnexpectedOptionKind {
	return UnexpectedOptionKindPanic
}

func NewUnexpected(err error, options ...UnexpectedOption) *Error {
	callerSkip := 3

	for _, option := range options {
		if option.Kind() == UnexpectedOptionKindPanic {
			callerSkip = 5
		}
	}

	return &Error{
		kind: KindUnexpected,
		unexpected: &Unexpected{
			message:  err.Error(),
			original: fmt.Sprintf("%+v", err),
			stack:    callers(callerSkip),
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
