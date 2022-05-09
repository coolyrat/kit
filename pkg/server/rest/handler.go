package rest

import "github.com/gin-gonic/gin"

func health(c *gin.Context) {
	c.JSON(200, gin.H{
		"server": "up",
	})
}

func DefaultHandler(e *gin.Engine) *gin.Engine {
	e.GET("/health", health)
	return e
}
