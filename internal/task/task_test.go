package task

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
	"todof/internal/auth"
	initializer "todof/internal/init"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var s TaskServiceInterface
var userCollection *mongo.Collection
var taskCollection *mongo.Collection
var ctx context.Context
var usersIds []primitive.ObjectID
var tasksIds []primitive.ObjectID

//service
func TestMain(m *testing.M) {
	taskCollection = initializer.Db.Collection("tasks")
	userCollection = initializer.Db.Collection("users")
	r := NewTaskRepo(initializer.Db)
	userRepo := auth.NewUserRepo(initializer.Db)
	s = NewTaskService(r, userRepo)
	ctx := context.Background()

	if _, err := taskCollection.DeleteMany(ctx, bson.M{}); err != nil {
		log.Fatalf("Erreur lors du nettoyage de la collection tasks : %v", err)
	}

	if _, err := userCollection.DeleteMany(ctx, bson.M{}); err != nil {
		log.Fatalf("Erreur lors du nettoyage de la collection users : %v", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (userID primitive.ObjectID)
		label string
		isErr bool
	}{
		{
			name:  "test success",
			label: "task test",
			isErr: false,
			setup: func() primitive.ObjectID {
				user := &auth.User{
					Email: "taskTest@gmail.com",
					Password: "password",
				}

				result, err := userCollection.InsertOne(ctx, user)
				if err != nil {
					t.Fatalf("Erreur lors de la création du user: %v", err)
				}
				usersIds = append(usersIds, result.InsertedID.(primitive.ObjectID))
				return result.InsertedID.(primitive.ObjectID)
			},
		},
	}

	for _, tt := range tests {
		userId := tt.setup()

		createDto := TaskCreateDto{
			Label: tt.label,
		}

		task, err := s.Create(ctx, createDto, userId)

		if task != nil {
			tasksIds = append(tasksIds, task.ID)
		}

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if err == nil && task.Label != tt.label {
			t.Errorf("%s: got label %s, expect label %s", tt.name, task.Label, tt.label)
		}
	}
}

func TestGetAllByUser(t *testing.T){
	userId := func() primitive.ObjectID {
		user := &auth.User{
			Email: "taskTest2@gmail.com",
			Password: "password",
		}

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			t.Fatalf("Erreur lors de la création du user: %v", err)
		}
		usersIds = append(usersIds, result.InsertedID.(primitive.ObjectID))
		return result.InsertedID.(primitive.ObjectID)
	}

	tests := []struct {
		name string
		userId primitive.ObjectID
		isTask bool
		isErr bool
	}{
		{"test success", usersIds[0], true, false},
		{"test avec user sans task", userId(), true, false},
	}

	for _, tt := range tests {
		task, err := s.GetAllByUser(ctx, tt.userId)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (task == nil) == tt.isTask{
			t.Errorf("%s: got task %v, expect task %v", tt.name, task, tt.isTask)
		}
	}
}

func TestUpdateOneDonePropertyByUser(t *testing.T){
	tests := []struct {
		name string
		idUser primitive.ObjectID
		idTask primitive.ObjectID
		isErr bool
	}{
		{"test success", usersIds[0], tasksIds[0], false},
	}

	for _, tt := range tests {
		err := s.UpdateOneDonePropertyByUser(ctx, tt.idUser, tt.idTask)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}
	}
}

func TestUpdateOneLabelPropertyByUser(t *testing.T){
	tests := []struct {
		name string
		idUser primitive.ObjectID
		idTask primitive.ObjectID
		label string
		isErr bool
	}{
		{"test success", usersIds[0], tasksIds[0], "label updated", false},
	}

	for _, tt := range tests {
		updateDto := TaskUpdateLabelDto{
			Label: tt.label,
		}

		err := s.UpdateOneLabelPropertyByUser(ctx, tt.idUser, tt.idTask, updateDto)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}
	}
}

func TestDeleteOneByUser(t *testing.T){
	tests := []struct {
		name string
		idUser primitive.ObjectID
		idTask primitive.ObjectID
		isErr bool
	}{
		{"test success", usersIds[0], tasksIds[0], false},
	}

	for _, tt := range tests {
		err := s.DeleteOneByUser(ctx, tt.idUser, tt.idTask)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}
	}
}

func TestDeleteManyByUser(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (userID primitive.ObjectID, taskID primitive.ObjectID)
		isErr bool
	}{
		{
			name: "test success",
			setup: func() (primitive.ObjectID, primitive.ObjectID) {
				user := &auth.User{
					Email:    "taskTest3@gmail.com",
					Password: "password",
				}

				result, err := userCollection.InsertOne(ctx, user)

				if err != nil {
					t.Fatalf("Erreur lors de la création du user: %v", err)
				}

				userID := result.InsertedID.(primitive.ObjectID)

				task := &Task{
					Label:    "Test task",
					Done:   false,
					CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
					UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
					IdUser: userID,
				}

				res, err := taskCollection.InsertOne(ctx, task)

				if err != nil {
					t.Fatalf("Erreur lors de la création de la tâche: %v", err)
				}
				taskID := res.InsertedID.(primitive.ObjectID)

				return userID, taskID
			},
			isErr: false,
		},
	}

	for _, tt := range tests {
		userID, taskID := tt.setup()

		err := s.DeleteManyByUser(ctx, userID, []primitive.ObjectID{taskID})

		if (err != nil) != tt.isErr {
			t.Errorf("got error %v, expected error: %v", err, tt.isErr)
		}
	}
}
