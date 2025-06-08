package taskcontroller

import (
	"fmt"
	"todof/internal/task"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateOneLabelPropertyByUser godoc
// @Summary Met à jour le label d'une task spécifique de l'utilisateur connecté
// @Description Met à jour la propriété `label` de la task indiquée appartenant à l'utilisateur authentifié via le token dans le header
// @Tags task
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID de la task à modifier"
// @Param label body task.TaskUpdateLabelDto true "Label de mise à jour"
// @Success 200 {object} ginresponse.JsonFormatterSwag "Tâche mise à jour avec succès"
// @Failure 401 {object} ginresponse.JsonFormatterSwag "Token invalide"
// @Failure 500 {object} ginresponse.JsonFormatterSwag "Erreur interne"
// @Router /task/{id}/label/user [put]
func (t *taskController) UpdateOneLabelPropertyByUser(c *gin.Context) {
	var taskUpdateDto task.TaskUpdateLabelDto
	if err := c.ShouldBindJSON(&taskUpdateDto); err != nil {
		logger.Ef("Erreur de validation : %s", err.Error())
		ginresponse.BadRequest(c, "Erreur de validation", ginresponse.ErrorModel{
			Message: err.Error(),
			Type:    "Validation",
			Detail:  fmt.Sprintf("%v", err),
		})
		return
	}

	id := c.Param("id")
	taskId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Ef("impossible de modifier la tâche : %s", err.Error())
		ginresponse.BadRequest(c, "impossible de modifier la tâche", err.Error())
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
