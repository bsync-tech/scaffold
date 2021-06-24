package server

import "github.com/gin-gonic/gin"

func LogMessage(c *gin.Context) {
	subsys := c.GetString("subsys")
	module := c.GetString("module")
	msg := c.GetString("msg")

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
