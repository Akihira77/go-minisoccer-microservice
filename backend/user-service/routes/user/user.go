package user

import (
	"user-service/controllers"
	"user-service/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IUserRoute interface {
	Run()
}

func NewUserRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IUserRoute {
	return &UserRoute{
		controller: controller,
		group:      group,
	}
}

func (ur *UserRoute) Run() {
	group := ur.group.Group("/auth")
	group.GET("/user", middlewares.Authenticate(), ur.controller.GetUserController().GetUserLogin)
	group.GET("/:uuid", middlewares.Authenticate(), ur.controller.GetUserController().GetUserByUUID)
	group.POST("/login", ur.controller.GetUserController().Login)
	group.POST("/register", ur.controller.GetUserController().Register)
	group.PUT("/:uuid", middlewares.Authenticate(), ur.controller.GetUserController().Update)
}
