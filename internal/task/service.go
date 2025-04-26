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
	UpdateOneDonePropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID) error
	UpdateOneLabelPropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID, taskUpdateDto TaskUpdateDto) error
	DeleteOneByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID) error
	DeleteManyByUser(ctx context.Context, idUser primitive.ObjectID, ids []primitive.ObjectID) error
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

	if tasks == nil {
		return []Task{}, nil
	}

	return tasks, nil
}

func (t *taskService) UpdateOneDonePropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID) error {
	if err := t.taskRepo.UpdateOneDonePropertyByUser(ctx, idUser, idTask); err != nil {
		return err
	}

	return nil
}

func (t *taskService) UpdateOneLabelPropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID, taskUpdateDto TaskUpdateDto) error {
	if err := t.taskRepo.UpdateOneLabelPropertyByUser(ctx, idUser, idTask, taskUpdateDto.Label); err != nil {
		return err
	}

	return nil
}

func (t *taskService) DeleteOneByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID) error {
	if err := t.taskRepo.DeleteOneByUser(ctx, idUser, idTask); err != nil {
		return err
	}

	return nil
}

func (t *taskService) DeleteManyByUser(ctx context.Context, idUser primitive.ObjectID, ids []primitive.ObjectID) error {
	if err := t.taskRepo.DeleteManyByUser(ctx, idUser, ids); err != nil {
		return err
	}

	return nil
}
