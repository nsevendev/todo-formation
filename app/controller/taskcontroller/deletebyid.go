package taskcontroller

import (
	"todof/internal/task"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t *taskController) DeleteById(c *gin.Context) {
	var idsDto task.TaskDeleteManyDto
	if err := c.ShouldBindJSON(&idsDto); err != nil {
		logger.Ef("Erreur lors de la validation des ids: %s", err.Error())
		ginresponse.BadRequest(c, "Erreur de validation des ids", "Erreur de validation des ids")
		return
	}

	var idsObjectIds []primitive.ObjectID
	for _, v := range idsDto.Ids {
		objectId, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			logger.Ef("Conversion id de la tâche impossible: %s", err.Error())
			ginresponse.BadRequest(c, "L'id de la tâche est invalide", err.Error())
			return
		}
		idsObjectIds = append(idsObjectIds, objectId)
	}
	
	if err := t.taskService.DeleteById(c, idsObjectIds); err != nil {
		logger.Ef("Erreur lors de la suppression des tâches : %s", err.Error())
		ginresponse.InternalServerError(c, "Impossible de supprimer les tâches", err.Error())
		return
	}

	ginresponse.Success(c, "Les tâches ont été supprimées avec succès", nil)
}