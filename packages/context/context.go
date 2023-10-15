package context

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go-gin-clean-arch/packages/errors"
	"gorm.io/gorm"
)

type Context interface {
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
	requestContext context.Context
	validationErr  *errors.Error
	getDB          func() *gorm.DB
	db             *gorm.DB
	uid            xid.ID
}

func New(c *gin.Context, getDB func() *gorm.DB) Context {
	var uid xid.ID
	claimsInterface, ok := c.Get("claims")
	if ok {
		if uidInterface, ok := claimsInterface.(map[string]interface{})["uid"]; ok {
			uid, _ = xid.FromString(uidInterface.(string))
		}
	}

	return &ctx{
		requestContext: c.Request.Context(),
		validationErr:  errors.NewValidation(),
		getDB:          getDB,
		uid:            uid,
	}
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
