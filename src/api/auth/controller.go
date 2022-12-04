package auth

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	apiUser "medusa/src/api/user"
	commonUser "medusa/src/common/user"
	"medusa/src/common/utils"
	"net/http"
	"strings"
	"time"
)

type AuthController struct {
	userRepository *commonUser.UserRepository
	jwt            *utils.Jwt
}

func NewAuthController(userRepository *commonUser.UserRepository, jwt *utils.Jwt) *AuthController {
	return &AuthController{
		userRepository: userRepository,
		jwt:            jwt,
	}
}

func (ac *AuthController) Register(routerGroup *gin.RouterGroup) {
	router := routerGroup.Group("/auth")

	router.POST("/register", ac.singUp)
	router.POST("/login", ac.signIn)
}

func (ac *AuthController) singUp(context *gin.Context) {
	var user *SignUpInput

	if err := context.ShouldBindJSON(&user); err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	newUser := &commonUser.UserDbModel{
		Name:      user.Name,
		Email:     strings.ToLower(user.Email),
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      "default",
		Verified:  false,
	}
	err = ac.userRepository.Insert(newUser)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	accessToken, err := ac.jwt.CreateJwt(newUser.ID)

	context.JSON(
		http.StatusCreated,
		gin.H{"status": "success", "user": apiUser.MapToUserResponse(newUser), "accessToken": accessToken},
	)
}

func (ac *AuthController) signIn(context *gin.Context) {
	var credentials *SignInInput

	if err := context.ShouldBindJSON(&credentials); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userRepository.FindOneBy(bson.M{"email": strings.ToLower(credentials.Email)})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err = utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
		return
	}

	accessToken, err := ac.jwt.CreateJwt(user.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}
