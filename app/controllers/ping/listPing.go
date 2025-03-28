package ping

import (
	"time"

	"github.com/gin-gonic/gin"
)

func ListPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"message": "pong",
		"time": time.Now().Format(time.RFC3339),
	})
}