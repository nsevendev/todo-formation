package taskcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateOneDonePropertyByUser godoc
// @Summary Met à jour la propriété done d'une task spécifique de l'utilisateur connecté
// @Description Met à jour la propriété `done` de la task indiquée appartenant à l'utilisateur authentifié via le token dans le header
// @Tags task
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID de la task à modifier"
// @Success 200 {object} doc.ResponseModel "Tâche mise à jour avec succès"
// @Failure 401 {object} doc.ResponseModel "Token invalide"
// @Failure 500 {object} doc.ResponseModel "Erreur interne"
// @Router /task/{id}/done/user [put]
func (t *taskController) UpdateOneDonePropertyByUser(c *gin.Context) {
	id := c.Param("id")
	taskId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
		logger.Ef("impossible de modifier la tâche : %s", err.Error())
        ginresponse.BadRequest(c, "impossible de modifier la tâche", err.Error())
        return
    }

	idUser := t.userService.GetIdUserInContext(c)
	
	if err := t.taskService.UpdateOneDonePropertyByUser(c, idUser, taskId); err != nil {
		logger.Ef("Erreur lors de la mise à jour de la tâche : %s", err.Error())
		ginresponse.InternalServerError(c, "Erreur lors de la mise à jour de la tâche", err.Error())
		return
	}

	ginresponse.Success(c, "Tâche mise à jour avec succès", nil)
}