package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

func RouteNotfound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() == http.StatusNotFound && c.Request.Method != "OPTIONS" {
			logger.Wf("Route inconnue : %s %s", c.Request.Method, c.Request.URL.Path)
			ginresponse.NotFound(c, "La route demandée n'existe pas.", nil)
			c.Abort()
		}
	}
}