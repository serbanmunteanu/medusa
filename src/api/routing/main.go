package routing

import (
	"github.com/gin-gonic/gin"
	"medusa/src/api/auth"
	"medusa/src/api/config"
	"medusa/src/api/swagger"
	apiUser "medusa/src/api/user"
	"medusa/src/common/mongodb"
	commonUser "medusa/src/common/user"
	"medusa/src/common/utils"
)

type RouteRegister interface {
	Register(routerGroup *gin.RouterGroup)
}

type RouteGroup struct {
	groupPrefix    string
	routeRegisters []RouteRegister
	middlewares    []gin.HandlerFunc
}

func Initialize(router *gin.Engine, config *config.WebServerConfig) {
	mongoClient, err := mongodb.NewMongoClient(config.MongoConfig)
	if err != nil {
		panic(err)
	}
	userRepository := commonUser.NewUserRepository(mongoClient, config.MongoConfig.Collections.UserCollection)
	jwt := utils.NewJwt(config.JwtConfig)
	authHandler := auth.NewAuthHandler(jwt, userRepository)
	routeGroups := []RouteGroup{
		{
			groupPrefix: "",
			routeRegisters: []RouteRegister{
				swagger.NewSwaggerController(),
				auth.NewAuthController(userRepository, jwt),
			},
			middlewares: []gin.HandlerFunc{},
		},
		{
			groupPrefix: "/api",
			routeRegisters: []RouteRegister{
				apiUser.NewUserController(userRepository),
			},
			middlewares: []gin.HandlerFunc{
				authHandler.GetAuthentication(),
				authHandler.GetAuthorization(),
			},
		},
	}

	for _, routeGroup := range routeGroups {
		group := router.Group(routeGroup.groupPrefix)
		for _, middle := range routeGroup.middlewares {
			group.Use(middle)
		}
		for _, route := range routeGroup.routeRegisters {
			route.Register(group)
		}
	}
}
