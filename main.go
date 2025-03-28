package main

import (
	"os"
	"strings"
	"todof/app/router"
	_ "todof/init"
	"todof/internal/logger"

	"github.com/gin-gonic/gin"
)

// permet de recuperer le nom de l'host traefik qui est dans le .env
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
 