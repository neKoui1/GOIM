package middlewares

import (
	"GOIM/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		userClaims, err := helper.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "认证失败",
			})
			return
		}
		c.Set("user_claims", userClaims)
		c.Next()
	}
}
