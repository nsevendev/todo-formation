package models

import (
	"context"
)

type Manager struct {
	TaskModel TaskInterface
}

func NewManager(ctx context.Context) *Manager {
	taskModel := NewTaskModel(ctx)

	return &Manager{
		TaskModel: taskModel,
	}
}