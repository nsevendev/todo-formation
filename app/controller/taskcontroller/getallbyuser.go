package taskcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t *taskController) GetAllByUser(c *gin.Context) {
	idUser, exists := c.Get("id_user")
	if !exists {
		logger.Ef("Erreur d'authentification : ID utilisateur non trouvé dans le contexte")
		ginresponse.Unauthorized(c, "Erreur d'authentification", "Vous n'avez pas les droits pour effectuer cette action")
		return
	}

	
	tasks, err := t.taskService.GetAllByUser(c, idUser.(primitive.ObjectID))
	if err != nil {
		logger.Ef("Erreur lors de la récupération des tâches : %v", err)
		ginresponse.InternalServerError(c, "Erreur interne", "Une erreur est survenue lors de la récupération des tâches")
		return
	}

	logger.If("Tâches récupérées avec succès pour l'utilisateur : %v", tasks)

	ginresponse.Success(c, "Tâches récupérées avec succès", tasks)
}