package taskcontroller

import (
	"fmt"
	"todof/internal/task"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create godoc
// @Summary Créer une task
// @Description Création d’une nouvelle task pour l'utilisateur connecté
// @Tags task
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body task.TaskCreateDto true "DTO de création de la task"
// @Success 201 {object} doc.ResponseModel "Tâche créée avec succès"
// @Failure 401 {object} doc.ResponseModel "Invalide token"
// @Failure 500 {object} doc.ResponseModel "Erreur interne"
// @Router /task [post]
func (t *taskController) Create(c *gin.Context) {
	var taskCreateDto task.TaskCreateDto
	if err := c.ShouldBindJSON(&taskCreateDto); err != nil {
		logger.Ef("Erreur de validation : %s", err.Error())
		ginresponse.BadRequest(c, "Erreur de validation", ginresponse.ErrorModel{
			Message: err.Error(),
			Type: "Validation",
			Detail: fmt.Sprintf("%v", err),
		})
		return
	}

	idUser, exists := c.Get("id_user")
	if !exists {
		logger.Ef("Erreur d'authentification : ID utilisateur non trouvé dans le contexte")
		ginresponse.Unauthorized(c, "Erreur d'authentification", "Vous n'avez pas les droits pour effectuer cette action")
		return
	}

	taskCreated, err := t.taskService.Create(c, taskCreateDto, idUser.(primitive.ObjectID))
	if err != nil {
		logger.Ef("Erreur de création de tâche : %s", err.Error())
		ginresponse.InternalServerError(c, err.Error(), err.Error())
		return
	}

	ginresponse.Created(c, "Tâche créée avec succès", taskCreated)
}