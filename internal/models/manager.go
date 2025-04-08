package models

import (
	"context"
)

type Manager struct {
	TaskModel TaskModelInterface
}

func NewManager(ctx context.Context) *Manager {
	taskModel := NewTaskModel(ctx)

	return &Manager{
		TaskModel: taskModel,
	}
}