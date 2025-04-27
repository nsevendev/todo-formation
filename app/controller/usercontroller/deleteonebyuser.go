package usercontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

func (u *userController) DeleteUser(c *gin.Context) {
	idUser := u.userService.GetIdUserInContext(c)

	if err := u.userService.DeleteOneByUser(c, idUser); err != nil {
		logger.Ef("Erreur lors de la suppression de l'utilisateur : %v", err)
		ginresponse.InternalServerError(c, "Erreur lors de la suppression de l'utilisateur", err.Error())
		return
	}

	ginresponse.NoContent(c, "Utilisateur supprimé avec succès")
}