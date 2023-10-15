package http

import (
	"github.com/gin-gonic/gin"

	"go-gin-clean-arch/adapter/presenter"
	"go-gin-clean-arch/config"
	"go-gin-clean-arch/packages/context"
	"go-gin-clean-arch/packages/http/middleware"
	"go-gin-clean-arch/packages/http/router"
	"go-gin-clean-arch/resource/request"
	"go-gin-clean-arch/usecase"
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

	r.Group("users", nil, func(r *router.Router) {
		r.Post("", handler.Create)
		r.Post("login", handler.Login)
		r.Post("refresh-token", handler.RefreshToken)
		r.Patch("reset-password-request", handler.ResetPasswordRequest)
		r.Patch("reset-password", handler.ResetPassword)

		r.Group("", []gin.HandlerFunc{middleware.Auth(true, config.UserRealm, true)}, func(r *router.Router) {
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
