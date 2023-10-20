package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type SessionOption struct {
	MaxAge time.Duration
}

func Session(name []string, secret string, option *SessionOption) gin.HandlerFunc {
	var (
		secure   bool
		sameSite http.SameSite
	)

	switch gin.Mode() {
	case gin.ReleaseMode:
		secure = true
		sameSite = http.SameSiteStrictMode
	case gin.TestMode:
		secure = true
		sameSite = http.SameSiteNoneMode
	case gin.DebugMode:
		secure = false
		sameSite = http.SameSiteLaxMode
	}

	maxAge := time.Hour * 24 * 365

	if option != nil {
		maxAge = option.MaxAge
	}

	store := cookie.NewStore([]byte(secret))
	store.Options(
		sessions.Options{
			Path:     "/",
			MaxAge:   int(maxAge),
			Secure:   secure,
			HttpOnly: true,
			SameSite: sameSite,
		},
	)
	return sessions.SessionsMany(name, store)
}
