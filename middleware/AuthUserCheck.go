package middleware

import (
	"gin_gorm_o/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthAdminCheck is auth admin middleware
func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("authorization")
		userClaim, err := helper.ParseToken(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Get token error: " + err.Error(),
			})
			return
		}

		if userClaim == nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "unauthorized user",
			})
			return
		}
		c.Set("user", userClaim)
		c.Next()
	}
}
