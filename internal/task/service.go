package task

import (
	"todof/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskService struct {
	taskModel models.TaskModelInterface
}

type TaskServiceInterface interface {
	Create(task *models.Task, IDUser primitive.ObjectID) error
}

func NewTaskService(modelManager *models.ManagerModel) TaskServiceInterface {
	taskModel := modelManager.TaskModel

	return &taskService{taskModel}
}

