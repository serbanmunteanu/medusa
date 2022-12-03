package user

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	commonUser "medusa/src/common/user"
	"net/http"
)

type UserController struct {
	userRepository *commonUser.UserRepository
}

func NewUserController(userRepository *commonUser.UserRepository) *UserController {
	return &UserController{
		userRepository: userRepository,
	}
}

func (us *UserController) Register(routerGroup *gin.RouterGroup) {
	router := routerGroup.Group("/users")

	router.GET("/", us.GetAll)
}

func (us *UserController) GetAll(context *gin.Context) {
	users, err := us.userRepository.Find(bson.M{})

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "success", "users": users})
}
