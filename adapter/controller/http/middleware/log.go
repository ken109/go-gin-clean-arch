package middleware

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-gin-clean-arch/packages/cerrors"
)

func Log(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(logger *zap.Logger) {
			_ = logger.Sync()
		}(logger)

		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		var (
			logFunc = logger.Info
			fields  = []zap.Field{
				zap.String("request-id", c.Writer.Header().Get("X-Request-Id")),
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(timeFormat)),
				zap.Duration("latency", latency),
			}
		)

		status := c.Writer.Status()
		if status >= 400 {
			var realErrors []error
			for _, err := range c.Errors {
				realErrors = append(realErrors, err.Err)
			}

			errorsJSON, _ := json.Marshal(realErrors)

			fields = append(fields, zap.Any("errors", json.RawMessage(errorsJSON)))
		}
		if status >= 500 {
			logFunc = logger.Error
		}

		logFunc(path, fields...)
	}
}

func RecoveryWithLog(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(
							strings.ToLower(se.Error()),
							"broken pipe",
						) || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(
						c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error))
					c.Abort()
					return
				}

				unexpectedErr := cerrors.NewUnexpected(fmt.Errorf("%+v", err), cerrors.WithUnexpectedPanic{})

				_ = c.Error(unexpectedErr)

				unexpectedErr.Response(c)
			}
		}()
		c.Next()
	}
}
