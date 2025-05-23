package taskcontroller

import (
	"todof/internal/auth"
	"todof/internal/task"

	"github.com/gin-gonic/gin"
)

type taskController struct {
	taskService task.TaskServiceInterface
	userService  auth.UserServiceInterface
}

type TaskControllerInterface interface {
	Create(c *gin.Context)
	GetAllByUser(c *gin.Context)
	UpdateOneDonePropertyByUser(c *gin.Context)
	UpdateOneLabelPropertyByUser(c *gin.Context)
	DeleteOneByUser(c *gin.Context)
	DeleteManyByUser(c *gin.Context)
	DeleteById(c *gin.Context)
	DeleteAllTasks(c *gin.Context)
}

func NewTaskController(taskService task.TaskServiceInterface, userService auth.UserServiceInterface) TaskControllerInterface {
	return &taskController{
		taskService,
		userService,
	}
}