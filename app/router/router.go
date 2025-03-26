package router

import (
	"github.com/gin-gonic/gin"
)

func Router(serv *gin.Engine) {
	serv.GET("/", "controllers")
}