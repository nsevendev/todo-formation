package main

import (
	"github.com/nsevenpack/env/env"
	"todof/docs"
	_ "todof/docs"
	_ "todof/internal/init"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/logger/v2/logger"
	"strings"
	"todof/app/router"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// @title API todo-formation
// @version 1.0
// @description API pour créer des todo avec utilisateurs. Pour tester les routes protégées, cliquez sur le bouton Authorize et saisissez : Bearer {votre token} (remplacez {votre token} par un token valide obtenu via la route /user/login).
// @schemes https
// @securityDefinitions.apikey BearerAuth
// @in headers
// @name Authorization
func main() {
	s := gin.Default()
	defer logger.Close()

	s.Use(CORSMiddleware())
	router.Router(s)
	logRoutes(s)

	run(s)
}

func logRoutes(s *gin.Engine) {
	for _, route := range s.Routes() {
		logger.If("%-6s -> %s\n", route.Method, route.Path)
	}
}

// run and log the server
func run(s *gin.Engine) {
	port := env.Get("PORT")
	hostTraefik := extractStringInBacktick(env.Get("HOST_TRAEFIK"))
	host := "0.0.0.0"

	docs.SwaggerInfo.Host = hostTraefik

	logger.S("Server is running on in container docker : " + host + ":" + port)
	logger.Sf("Server is running on navigator on : https://%v", hostTraefik)

	logger.I("Démarrage du serveur ...")
	if err := s.Run(host + ":" + port); err != nil {
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
