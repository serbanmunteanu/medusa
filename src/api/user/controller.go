package user

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"medusa/src/common/redis"
	commonUser "medusa/src/common/user"
	"net/http"
)

type UserController struct {
	userRepository *commonUser.UserRepository
	redisService   *redis.Service
}

func NewUserController(userRepository *commonUser.UserRepository, redisService *redis.Service) *UserController {
	return &UserController{
		userRepository: userRepository,
		redisService:   redisService,
	}
}

func (us *UserController) Register(routerGroup *gin.RouterGroup) {
	router := routerGroup.Group("/users")

	router.GET("/", us.GetAll)
}

func (us *UserController) GetAll(context *gin.Context) {
	var users []commonUser.UserDbModel
	found, err := us.redisService.HasKey("users:getAll")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	log.Info(found)
	users, err = us.userRepository.Find(bson.M{})
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	err = us.redisService.Set("users:getAll", users, 500)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "success", "users": MapToUsersResponse(users)})
}
