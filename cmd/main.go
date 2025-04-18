package main

import (
	"os"
	"strings"
	_ "todof/internal/init"
	"todof/mod/middleware"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/logger/v2/logger"
)

func main() {
	s := gin.Default()
	defer logger.Close()

	s.Use(middleware.RouteNotfound())

	run(s)
}

// run and log the server
func run(s *gin.Engine) {
	port := os.Getenv("PORT")
	hostTraefik := extractStringInBacktick(os.Getenv("HOST_TRAEFIK"))
	host := "0.0.0.0"

	logger.S("Server is running on in container docker : " + host + ":" + port)
	logger.Sf("Server is running on navigator on : https://%v", hostTraefik)

	logger.I("Démarrage du serveur ...")		
	if err := s.Run(host + ":" + port);err != nil {
		logger.Ff("Erreur lors du démarrage du serveur : %v", err)
	}
}

// utiliser pour recuperer une string entre des backtick
func extractStringInBacktick(s string) string {
	start := strings.Index(s, "`")
	end := strings.LastIndex(s, "`")

	if start == -1 || end == -1 || start == end {
		return ""
	}

	return s[start+1 : end]
}