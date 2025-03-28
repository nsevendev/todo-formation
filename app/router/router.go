package router

import (
	"todof/app/controllers/ping"

	"github.com/gin-gonic/gin"
)

func Router(serv *gin.Engine) {
	serv.GET("/", ping.ListPing)
}