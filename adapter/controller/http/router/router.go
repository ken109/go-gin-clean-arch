package router

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"gorm.io/gorm"

	"go-gin-clean-arch/config"
	"go-gin-clean-arch/packages/cerrors"
)

type Router struct {
	g  *gin.RouterGroup
	db *gorm.DB
}

func New(engine *gin.Engine, db *gorm.DB) *Router {
	return &Router{
		g:  engine.Group(""),
		db: db,
	}
}

func (r *Router) Group(relativePath string, handlers []gin.HandlerFunc, fn func(r *Router)) {
	if handlers == nil {
		handlers = []gin.HandlerFunc{}
	}
	fn(&Router{
		g:  r.g.Group(relativePath, handlers...),
		db: r.db,
	})
}

type HandlerFunc func(ctx context.Context, c *gin.Context) error

func (r *Router) Get(relativePath string, handlerFunc HandlerFunc) {
	r.g.GET(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Post(relativePath string, handlerFunc HandlerFunc) {
	r.g.POST(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Put(relativePath string, handlerFunc HandlerFunc) {
	r.g.PUT(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Patch(relativePath string, handlerFunc HandlerFunc) {
	r.g.PATCH(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Delete(relativePath string, handlerFunc HandlerFunc) {
	r.g.DELETE(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Options(relativePath string, handlerFunc HandlerFunc) {
	r.g.OPTIONS(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) Head(relativePath string, handlerFunc HandlerFunc) {
	r.g.HEAD(relativePath, r.wrapperFunc(handlerFunc))
}

func (r *Router) wrapperFunc(handlerFunc HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), config.DBKey, r.db)

		var uid xid.ID
		claimsInterface, ok := c.Get("claims")
		if ok {
			if uidInterface, ok := claimsInterface.(map[string]interface{})["uid"]; ok {
				uid, _ = xid.FromString(uidInterface.(string))
			}
		}
		ctx = context.WithValue(ctx, config.UIDKey, uid)
		ctx = context.WithValue(ctx, config.ErrorKey, cerrors.NewValidation())

		err := handlerFunc(ctx, c)
		if err != nil {
			var v *cerrors.Error
			switch {
			case errors.As(err, &v):
				v.Response(c)
			default:
				cerrors.NewUnexpected(v).Response(c)
			}

			_ = c.Error(err)
		}
	}
}
