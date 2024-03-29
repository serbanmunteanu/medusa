package swagger

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
)

type SwaggerController struct {
}

func NewSwaggerController() *SwaggerController {
	return &SwaggerController{}
}

func (sc *SwaggerController) Register(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/swagger", func(context *gin.Context) {
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
