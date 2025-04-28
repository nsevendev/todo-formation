package task

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskService struct {
	taskRepo taskRepoInterface
	db *mongo.Database
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

func NewTaskService(taskRepo taskRepoInterface, db *mongo.Database) TaskServiceInterface {
	return &taskService{
		taskRepo: taskRepo,
		db: db,
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
	userCollection := t.db.Collection("users")
	filter := bson.M{"role": bson.M{"$ne": "admin"}}
	cursor, err := userCollection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	
	var userIDs []primitive.ObjectID
	if err := cursor.All(ctx, &userIDs); err != nil {
		return err
	}

	if cursor.Err() != nil {
		return cursor.Err()
	}

	if err := t.taskRepo.DeleteAllTasks(ctx, userIDs); err != nil {
		return err
	}

	return nil
}