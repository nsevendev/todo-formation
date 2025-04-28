package usercontroller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

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
