package routing

import (
	"github.com/gin-gonic/gin"
	"medusa/src/api/auth"
	"medusa/src/api/config"
	"medusa/src/api/swagger"
	"medusa/src/common/mongodb"
	"medusa/src/common/users"
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
	userRepository := users.NewUserRepository(mongoClient, config.MongoConfig.Collections.UserCollection)
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
				auth.NewAuthController(userRepository),
			},
		},
	}

	for _, routeGroup := range routeGroups {
		for _, route := range routeGroup.routeRegisters {
			route.Register(router.Group(routeGroup.groupPrefix))
		}
	}
}
