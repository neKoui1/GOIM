package router

import (
	"GOIM/middlewares"
	"GOIM/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/login", service.Login)
			user.POST("/register", service.Register)
			user.POST("/sendcode", service.SendCode)
			auth := user.Group("/", middlewares.AuthCheck())
			{
				auth.GET("/info", service.GetUserInfo)
				auth.GET("/websocket/message", service.WebSocketMessage)
			}
		}
	}


	return r
}
