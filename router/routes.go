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

	//todo routes
	v1.POST("/todo/create", handler.AuthMiddleware(), handler.CreateTodoHandler)
	v1.POST("/todo/update", handler.AuthMiddleware(), handler.UpdateTodoHandler)
	v1.GET("/todo/list", handler.AuthMiddleware(), handler.ListTodosHandler)
	v1.DELETE("/todo/delete/:id", handler.AuthMiddleware(), handler.DeleteTodoHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
