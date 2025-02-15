package cerrors

import (
	"encoding/json"
	"fmt"
)

type Unexpected struct {
	message string
	stack   []Frame
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
			message: err.Error(),
			stack:   callers(callerSkip),
		},
	}
}

func (e Unexpected) Error() string {
	return fmt.Sprintf("message: %s, stack: %s", e.message, e.stack)
}

func (e Unexpected) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string  `json:"message"`
		Stack   []Frame `json:"stack"`
	}{
		Message: e.message,
		Stack:   e.stack,
	})
}
