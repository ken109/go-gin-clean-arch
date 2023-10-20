package http

import (
	"github.com/gin-gonic/gin"

	"context"
	"github.com/gin-contrib/sessions"
	"go-gin-clean-arch/adapter/controller/http/middleware"
	"go-gin-clean-arch/adapter/controller/http/router"
	"go-gin-clean-arch/config"
	"go-gin-clean-arch/packages/util"
	"go-gin-clean-arch/resource/request"
	"go-gin-clean-arch/usecase"
	"net/http"
)

type user struct {
	usecase usecase.User
}

func NewUser(r *router.Router, usecase usecase.User) {
	handler := user{
		usecase: usecase,
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

	id, err := u.usecase.Create(ctx, &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, id)
	return nil
}

func (u user) ResetPasswordRequest(ctx context.Context, c *gin.Context) error {
	var req request.UserResetPasswordRequest

	if !bind(c, &req) {
		return nil
	}

	res, err := u.usecase.ResetPasswordRequest(ctx, &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, res)
	return nil
}

func (u user) ResetPassword(ctx context.Context, c *gin.Context) error {
	var req request.UserResetPassword

	if !bind(c, &req) {
		return nil
	}

	err := u.usecase.ResetPassword(ctx, &req)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (u user) Login(ctx context.Context, c *gin.Context) error {
	var req request.UserLogin

	if !bind(c, &req) {
		return nil
	}

	res, err := u.usecase.Login(ctx, &req)
	if err != nil {
		return err
	}

	if res == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	if req.Session {
		session := sessions.DefaultMany(c, config.UserRealm)
		session.Set("token", res.Token)
		session.Set("refresh_token", res.RefreshToken)
		if err := session.Save(); err != nil {
			return err
		}
		c.Status(http.StatusOK)
	} else {
		c.JSON(http.StatusOK, res)
	}

	return nil
}

func (u user) RefreshToken(_ context.Context, c *gin.Context) error {
	var req request.UserRefreshToken

	if !bind(c, &req) {
		return nil
	}

	res, err := u.usecase.RefreshToken(&req)
	if err != nil {
		return err
	}

	if res == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	if req.Session {
		session := sessions.DefaultMany(c, config.UserRealm)
		session.Set("token", res.Token)
		session.Set("refresh_token", res.RefreshToken)
		if err := session.Save(); err != nil {
			return err
		}
		c.Status(http.StatusOK)
	} else {
		c.JSON(http.StatusOK, res)
	}

	return nil
}

func (u user) GetMe(ctx context.Context, c *gin.Context) error {
	res, err := u.usecase.GetByID(ctx, util.UID(ctx))
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, res)
	return nil
}
