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
}

func MapUrls(router *gin.Engine, config *config.WebServerConfig) {
	mongoClient, err := mongodb.NewMongoClient(config.MongoConfig)
	if err != nil {
		panic(err)
	}
	userRepository := commonUser.NewUserRepository(mongoClient, config.MongoConfig.Collections.UserCollection)
	jwt := utils.NewJwt(config.JwtConfig)
	routeGroups := []RouteGroup{
		{
			groupPrefix: "",
			routeRegisters: []RouteRegister{
				swagger.NewSwaggerController(),
			},
		},
		{
			groupPrefix: "/api",
			routeRegisters: []RouteRegister{
				auth.NewAuthController(userRepository, jwt),
				apiUser.NewUserController(userRepository),
			},
		},
	}

	for _, routeGroup := range routeGroups {
		for _, route := range routeGroup.routeRegisters {
			route.Register(router.Group(routeGroup.groupPrefix))
		}
	}
}
