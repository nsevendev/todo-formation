package taskcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteOneByUser godoc
// @Summary Supprime une task spécifique de l'utilisateur connecté
// @Description Supprime la task indiqué de l'utilisateur authentifié via le token dans le header
// @Tags task
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID de la task à supprimer"
// @Success 200 {object} ginresponse.JsonFormatterSwag "Tâche supprimée avec succès"
// @Failure 401 {object} ginresponse.JsonFormatterSwag "Token invalide"
// @Failure 500 {object} ginresponse.JsonFormatterSwag "Erreur interne"
// @Router /task/{id}/user [delete]
func (t *taskController) DeleteOneByUser(c *gin.Context) {
	id := c.Param("id")
	taskId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Ef("Erreur lors de la conversion de l'ID de la tâche : %s", err.Error())
		ginresponse.BadRequest(c, "Id de la tâche invalide", err.Error())
		return
	}

	idUser := t.userService.GetIdUserInContext(c)

	if err := t.taskService.DeleteOneByUser(c, idUser, taskId); err != nil {
		logger.Ef("Erreur lors de la suppression de la tâche : %s", err.Error())
		ginresponse.InternalServerError(c, "Erreur lors de la suppression de la tâche", err.Error())
		return
	}

	ginresponse.Success(c, "Tâche supprimée avec succès", nil)
}
