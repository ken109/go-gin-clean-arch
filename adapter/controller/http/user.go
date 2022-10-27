package http

import (
	"github.com/gin-gonic/gin"
	"go-gin-ddd/adapter/presenter"
	"go-gin-ddd/config"
	"go-gin-ddd/packages/http/middleware"
	"go-gin-ddd/packages/http/router"
	"go-gin-ddd/usecase"

	"go-gin-ddd/packages/context"

	"go-gin-ddd/resource/request"
)

type user struct {
	inputFactory  usecase.UserInputFactory
	outputFactory func(c *gin.Context) usecase.UserOutputPort
}

func NewUser(r *router.Router, inputFactory usecase.UserInputFactory, outputFactory presenter.UserOutputFactory) {
	handler := user{
		inputFactory:  inputFactory,
		outputFactory: outputFactory,
	}

	r.Group("user", nil, func(r *router.Router) {
		r.Post("", handler.Create)
		r.Post("login", handler.Login)
		r.Post("refresh-token", handler.RefreshToken)
		r.Patch("reset-password-request", handler.ResetPasswordRequest)
		r.Patch("reset-password", handler.ResetPassword)
	})

	r.Group("", []gin.HandlerFunc{middleware.Auth(true, config.UserRealm, true)}, func(r *router.Router) {
		r.Group("user", nil, func(r *router.Router) {
			r.Get("me", handler.GetMe)
		})
	})
}

func (u user) Create(ctx context.Context, c *gin.Context) error {
	var req request.UserCreate

	if !bind(c, &req) {
		return nil
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.Create(ctx, &req)
}

func (u user) ResetPasswordRequest(ctx context.Context, c *gin.Context) error {
	var req request.UserResetPasswordRequest

	if !bind(c, &req) {
		return nil
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.ResetPasswordRequest(ctx, &req)
}

func (u user) ResetPassword(ctx context.Context, c *gin.Context) error {
	var req request.UserResetPassword

	if !bind(c, &req) {
		return nil
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.ResetPassword(ctx, &req)
}

func (u user) Login(ctx context.Context, c *gin.Context) error {
	var req request.UserLogin

	if !bind(c, &req) {
		return nil
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.Login(ctx, &req)
}

func (u user) RefreshToken(_ context.Context, c *gin.Context) error {
	var req request.UserRefreshToken

	if !bind(c, &req) {
		return nil
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.RefreshToken(&req)
}

func (u user) GetMe(ctx context.Context, c *gin.Context) error {
	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.GetByID(ctx, ctx.UID())
}
