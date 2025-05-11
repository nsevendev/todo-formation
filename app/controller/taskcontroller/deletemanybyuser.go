package taskcontroller

import (
	"todof/internal/task"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteManyByUser godoc
// @Summary Supprime une ou plusieurs task(s) spécifique de l'utilisateur connecté
// @Description Supprime la ou les task(s) indiqué de l'utilisateur authentifié via le token dans le header
// @Tags task
// @Security BearerAuth
// @Produce json
// @Param ids body task.TaskDeleteManyDto true "Ids des tasks à supprimer"
// @Success 200 {object} ginresponse.JsonFormatterSwag "Les tâches ont été supprimées avec succès"
// @Failure 401 {object} ginresponse.JsonFormatterSwag "Token invalide"
// @Failure 500 {object} ginresponse.JsonFormatterSwag "Erreur interne"
// @Router /task/delete/user [post]
func (t *taskController) DeleteManyByUser(c *gin.Context) {
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

	idUser := t.userService.GetIdUserInContext(c)

	if err := t.taskService.DeleteManyByUser(c, idUser, idsObjectIds); err != nil {
		ginresponse.InternalServerError(c, "Impossible de supprimer les tâches", err.Error())
		return
	}

	ginresponse.Success(c, "Les tâches ont été supprimées avec succès", nil)
}
