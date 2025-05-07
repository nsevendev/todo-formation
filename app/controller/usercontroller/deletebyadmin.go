package usercontroller

import (
	"fmt"
	"todof/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

// DeleteByAdmin godoc
// @Summary Supprime un ou plusieurs utilisateur(s) spécifique par un admin
// @Description Supprime le ou les utilisateur(s) indiqué par un utilisateur avec role = admin
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param ids body auth.UserDeleteDto true "Ids des utilisateurs à supprimer"
// @Success 200 {object} doc.ResponseModel "X utilisateurs supprimés"
// @Failure 401 {object} doc.ResponseModel "Token invalide"
// @Failure 403 {object} doc.ResponseModel "Insufficient permissions"
// @Failure 500 {object} doc.ResponseModel "Erreur interne"
// @Router /user/users [post]
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

	succesMessage := fmt.Sprintf("%v utilisateurs supprimés", deletedCount)
	ginresponse.Success(c, succesMessage, deletedCount)
}