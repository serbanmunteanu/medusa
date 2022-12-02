package routing

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
)

func MapUrls(router *gin.Engine) {
	router.GET("/swagger", func(context *gin.Context) {
		filePath, err := getSwaggerFilepath()
		if err != nil {
			log.Error(err)
			context.AbortWithStatus(http.StatusInternalServerError)
		}
		swaggerDocs, err := os.ReadFile(filePath)
		if err != nil {
			log.Error(err)
			context.AbortWithStatus(http.StatusInternalServerError)
		}

		context.Data(
			http.StatusOK,
			gin.MIMEJSON,
			swaggerDocs,
		)
	})
}

func getSwaggerFilepath() (string, error) {
	if os.Getenv("APP_ENV") == "prod" {
		return filepath.Abs("./docs/swagger.json")
	}

	return filepath.Abs("api/docs/swagger.json")
}
