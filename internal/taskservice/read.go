package taskservice

import (
	"context"
	"todof/internal/models"
)


func (t *taskService) GetOneById(ctx context.Context, id string) (*models.Task, error) {
	task, err := t.taskModel.GetOneByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}