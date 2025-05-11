package taskcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

// DeleteAllTasks godoc
// @Summary Supprime les tâches par un admin
// @Description Supprime les tâches par un utilisateur avec role = admin
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} ginresponse.JsonFormatterSwag "X tâches supprimés"
// @Failure 401 {object} ginresponse.JsonFormatterSwag "Token invalide"
// @Failure 403 {object} ginresponse.JsonFormatterSwag "Insufficient permissions"
// @Failure 500 {object} ginresponse.JsonFormatterSwag "Erreur interne"
// @Router /task/delete/all [delete]
func (t *taskController) DeleteAllTasks(c *gin.Context) {
	err := t.taskService.DeleteAllTasks(c)
	if err != nil {
		logger.Ef("Impossible de supprimer les tâches: %s", err.Error())
		ginresponse.InternalServerError(c, "Impossible de supprimer les tâches", err.Error())
		return
	}

	ginresponse.Success(c, "Les tâches ont été supprimées avec succès", nil)
}
