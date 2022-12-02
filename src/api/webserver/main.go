package webserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"medusa/src/api/auth"
	"medusa/src/api/config"
	"medusa/src/api/routing"
	"medusa/src/common/env"
	"medusa/src/common/logger"
	"medusa/src/common/middleware"
	"medusa/src/common/middleware/rateLimit"
	"medusa/src/common/middleware/requestLog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Boot(envFiles ...string) {
	env.Load(envFiles...)
	webConfig := &config.WebServerConfig{}
	config.Load("web", webConfig)

	logger.SetupErrorLog()
	router := gin.New()

	middlewares := []middleware.RouterMiddleware{
		&rateLimit.RateLimit{},
		&requestLog.RequestLog{},
	}

	for _, middle := range middlewares {
		middle.Register(router)
	}

	router.Use(gin.Recovery())

	authService := auth.NewAuth()

	router.Use(func(context *gin.Context) {
		authenticationError := authService.Authenticate(context.Request)
		if authenticationError != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, authenticationError)
		}
		context.Next()
	})
	router.Use(func(context *gin.Context) {
		authorizationError := authService.Authorize(context.Request)
		if authorizationError != nil {
			context.AbortWithStatusJSON(http.StatusForbidden, authorizationError)
		}
		context.Next()
	})

	routing.MapUrls(router)

	log.Info("Starting server on port ", os.Getenv("SERVER_PORT"))

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler: router,
	}

	go gracefulShutdown(
		server,
		quit,
		done,
	)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Info("Graceful shutdown completed")
}

func gracefulShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v", err)
	}

	close(done)
}
