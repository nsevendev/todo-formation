package usercontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

// DeleteOneByUser godoc
// @Summary Supprime l'utilisateur connecté
// @Description Supprime l'utilisateur connecté via le token utilisé dans le header
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 204 "1"
// @Failure 401 {object} ginresponse.JsonFormatterSwag "Token invalide"
// @Failure 500 {object} ginresponse.JsonFormatterSwag "Erreur interne"
// @Router /user [delete]
func (u *userController) DeleteOneByUser(c *gin.Context) {
	idUser := u.userService.GetIdUserInContext(c)

	if err := u.userService.DeleteOneByUser(c, idUser); err != nil {
		logger.Ef("Erreur lors de la suppression de l'utilisateur : %v", err)
		ginresponse.InternalServerError(c, "Erreur lors de la suppression de l'utilisateur", err.Error())
		return
	}

	ginresponse.NoContent(c, "Utilisateur supprimé avec succès")
}
