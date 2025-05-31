package task

import (
	"context"

	"todof/internal/auth"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskService struct {
	taskRepo TaskRepoInterface
	userRepo auth.UserRepoInterface
}

type TaskServiceInterface interface {
	Create(ctx context.Context, taskCreateDto TaskCreateDto, IdUser primitive.ObjectID) (*Task, error)
	GetAllByUser(ctx context.Context, idUser primitive.ObjectID) ([]Task, error)
	UpdateOneDonePropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID) error
	UpdateOneLabelPropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID, taskUpdateDto TaskUpdateLabelDto) error
	DeleteOneByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID) error
	DeleteManyByUser(ctx context.Context, idUser primitive.ObjectID, ids []primitive.ObjectID) error
	DeleteById(ctx context.Context, ids []primitive.ObjectID) error
	DeleteAllTasks(ctx context.Context) error
}

func NewTaskService(taskRepo TaskRepoInterface, userRepo auth.UserRepoInterface) TaskServiceInterface {
	return &taskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

func (t *taskService) Create(ctx context.Context, taskCreateDto TaskCreateDto, IdUser primitive.ObjectID) (*Task, error) {
	task := &Task{
		Label:  taskCreateDto.Label,
		Done:   false,
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

func (t *taskService) UpdateOneLabelPropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID, taskUpdateDto TaskUpdateLabelDto) error {
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

func (t *taskService) DeleteById(ctx context.Context, ids []primitive.ObjectID) error {
	if err := t.taskRepo.DeleteById(ctx, ids); err != nil {
		return err
	}

	return nil
}

func (t *taskService) DeleteAllTasks(ctx context.Context) error {
	users, err := t.userRepo.FindNonAdmin(ctx)
	if err != nil {
		return err
	}

	var userIDs []primitive.ObjectID
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	if err := t.taskRepo.DeleteAllTasks(ctx, userIDs); err != nil {
		return err
	}

	return nil
}
