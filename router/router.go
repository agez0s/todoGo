package router

import "github.com/gin-gonic/gin"

func Initialize() {
	r := gin.New()
	r.Use(gin.Recovery())
	initializeRoutes(r)
	r.Run()
}
