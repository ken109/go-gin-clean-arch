package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/ken109/gin-jwt"
	mysqlRepository "go-gin-ddd/infrastructure/mysql"

	"go-gin-ddd/packages/http/middleware"
	"go-gin-ddd/packages/http/router"

	"go-gin-ddd/config"
	httpController "go-gin-ddd/controller/http"
	"go-gin-ddd/driver/rdb"
	"go-gin-ddd/infrastructure/email"
	"go-gin-ddd/infrastructure/log"
	"go-gin-ddd/usecase"
)

func Execute() {
	logger := log.Logger()

	err := jwt.SetUp(
		jwt.Option{
			Realm:            config.UserRealm,
			SigningAlgorithm: jwt.HS256,
			SecretKey:        []byte(config.Env.App.Secret),
		},
	)
	if err != nil {
		panic(err)
	}
	logger.Info("Succeeded in setting up JWT.")

	engine := gin.New()

	engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	engine.Use(middleware.Log(log.ZapLogger(), time.RFC3339, false))
	engine.Use(middleware.RecoveryWithLog(log.ZapLogger(), true))

	// cors
	engine.Use(middleware.Cors(nil))

	// cookie
	engine.Use(middleware.Session([]string{config.UserRealm}, config.Env.App.Secret, nil))

	r := router.New(engine, rdb.Get)

	// dependencies injection
	// ----- infrastructure -----
	emailDriver := email.New()

	// mysql
	userRepository := mysqlRepository.NewUser()

	// ----- usecase -----
	userUseCase := usecase.NewUser(emailDriver, userRepository)

	// ----- controller -----
	httpController.NewUser(r, userUseCase)

	logger.Info("Succeeded in dependencies injection.")

	// serve
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Env.Port),
		Handler: engine,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	logger.Info("Succeeded in listen and serve.")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %+v", err)
	}

	logger.Info("Server exiting")
}
