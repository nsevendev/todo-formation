package middleware

import (
	"net/http"
	"todof/mod/apiresponse"
	"todof/mod/logger"

	"github.com/gin-gonic/gin"
)

func RouteNotfound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() == http.StatusNotFound && c.Request.Method != "OPTIONS" {
			logger.Warnf("Route inconnue : %s %s", c.Request.Method, c.Request.URL.Path)

			apiresponse.Error(c, http.StatusNotFound, "La route demandée n'existe pas.", "Not Found")

			c.Abort()
		}
	}
}