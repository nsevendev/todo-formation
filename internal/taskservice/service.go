package taskservice

import (
	"context"
	"todof/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskService struct {
	taskModel models.TaskModelInterface
}

type TaskServiceInterface interface {
	Create(ctx context.Context, label string, IDUser primitive.ObjectID) (*models.Task, error)
	GetOneById(ctx context.Context, id string) (*models.Task, error)
}

func NewTaskService(modelManager *models.ManagerModel) TaskServiceInterface {
	return &taskService{
		taskModel: modelManager.TaskModel,
	}
}

