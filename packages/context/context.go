package context

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"packages/errors"
)

type Context interface {
	RequestID() string
	Authenticated() bool
	UID() uint

	Validate(request interface{}) (invalid bool)
	FieldError(fieldName string, message string)
	IsInValid() bool
	ValidationError() error

	DB() *gorm.DB
	Transaction(fn func(ctx Context) error) error
}

type ctx struct {
	id    string
	verr  *errors.Error
	getDB func() *gorm.DB
	db    *gorm.DB
	uid   uint
}

func New(c *gin.Context, getDB func() *gorm.DB) Context {
	requestID := c.GetHeader("X-Request-Id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	var uid uint
	claimsInterface, ok := c.Get("claims")
	if ok {
		if uidInterface, ok := claimsInterface.(map[string]interface{})["uid"]; ok {
			uid = uint(uidInterface.(float64))
		}
	}

	return &ctx{
		id:    requestID,
		verr:  errors.NewValidation(),
		getDB: getDB,
		uid:   uid,
	}
}

func (c ctx) RequestID() string {
	return c.id
}

func (c ctx) Authenticated() bool {
	return c.uid != 0
}

func (c ctx) UID() uint {
	return c.uid
}
