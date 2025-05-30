package task

import (
	"context"
	"errors"

	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepo struct {
	collection *mongo.Collection
}

type TaskRepoInterface interface {
	Create(ctx context.Context, task *Task) error
	GetAllByUser(ctx context.Context, idUser primitive.ObjectID) ([]Task, error)
	UpdateOneDonePropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID) error
	UpdateOneLabelPropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID, labelUpdate string) error
	DeleteOneByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID) error
	DeleteManyByUser(ctx context.Context, idUser primitive.ObjectID, ids []primitive.ObjectID) error
	DeleteById(ctx context.Context, ids []primitive.ObjectID) error
	DeleteAllTasks(ctx context.Context, userIDs []primitive.ObjectID) error
}

func NewTaskRepo(db *mongo.Database) TaskRepoInterface {
	return &taskRepo{
		collection: db.Collection("tasks"),
	}
}

func (t *taskRepo) Create(ctx context.Context, task *Task) error {
	task.ID = primitive.NewObjectID()
	task.SetTimeStamps()

	_, err := t.collection.InsertOne(ctx, task)
	if err != nil {
		logger.Ef("impossible de créer la tâche _id: %s, id_user: %v", task.ID.Hex(), task.IdUser.Hex())
		return errors.New("impossible de créer la tâche : " + err.Error())
	}
	return nil
}

func (t *taskRepo) GetAllByUser(ctx context.Context, idUser primitive.ObjectID) ([]Task, error) {
	cursor, err := t.collection.Find(ctx, bson.M{"id_user": idUser})
	if err != nil {
		logger.Ef("impossible de récupérer les tâches de l'utilisateur _id: %s", idUser.Hex())
		return nil, errors.New("impossible de récupérer les tâches")
	}
	defer cursor.Close(ctx)

	var tasks []Task
	if err := cursor.All(ctx, &tasks); err != nil {
		logger.Ef("impossible de récupérer les tâches de l'utilisateur _id: %s", idUser.Hex())
		return nil, errors.New("impossible de récupérer les tâches")
	}

	return tasks, nil
}

func (t *taskRepo) UpdateOneDonePropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID) error {
	filter := bson.M{"_id": idTask, "id_user": idUser}

	var task Task

	if err := t.collection.FindOne(ctx, filter).Decode(&task); err != nil {
		logger.Ef("impossible de trouver la tâche _id: %s, id_user: %s, error: %s", idTask.Hex(), idUser.Hex(), err.Error())
		return errors.New("aucune tâche trouvée")
	}

	update := bson.M{
		"$set": bson.M{
			"done": !task.Done,
		},
	}

	result, err := t.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Ef("Erreur lors de la mise à jour de la tâche _id: %s, id_user: %s, error: %s", idTask.Hex(), idUser.Hex(), err.Error())
		return errors.New("impossible de mettre à jour la tâche")
	}

	if result.MatchedCount == 0 {
		logger.Ef("Aucune tâche modifié _id: %s, id_user: %s", idTask.Hex(), idUser.Hex())
		return errors.New("aucune tâche mise à jour")
	}

	logger.Sf("tâche mise à jour _id: %s, id_user: %s", idTask.Hex(), idUser.Hex())

	return nil
}

func (t *taskRepo) UpdateOneLabelPropertyByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID, labelUpdate string) error {
	filter := bson.M{"_id": idTask, "id_user": idUser}

	var task Task

	if err := t.collection.FindOne(ctx, filter).Decode(&task); err != nil {
		logger.Ef("impossible de trouver la tâche _id: %s, id_user: %s, error: %s", idTask.Hex(), idUser.Hex(), err.Error())
		return errors.New("aucune tâche trouvée")
	}

	update := bson.M{
		"$set": bson.M{
			"label": labelUpdate,
		},
	}

	result, err := t.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Ef("Erreur lors de la mise à jour de la tâche _id: %s, id_user: %s, error: %s", idTask.Hex(), idUser.Hex(), err.Error())
		return errors.New("impossible de mettre à jour la tâche")
	}

	if result.MatchedCount == 0 {
		logger.Ef("Aucune tâche modifié _id: %s, id_user: %s", idTask.Hex(), idUser.Hex())
		return errors.New("aucune tâche mise à jour")
	}

	logger.Sf("tâche mise à jour _id: %s, id_user: %s", idTask.Hex(), idUser.Hex())

	return nil
}

func (t *taskRepo) DeleteOneByUser(ctx context.Context, idUser primitive.ObjectID, idTask primitive.ObjectID) error {
	filter := bson.M{"_id": idTask, "id_user": idUser}

	result, err := t.collection.DeleteOne(ctx, filter)
	if err != nil {
		logger.Ef("impossible de supprimer la tâche _id: %s, id_user: %s", idTask.Hex(), idUser.Hex())
		return errors.New("impossible de supprimer la tâche")
	}

	if result.DeletedCount == 0 {
		logger.Ef("aucune tâche supprimée _id: %s, id_user: %s", idTask.Hex(), idUser.Hex())
		return errors.New("aucune tâche supprimée")
	}

	logger.Sf("tâche supprimée _id: %s, id_user: %s", idTask.Hex(), idUser.Hex())

	return nil
}

func (t *taskRepo) DeleteManyByUser(ctx context.Context, idUser primitive.ObjectID, ids []primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": ids}, "id_user": idUser}

	result, err := t.collection.DeleteMany(ctx, filter)
	if err != nil {
		logger.Ef("impossible de supprimer les tâches _id: %s, id_user: %s", ids, idUser.Hex())
		return errors.New("impossible de supprimer les tâches")
	}

	if result.DeletedCount == 0 {
		logger.Ef("aucune tâche supprimée _id: %s, id_user: %s", ids, idUser.Hex())
		return errors.New("aucune tâche supprimée")
	}

	logger.Sf("tâches supprimées _id: %s, id_user: %s", ids, idUser.Hex())

	return nil
}

func (t *taskRepo) DeleteById(ctx context.Context, ids []primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	result, err := t.collection.DeleteMany(ctx, filter)
	if err != nil {
		logger.Ef("impossible de supprimer les tâches _id: %v", ids)
		return errors.New("impossible de supprimer les tâches")
	}

	if result.DeletedCount == 0 {
		logger.Ef("aucune tâche supprimée _id: %v", ids)
		return errors.New("aucune tâche supprimée")
	}

	logger.Sf("tâches supprimées _id: %v", ids)
	return nil
}

func (t *taskRepo) DeleteAllTasks(ctx context.Context, userIDs []primitive.ObjectID) error {
	_, err := t.collection.DeleteMany(ctx, bson.M{"id_user": bson.M{"$in": userIDs}})
	if err != nil {
		logger.Ef("impossible de supprimer les tâches pour les users: %v", userIDs)
		return errors.New("impossible de supprimer les tâches")
	}

	return nil
}
