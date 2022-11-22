package rateLimit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

type RateLimit struct {
}

func (rt *RateLimit) Register(router *gin.Engine) {
	concurrentRequests, err := strconv.ParseInt(os.Getenv("CONCURRENT_REQUESTS"), 10, 64)

	if err != nil {
		log.Warning(fmt.Sprintf("Could not read environment variable CONCURRENT_REQUESTS (%s). Defaulting to 100.", err.Error()))
		concurrentRequests = 100
	}

	if concurrentRequests < 2 {
		concurrentRequests = 2
	}

	/*
		A buffered channel of size 1 will accept 2 messages at a time.
		So in order to limit the server to X number of concurrent requests, we need to create a buffered channel of size X - 1.
	*/
	buffer := make(chan bool, concurrentRequests-1)

	router.Use(func(context *gin.Context) {
		select {
		case buffer <- true:
			context.Next()
			<-buffer
		default:
			log.Warning("Too many concurrent requests. Aborting.")
			context.Abort()
			context.JSON(http.StatusTooManyRequests, "test")
		}
	})
}
