package taskcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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