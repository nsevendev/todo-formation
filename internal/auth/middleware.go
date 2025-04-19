package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type authMiddleware struct {
	authService UserServiceInterface
}

type AuthMiddlewareInterface interface {
	RequireAuth() gin.HandlerFunc
	RequireRole(roles ...string) gin.HandlerFunc
}

func NewAuthMiddleware(authService UserServiceInterface) AuthMiddlewareInterface {
	return &authMiddleware{
		authService: authService,
	}
}

func (m *authMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraire le token du header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			ginresponse.Unauthorized(c, "Authorization header est requis", ginresponse.ErrorModel{
				Message: "Authorization header est requis",
				Type:    "Authorization",
				Detail:  "Authorization header n'est pas présent dans la requête",
			})
			c.Abort()
			return
		}

		// Format attendu: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ginresponse.Unauthorized(c, "Authorization header doit être de format 'Bearer <token>'", ginresponse.ErrorModel{
				Message: "Authorization header doit être de format 'Bearer <token>'",
				Type:    "Authorization",
				Detail:  "Authorization header n'est pas au bon format",
			})
			c.Abort()
			return
		}

		// Valider le token
		claims, err := m.authService.ValidateToken(parts[1])
		if err != nil {
			ginresponse.Unauthorized(c, "Invalide token", ginresponse.ErrorModel{
				Message: "Invalide token",
				Type:    "Token",
				Detail:  fmt.Sprintf("Token invalide ou expiré: %v", err),
			})
			c.Abort()
			return
		}

		// Convertir l'ID de l'utilisateur en ObjectID
		idUser, err := primitive.ObjectIDFromHex(claims.IdUser)
		if err != nil {
			ginresponse.Unauthorized(c, "Contenu du token invalide", ginresponse.ErrorModel{
				Message: "Contenu du token invalide",
				Type:    "Token",
				Detail:  fmt.Sprintf("Contenu du token invalide: %v", err),
			})
			c.Abort()
			return
		}

		// Ajouter les claims au contexte de la requête
		c.Set("id_user", idUser)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func (m *authMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		authorized := false
		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role format"})
			c.Abort()
			return
		}

		for _, allowedRole := range roles {
			if roleStr == allowedRole {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}