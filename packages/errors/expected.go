package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Expected struct {
	statusCode int
	msg        string
}

func NewExpected(statusCode int, message string) *Error {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	return &Error{
		kind:     KindExpected,
		expected: &Expected{statusCode: statusCode, msg: message},
	}
}

func (e Expected) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.statusCode, e.msg)
}

func (e Expected) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: e.msg,
	})
}

func (e Expected) StatusCode() int {
	return e.statusCode
}

func (e Expected) Message() string {
	return e.msg
}

func (e *Expected) ChangeStatus(before int, after int) bool {
	if e.statusCode == before {
		e.statusCode = after
		return true
	}
	return false
}

func (e Expected) StatusOk() bool {
	return e.statusCode < 300
}

// expected errors

func NotFound() *Error {
	return NewExpected(http.StatusNotFound, "リソースが見つかりませんでした")
}

func Forbidden() *Error {
	return NewExpected(http.StatusForbidden, "リソースに対する権限がありません")
}
