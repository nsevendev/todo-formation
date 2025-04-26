package taskcontroller

import (
	"fmt"
	"todof/internal/task"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t *taskController) UpdateOneLabelPropertyByUser(c *gin.Context) {
	var taskUpdateDto task.TaskUpdateDto
	if err := c.ShouldBindJSON(&taskUpdateDto); err != nil {
		logger.Ef("Erreur de validation : %s", err.Error())
		ginresponse.BadRequest(c, "Erreur de validation", ginresponse.ErrorModel{
			Message: err.Error(),
			Type: "Validation",
			Detail: fmt.Sprintf("%v", err),
		})
		return
	}

	id := c.Param("id")
	taskId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
		logger.Ef("impossible de modifier la tâche : %s", err.Error())
        ginresponse.BadRequest(c, "Id de la tâche invalide", err.Error())
        return
    }

	idUser := t.userService.GetIdUserInContext(c)
	
	if err := t.taskService.UpdateOneLabelPropertyByUser(c, idUser, taskId, taskUpdateDto); err != nil {
		logger.Ef("Erreur lors de la mise à jour de la tâche : %s", err.Error())
		ginresponse.InternalServerError(c, "Erreur lors de la mise à jour de la tâche", err.Error())
		return
	}

	ginresponse.Success(c, "Tâche mise à jour avec succès", nil)
}