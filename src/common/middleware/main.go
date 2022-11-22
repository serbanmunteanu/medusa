package middleware

import "github.com/gin-gonic/gin"

type RouterMiddleware interface {
	Register(router *gin.Engine)
}
