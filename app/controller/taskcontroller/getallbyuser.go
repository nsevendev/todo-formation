package taskcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllByUser godoc
// @Summary Récupérer toutes les tasks de l'utilisateur connecté
// @Description Récupére les tasks crée par l'utilisateur authentifié via le token dans le header
// @Tags task
// @Produce json
// @Security BearerAuth
// @Success 200 {object} doc.ResponseModel "Tâches récupérées avec succès"
// @Failure 401 {object} doc.ResponseModel "Invalide token"
// @Failure 500 {object} doc.ResponseModel "Erreur interne"
// @Router /user [get]
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

	ginresponse.Success(c, "Tâches récupérées avec succès", tasks)
}