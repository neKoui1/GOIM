package router

import (
	"GOIM/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/login", service.Login)

	return r
}
