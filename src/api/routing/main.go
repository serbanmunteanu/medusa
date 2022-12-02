package routing

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type EmailTaskPayload struct {
	// ID for the email recipient.
	UserID int
}

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
	router.GET("/test-asynq", func(context *gin.Context) {
		client := asynq.NewClient(asynq.RedisClientOpt{
			Addr:     "localhost:6379",
			Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		})

		// Create a task with typename and payload.
		payload, err := json.Marshal(EmailTaskPayload{UserID: 42})
		if err != nil {
			log.Fatal(err)
		}
		t1 := asynq.NewTask("email:welcome", payload)

		t2 := asynq.NewTask("email:reminder", payload)

		// Process the task immediately.
		info, err := client.Enqueue(t1)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(" [*] Successfully enqueued task: %+v", info)

		// Process the task 24 hours later.
		info, err = client.Enqueue(t2, asynq.ProcessIn(24*time.Hour))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(" [*] Successfully enqueued task: %+v", info)
	})
}

func getSwaggerFilepath() (string, error) {
	if os.Getenv("APP_ENV") == "prod" {
		return filepath.Abs("./docs/swagger.json")
	}

	return filepath.Abs("api/docs/swagger.json")
}
