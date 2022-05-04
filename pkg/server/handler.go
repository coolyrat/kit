package server

import "github.com/gin-gonic/gin"

func health(c *gin.Context) {
	c.JSON(200, gin.H{
		"server": "up",
	})
}

func defaultHandler(e *gin.Engine) *gin.Engine {
	e.GET("/health", health)
	return e
}
