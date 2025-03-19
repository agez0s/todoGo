package router

import (
	"github.com/agez0s/todoGo/docs"
	"github.com/agez0s/todoGo/handler"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initializeRoutes(router *gin.Engine) {
	handler.Init()
	basePath := "/api"
	docs.SwaggerInfo.BasePath = basePath
	v1 := router.Group(basePath + "/v1/")

	//authorization routes
	v1.POST("/auth/newUser", handler.CreateUserHandler)
	v1.POST("/auth/login", handler.LoginUserHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
