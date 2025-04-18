package models

import (
	"context"
	"errors"
	initializer "todof/internal/init"
	"todof/mod/mongotool/mongodate"

	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Task struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Label	 string             `bson:"label" json:"label"`
	Done	 bool               `bson:"done" json:"done"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`

	IDUser primitive.ObjectID `bson:"id_user" json:"id_user"`
}

type taskModel struct {
	col *mongo.Collection
}

type TaskModelInterface interface {
	Create(ctx context.Context, task *Task) error
	GetOneByID(ctx context.Context, IDTask string) (*Task, error)
	GetAllByUser(ctx context.Context, IDUser primitive.ObjectID) ([]Task, error)
	UpdateLabel(ctx context.Context, IDTask primitive.ObjectID, IDUser primitive.ObjectID, newLabel string) error
	UpdateDone(ctx context.Context, IDTask primitive.ObjectID, IDUser primitive.ObjectID, done bool) error
	DeleteOneByID(ctx context.Context, IDTask primitive.ObjectID, IDUser primitive.ObjectID) error
	DeleteAllByUser(ctx context.Context, IDUser primitive.ObjectID) error
	DeleteAll(ctx context.Context,) error
}

func newTaskModel() TaskModelInterface {
	return &taskModel{
		col: initializer.Db.Collection("tasks"),
	}
}

func (t *taskModel) Create(ctx context.Context, task *Task) error {
	task.ID = primitive.NewObjectID()
	task.CreatedAt = mongodate.Now()
	task.UpdatedAt = mongodate.Now()

	_, err := t.col.InsertOne(ctx, task)
	if err != nil {
		logger.Ef("impossible de créer la tâche _id: %s, id_user: %v", task.ID.Hex(), task.IDUser.Hex())
		return errors.New("impossible de créer la tâche")
	}
	return nil
}

func (t *taskModel) GetOneByID(ctx context.Context, IDTask string) (*Task, error) {
	var task Task
	if err := t.col.FindOne(ctx, bson.M{"_id": IDTask}).Decode(&task); err == mongo.ErrNoDocuments {
		logger.Ef("impossible de trouver la tâche _id: %s", IDTask)
		return nil, errors.New("une erreur est survenue")
	}
	return &task, nil
}

func (t *taskModel) GetAllByUser(ctx context.Context, IDUser primitive.ObjectID) ([]Task, error) {
	cur, err := t.col.Find(ctx, bson.M{"id_user": IDUser})
	if err != nil {
		logger.Ef("impossible de trouver les tâches de l'utilisateur id_user: %v", IDUser.Hex())
		return nil, err
	}
	defer cur.Close(ctx)

	var tasks []Task
	for cur.Next(ctx) {
		var task Task
		if err := cur.Decode(&task); err != nil {
			logger.Ef("impossible de décoder la tâche id_user: %v", IDUser.Hex())
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (t *taskModel) UpdateLabel(ctx context.Context, IDTask primitive.ObjectID, IDUser primitive.ObjectID, newLabel string) error {
	res, err := t.col.UpdateOne(ctx,
		bson.M{"_id": IDTask, "id_user": IDUser},
		bson.M{
			"$set": bson.M{
				"label":      newLabel,
				"updated_at": mongodate.Now(),
			},
		},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		logger.Ef("impossible de modifier la tâche _id: %s, id_user: %v", IDTask.Hex(), IDUser.Hex())
		return errors.New("imposible de modifier la tâche")
	}
	return err
}

func (t *taskModel) UpdateDone(ctx context.Context, IDTask primitive.ObjectID, IDUser primitive.ObjectID, done bool) error {
	res, err := t.col.UpdateOne(ctx,
		bson.M{"_id": IDTask, "id_user": IDUser},
		bson.M{
			"$set": bson.M{
				"done":       done,
				"updated_at": mongodate.Now(),
			},
		},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		logger.Ef("impossible de modifier la tâche _id: %s, id_user: %v", IDTask.Hex(), IDUser.Hex())
		return errors.New("imposible de modifier la tâche")
	}
	return nil
}

func (t *taskModel) DeleteOneByID(ctx context.Context, IDTask primitive.ObjectID, IDUser primitive.ObjectID) error {
	_, err := t.col.DeleteOne(ctx, bson.M{"_id": IDTask, "id_user": IDUser})
	if err != nil {
		logger.Ef("impossible de supprimer la tâche _id: %s, id_user: %v", IDTask.Hex(), IDUser.Hex())
		return errors.New("impossible de supprimer la tâche")
	}
	return nil
}

func (t *taskModel) DeleteAllByUser(ctx context.Context, IDUser primitive.ObjectID) error {
	_, err := t.col.DeleteMany(ctx, bson.M{"id_user": IDUser})
	if err != nil {
		logger.Ef("impossible de supprimer les tâches id_user: %s", IDUser.Hex())
		return errors.New("impossible de supprimer les tâches")
	}
	return nil
}

func (t *taskModel) DeleteAll(ctx context.Context,) error {
	_, err := t.col.DeleteMany(ctx, bson.M{})
	if err != nil {
		logger.Ef("impossible de supprimer toutes les tâches")
		return errors.New("impossible de supprimer toutes les tâches")
	}
	return err
}
