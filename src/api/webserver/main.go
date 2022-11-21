package webserver

import (
    "context"
    "fmt"
    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
    "medusa/src/api/routing"
    "medusa/src/common/env"
    "medusa/src/common/logger"
    "net/http"
    "os"
    "os/signal"
    "time"
)

func Boot(envFiles ...string) {
    env.Load(envFiles...)

    logger.SetupErrorLog()
    router := gin.New()

    router.Use(gin.Recovery())
    accessLog := logger.SetupAccessLog()

    router.Use(func(context *gin.Context) {
        context.Next()

        accessLog.WithFields(log.Fields{
            "ip":     context.Request.RemoteAddr,
            "url":    context.Request.URL.Path,
            "method": context.Request.Method,
            "status": context.Writer.Status(),
        }).Info("")
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
