package taskcontroller

import (
	"todof/internal/task"

	"github.com/gin-gonic/gin"
)

type taskController struct {
	taskService task.TaskServiceInterface
}

type TaskControllerInterface interface {
	Create(c *gin.Context)
	GetAllByUser(c *gin.Context)
}

func NewTaskController(taskService task.TaskServiceInterface) TaskControllerInterface {
	return &taskController{
		taskService,
	}
}