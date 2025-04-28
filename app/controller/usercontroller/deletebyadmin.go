package usercontroller

import (
	"todof/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

func (u *userController) DeleteByAdmin(c *gin.Context) {
	var idsDto auth.UserDeleteDto
	if err := c.ShouldBindJSON(&idsDto); err != nil {
		logger.Ef("Erreur lors de la récupération des IDs des utilisateurs : %v", err)
		ginresponse.BadRequest(c, "Erreur de validation des IDs", "Erreur de validation des IDs")
		return
	}

	deletedCount, err := u.userService.DeleteByAdmin(c, idsDto.Ids)
	if err != nil {
		logger.Ef("Erreur lors de la suppression des utilisateurs : %v", err)
		ginresponse.InternalServerError(c, "Impossible de supprimer les utilisateurs", err.Error())
		return
	}

	ginresponse.Success(c, "Utilisateurs supprimées avec succès", deletedCount)
}
