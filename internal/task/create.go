package task

import (
	"todof/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t *taskService) Create(label string, IDUser primitive.ObjectID) (*models.Task, error) {
	task := &models.Task{
		Label: label,
		Done:  false,
		IDUser: IDUser,
	}

	if err := t.taskModel.Create(task, IDUser); err != nil {
		return nil, err
	}
	
	return task, nil
}