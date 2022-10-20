package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-gin-ddd/config"
	"go-gin-ddd/packages/http/middleware"
	"go-gin-ddd/packages/http/router"

	"go-gin-ddd/packages/context"

	"go-gin-ddd/resource/request"
	"go-gin-ddd/usecase"
)

type user struct {
	userUseCase usecase.IUser
}

func NewUser(r *router.Router, uuc usecase.IUser) {
	handler := user{userUseCase: uuc}

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

	id, err := u.userUseCase.Create(ctx, &req)
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

	res, err := u.userUseCase.ResetPasswordRequest(ctx, &req)
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

	err := u.userUseCase.ResetPassword(ctx, &req)
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

	res, err := u.userUseCase.Login(ctx, &req)
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
		if err = session.Save(); err != nil {
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

	res, err := u.userUseCase.RefreshToken(req.RefreshToken)
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
		if err = session.Save(); err != nil {
			return err
		}
		c.Status(http.StatusOK)
	} else {
		c.JSON(http.StatusOK, res)
	}

	return nil
}

func (u user) GetMe(ctx context.Context, c *gin.Context) error {
	user, err := u.userUseCase.GetByID(ctx, ctx.UID())
	if err != nil {
		return err
	}

	c.JSONP(http.StatusOK, user)
	return nil
}
