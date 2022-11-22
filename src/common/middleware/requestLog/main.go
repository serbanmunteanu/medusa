package requestLog

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"medusa/src/common/logger"
)

type RequestLog struct {
}

func (rl *RequestLog) Register(router *gin.Engine) {
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
}
