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
				// 下载包括所有用户的Excel
				user.GET("/download", service.ExportUserExcel)
				auth.GET("/info", service.GetUserInfo)
				// 查询指定用户的个人信息
				auth.GET("/param/:account", service.UserParam)
				// 添加好友
				auth.POST("/add", service.UserAdd)
				// 删除好友
				auth.POST("/delete", service.UserDelete)

				auth.GET("/websocket/message", service.WebSocketMessage)
				auth.GET("/chat/list", service.ChatList)
			}
		}
	}

	return r
}
