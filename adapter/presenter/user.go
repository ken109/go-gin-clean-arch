package presenter

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-gin-ddd/config"
	"go-gin-ddd/domain"
	"go-gin-ddd/resource/response"
	"go-gin-ddd/usecase"
)

type user struct {
	c *gin.Context
}

type UserOutputFactory func(c *gin.Context) usecase.UserOutputPort

func NewUserOutputFactory() UserOutputFactory {
	return func(c *gin.Context) usecase.UserOutputPort {
		return &user{c: c}
	}
}

func (u *user) Create(id uint) error {
	u.c.JSON(http.StatusCreated, id)
	return nil
}

func (u *user) ResetPasswordRequest(res *response.UserResetPasswordRequest) error {
	u.c.JSON(http.StatusOK, res)
	return nil
}

func (u *user) ResetPassword() error {
	u.c.Status(http.StatusOK)
	return nil
}

func (u *user) Login(isSession bool, res *response.UserLogin) error {
	if res == nil {
		u.c.Status(http.StatusUnauthorized)
		return nil
	}

	if isSession {
		session := sessions.DefaultMany(u.c, config.UserRealm)
		session.Set("token", res.Token)
		session.Set("refresh_token", res.RefreshToken)
		if err := session.Save(); err != nil {
			return err
		}
		u.c.Status(http.StatusOK)
	} else {
		u.c.JSON(http.StatusOK, res)
	}

	return nil
}

func (u *user) RefreshToken(isSession bool, res *response.UserLogin) error {
	if res == nil {
		u.c.Status(http.StatusUnauthorized)
		return nil
	}

	if isSession {
		session := sessions.DefaultMany(u.c, config.UserRealm)
		session.Set("token", res.Token)
		session.Set("refresh_token", res.RefreshToken)
		if err := session.Save(); err != nil {
			return err
		}
		u.c.Status(http.StatusOK)
	} else {
		u.c.JSON(http.StatusOK, res)
	}

	return nil
}

func (u *user) GetByID(res *domain.User) error {
	u.c.JSONP(http.StatusOK, res)
	return nil
}
