package routes

import (
	"user-service/controllers"
	"user-service/routes/user"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IRouteRegistry interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup) IRouteRegistry {
	return &Registry{
		controller: controller,
		group:      group,
	}
}

func (r *Registry) userRoute() user.IUserRoute {
	return user.NewUserRoute(r.controller, r.group)
}

func (r *Registry) Serve() {
	r.userRoute().Run()
}
