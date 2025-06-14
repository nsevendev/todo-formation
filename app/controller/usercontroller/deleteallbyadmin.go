package usercontroller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

// DeleteAllByAdmin godoc
// @Summary Supprime les utilisateurs par un admin
// @Description Supprime les utilisateurs par un utilisateur avec role = admin
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} ginresponse.JsonFormatterSwag "X utilisateurs supprimés"
// @Failure 401 {object} ginresponse.JsonFormatterSwag "Token invalide"
// @Failure 403 {object} ginresponse.JsonFormatterSwag "Insufficient permissions"
// @Failure 500 {object} ginresponse.JsonFormatterSwag "Erreur interne"
// @Router /user/users/all [delete]
func (u *userController) DeleteAllByAdmin(c *gin.Context) {

	deletedCount, err := u.userService.DeleteAllByAdmin(c)
	if err != nil {
		logger.Ef("Erreur lors de la suppression des utilisateurs : %v", err)
		ginresponse.InternalServerError(c, "Impossible de supprimer les utilisateurs", err.Error())
		return
	}

	succesMessage := fmt.Sprintf("%v utilisateurs supprimés", deletedCount)
	ginresponse.Success(c, succesMessage, deletedCount)
}
