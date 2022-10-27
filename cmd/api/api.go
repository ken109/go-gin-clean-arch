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
	httpController "go-gin-ddd/adapter/controller/http"
	"go-gin-ddd/adapter/gateway/mail"
	mysqlRepository "go-gin-ddd/adapter/gateway/mysql"
	"go-gin-ddd/adapter/presenter"
	"go-gin-ddd/driver"
	"go-gin-ddd/packages/log"
	"go-gin-ddd/usecase"

	"go-gin-ddd/packages/http/middleware"
	"go-gin-ddd/packages/http/router"

	"go-gin-ddd/config"
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

	r := router.New(engine, driver.GetRDB)

	// dependencies injection
	// ----- gateway -----
	mailAdapter := mail.New()

	// mysql
	userRepository := mysqlRepository.NewUser()

	// ----- usecase -----
	userInputFactory := usecase.NewUserInputFactory(userRepository, mailAdapter)
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
		logger.Fatal(fmt.Sprintf("Server forced to shutdown: %+v", err))
	}

	logger.Info("Server exiting")
}
