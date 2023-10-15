package context

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"gorm.io/gorm"

	"go-gin-clean-arch/packages/errors"
)

type Context interface {
	RequestID() string
	RequestContext() context.Context
	Authenticated() bool
	UID() xid.ID

	Validate(request interface{}) (invalid bool)
	FieldError(fieldName string, message string)
	IsInValid() bool
	ValidationError() error

	DB() *gorm.DB
	Transaction(fn func(ctx Context) error) error
}

type ctx struct {
	requestID      string
	requestContext context.Context
	validationErr  *errors.Error
	getDB          func() *gorm.DB
	db             *gorm.DB
	uid            xid.ID
}

func New(c *gin.Context, getDB func() *gorm.DB) Context {
	requestID := c.GetHeader("X-Request-Id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	var uid xid.ID
	claimsInterface, ok := c.Get("claims")
	if ok {
		if uidInterface, ok := claimsInterface.(map[string]interface{})["uid"]; ok {
			uid, _ = xid.FromString(uidInterface.(string))
		}
	}

	return &ctx{
		requestID:      requestID,
		requestContext: c.Request.Context(),
		validationErr:  errors.NewValidation(),
		getDB:          getDB,
		uid:            uid,
	}
}

func (c *ctx) RequestID() string {
	return c.requestID
}

func (c *ctx) RequestContext() context.Context {
	return c.requestContext
}

func (c *ctx) Authenticated() bool {
	return !c.uid.IsNil()
}

func (c *ctx) UID() xid.ID {
	return c.uid
}
