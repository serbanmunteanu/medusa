package auth

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	apiUser "medusa/src/api/users"
	commonUser "medusa/src/common/users"
	"net/http"
	"time"
)

type AuthController struct {
	userRepository *commonUser.UserRepository
}

func NewAuthController(userRepository *commonUser.UserRepository) *AuthController {
	return &AuthController{
		userRepository: userRepository,
	}
}

func (ac *AuthController) Register(routerGroup *gin.RouterGroup) {
	router := routerGroup.Group("/auth")

	router.POST("/register", ac.singUpUser)
}

func (ac *AuthController) singUpUser(context *gin.Context) {
	var user *apiUser.SignUpInput

	if err := context.ShouldBindJSON(&user); err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "messages": err.Error()})
		return
	}
	newUser := &commonUser.UserDbModel{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      "default",
		Verified:  false,
	}
	err := ac.userRepository.Insert(newUser)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "messages": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": apiUser.MapToUserResponse(newUser)}})
}
