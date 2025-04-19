package task

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskService struct {
	taskRepo taskRepoInterface
}

type TaskServiceInterface interface {
	Create(ctx context.Context, taskCreateDto TaskCreateDto, IdUser primitive.ObjectID) (*Task, error)
	GetAllByUser(ctx context.Context, idUser primitive.ObjectID) ([]Task, error)
}

func NewTaskService(taskRepo taskRepoInterface) TaskServiceInterface {
	return &taskService{
		taskRepo: taskRepo,
	}
}

func (t *taskService) Create(ctx context.Context, taskCreateDto TaskCreateDto, IdUser primitive.ObjectID) (*Task, error) {
	task := &Task{
		Label: taskCreateDto.Label,
		Done:  false,
		IdUser: IdUser,
	}

	if err := t.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	
	return task, nil
}

func (t *taskService) GetAllByUser(ctx context.Context, idUser primitive.ObjectID) ([]Task, error) {
	tasks, err := t.taskRepo.GetAllByUser(ctx, idUser)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

