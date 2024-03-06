package service

import "github.com/gin-gonic/gin"

func Hello(c *gin.Context) {
	c.JSON(200, "hello,胡海峰")
}
