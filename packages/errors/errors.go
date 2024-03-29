package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Kind string

const (
	KindUnexpected Kind = "unexpected"
	KindExpected   Kind = "expected"
	KindValidation Kind = "validation"
)

type Error struct {
	kind       Kind
	unexpected *Unexpected
	expected   *Expected
	validation *Validation
}

func (e Error) Error() string {
	message := fmt.Sprintf("%s error: ", e.kind)
	switch e.kind {
	case KindUnexpected:
		message += e.unexpected.message
	case KindExpected:
		message += e.expected.Error()
	case KindValidation:
		message += e.validation.Error()
	}
	return message
}

func (e Error) MarshalJSON() ([]byte, error) {
	switch e.kind {
	case KindUnexpected:
		if gin.Mode() == gin.ReleaseMode {
			return json.Marshal(struct {
				Kind string `json:"kind"`
			}{
				Kind: string(e.kind),
			})
		} else {
			return json.Marshal(struct {
				Kind       string      `json:"kind"`
				Unexpected *Unexpected `json:"unexpected"`
			}{
				Kind:       string(e.kind),
				Unexpected: e.unexpected,
			})
		}
	case KindExpected:
		return json.Marshal(struct {
			Kind     string    `json:"kind"`
			Expected *Expected `json:"expected"`
		}{
			Kind:     string(e.kind),
			Expected: e.expected,
		})
	case KindValidation:
		a, err := json.Marshal(struct {
			Kind       string      `json:"kind"`
			Validation *Validation `json:"validation"`
		}{
			Kind:       string(e.kind),
			Validation: e.validation,
		})
		return a, err
	}
	return nil, errors.New("エラー種別が正しくありません。")
}

func (e Error) Kind() Kind {
	return e.kind
}

func (e Error) Unexpected() error {
	return e.unexpected
}

func (e Error) Expected() *Expected {
	return e.expected
}

func (e Error) Validation() *Validation {
	return e.validation
}

func (e Error) Response() ErrorResponse {
	var statusCode int

	switch e.kind {
	case KindUnexpected:
		statusCode = http.StatusInternalServerError
	case KindExpected:
		statusCode = e.Expected().statusCode
	case KindValidation:
		statusCode = http.StatusBadRequest
	}

	return ErrorResponse{
		status: statusCode,
		err:    e,
	}
}

type ErrorResponse struct {
	status int
	err    Error
}

func (r ErrorResponse) Do(c *gin.Context) {
	c.PureJSON(r.status, struct {
		RequestID string `json:"request_id"`
		Error
	}{
		RequestID: requestid.Get(c),
		Error:     r.err,
	})
}
