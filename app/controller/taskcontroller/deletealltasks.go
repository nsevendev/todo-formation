package taskcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

func (t *taskController) DeleteAllTasks(c *gin.Context) {
	err := t.taskService.DeleteAllTasks(c)
	if err != nil {
		logger.Ef("Impossible de supprimer les tâches: %s", err.Error())
		ginresponse.InternalServerError(c, "Impossible de supprimer les tâches", err.Error())
		return
	}

	ginresponse.Success(c, "Les tâches ont été supprimées avec succès", nil)
}
