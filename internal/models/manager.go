package models

import (
	"context"
)

type ManagerModel struct {
	TaskModel TaskModelInterface
}

func NewManager(ctx context.Context) *ManagerModel {
	taskModel := NewTaskModel(ctx)

	return &ManagerModel{
		TaskModel: taskModel,
	}
}