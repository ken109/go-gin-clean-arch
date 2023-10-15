package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/ken109/gin-jwt"
	httpController "go-gin-clean-arch/adapter/controller/http"
	"go-gin-clean-arch/adapter/gateway/mail"
	mysqlRepository "go-gin-clean-arch/adapter/gateway/mysql"
	"go-gin-clean-arch/adapter/presenter"
	"go-gin-clean-arch/driver"
	"go-gin-clean-arch/packages/log"
	"go-gin-clean-arch/usecase"

	"go-gin-clean-arch/packages/http/middleware"
	"go-gin-clean-arch/packages/http/router"

	"go-gin-clean-arch/config"
)

func main() {
	logger := log.Logger()
	defer logger.Sync()

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

	// cors
	engine.Use(middleware.Cors(nil))

	// health check
	engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// middlewares
	engine.Use(middleware.Log(log.ZapLogger(), time.RFC3339, false))
	engine.Use(middleware.RecoveryWithLog(log.ZapLogger(), true))

	// cookie
	engine.Use(middleware.Session([]string{config.UserRealm}, config.Env.App.Secret, nil))

	r := router.New(engine, driver.GetRDB)

	// dependencies injection
	// ----- gateway -----
	mailSender := mail.NewSender()

	// mysql
	userRepository := mysqlRepository.NewUser()

	// ----- usecase -----
	userInputFactory := usecase.NewUserInputFactory(userRepository, mailSender)
	userOutputFactory := presenter.NewUserOutputFactory()

	// ----- controller -----
	httpController.NewUser(r, userInputFactory, userOutputFactory)

	logger.Info("Succeeded in dependencies injection.")

	// serve
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Env.Port),
		Handler: engine,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
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
		logger.Fatal(fmt.Sprintf("Server forced to shutdown: %+v", err))
	}

	logger.Info("Server exiting")
}
