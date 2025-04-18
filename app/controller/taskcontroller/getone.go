package taskcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
)

func (co *taskController) GetOneById(c *gin.Context) {
	id := c.Param("id")
	task, err := co.taskService.GetOneById(c.Request.Context(), id)
	if err != nil {
		ginresponse.NotFound(c, "Une erreur est survenue.", []ginresponse.ErrorModel{{
			Message: "Une erreur est survenue.",
			Type: "task",
			Field: "id",
			Detail: err.Error(),
		}})
		return 
	}

	ginresponse.Success(c, "Récupération de la tâche avec succès", task)
}