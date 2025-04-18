package taskcontroller

import (
	"todof/internal/taskservice"

	"github.com/gin-gonic/gin"
)

type taskController struct {
	taskService taskservice.TaskServiceInterface
}

type TaskControllerInterface interface {
	GetOneById(c *gin.Context)
}

func NewTaskController(taskService taskservice.TaskServiceInterface) TaskControllerInterface {
	return &taskController{
		taskService: taskService,
	}
}