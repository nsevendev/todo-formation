package usercontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetProfilCurrentUser godoc
// @Summary Récupérer le profil de l'utilisateur connecté
// @Description Récupére le profil de l'utilisateur connecté via le token utilisé dans le header
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ginresponse.JsonFormatterSwag "Profil utilisateur récupéré avec succès"
// @Failure 401 {object} ginresponse.JsonFormatterSwag "Token invalide"
// @Failure 500 {object} ginresponse.JsonFormatterSwag "Erreur interne"
// @Router /user/profil [get]
func (u *userController) GetProfilCurrentUser(c *gin.Context) {
	idUser, exists := c.Get("id_user")
	if !exists {
		logger.Ef("Erreur d'authentification : ID utilisateur non trouvé dans le contexte")
		ginresponse.Unauthorized(c, "Erreur d'authentification", "Vous n'avez pas les droits pour effectuer cette action")
		return
	}

	user, err := u.userService.GetProfilCurrentUser(c, idUser.(primitive.ObjectID))
	if err != nil {
		logger.Ef("Erreur à la récupération du profil de l'utilisateur: %v", err)
		ginresponse.InternalServerError(c, "Erreur à la récupération du profil de l'utilisateur", err.Error())
		return
	}

	ginresponse.Success(c, "Profil utilisateur récupéré avec succès", user)
}
