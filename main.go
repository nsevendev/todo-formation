package main

import (
	"todo_formation/app/router"
	_ "todo_formation/init"
	"todo_formation/internal/logger"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func extractBacktickContent(s string) string {
	start := strings.Index(s, "`")
	end := strings.LastIndex(s, "`")

	if start == -1 || end == -1 || start == end {
		return ""
	}

	return s[start+1 : end]
}

func main() {	
	serv := gin.Default()

	router.Router(serv)

	port := os.Getenv("PORT")
	host := "0.0.0.0"
	hostTraefik := extractBacktickContent(os.Getenv("HOST_TRAEFIK"))

	logger.Success("Server is running on " + host + ":" + port)
	logger.Successf("Server is running on https://%v", hostTraefik)

	serv.Run(host + ":" + port)
}
 