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

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	jwt "github.com/ken109/gin-jwt"

	httpController "go-gin-clean-arch/adapter/controller/http"
	"go-gin-clean-arch/adapter/controller/http/middleware"
	"go-gin-clean-arch/adapter/controller/http/router"
	"go-gin-clean-arch/adapter/gateway/mail"
	mysqlRepository "go-gin-clean-arch/adapter/gateway/mysql"
	"go-gin-clean-arch/config"
	"go-gin-clean-arch/driver"
	"go-gin-clean-arch/packages/log"
	"go-gin-clean-arch/usecase"
)

func main() {
	logger := log.ZapLogger()
	defer logger.Sync()

	err := jwt.SetUp(
		jwt.Option{
			Realm:            config.UserRealm,
			SigningAlgorithm: jwt.HS256,
			SecretKey:        []byte(config.Env.App.Secret),
		},
	)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to set up JWT: %+v", err))
	}
	logger.Info("Succeeded in setting up JWT.")

	engine := gin.New()

	// cors
	engine.Use(middleware.Cors(nil))

	// health check
	engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// middlewares
	engine.Use(requestid.New())
	engine.Use(middleware.Log(log.ZapLogger(), time.RFC3339, false))
	engine.Use(middleware.RecoveryWithLog(log.ZapLogger()))

	// cookie
	engine.Use(middleware.Session([]string{config.UserRealm}, config.Env.App.Secret, nil))

	db, err := driver.NewRDB()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to DB: %+v", err))
	}

	// dependencies injection
	// ----- gateway -----
	mailSender := mail.NewMailSender()

	// mysql
	transactionRepository := mysqlRepository.NewTransaction()
	userRepository := mysqlRepository.NewUser()

	// ----- usecase -----
	userUsecase := usecase.NewUser(transactionRepository, userRepository, mailSender)

	// ----- controller -----
	r := router.New(engine, db)

	httpController.NewUser(r, userUsecase)

	// serve
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.Env.Port),
		Handler:           engine,
		ReadHeaderTimeout: time.Second * 20,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(fmt.Sprintf("Failed to listen and serve: %+v", err))
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
