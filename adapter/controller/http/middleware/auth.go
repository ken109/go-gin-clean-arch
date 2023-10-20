package middleware

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/ken109/gin-jwt"
)

func Auth(must bool, realm string, session bool) gin.HandlerFunc {
	verifyFunc := jwt.TryVerify
	if must {
		verifyFunc = jwt.MustVerify
	}

	if session {
		return func(c *gin.Context) {
			if c.GetHeader("Authorization") == "" {
				session := sessions.DefaultMany(c, realm)
				token := session.Get("token")
				if token, ok := token.(string); ok {
					c.Request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
				}
			}

			verifyFunc(realm)(c)
		}
	} else {
		return verifyFunc(realm)
	}
}
