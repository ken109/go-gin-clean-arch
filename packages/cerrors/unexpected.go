package cerrors

import (
	"encoding/json"
	"fmt"
)

type Unexpected struct {
	err   error
	stack []Frame
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
			err:   err,
			stack: callers(callerSkip),
		},
	}
}

func (e Unexpected) Error() string {
	return fmt.Sprintf("message: %s, stack: %s", e.err.Error(), e.stack)
}

func (e Unexpected) Unwrap() error {
	return e.err
}

func (e Unexpected) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string  `json:"message"`
		Stack   []Frame `json:"stack"`
	}{
		Message: e.err.Error(),
		Stack:   e.stack,
	})
}
