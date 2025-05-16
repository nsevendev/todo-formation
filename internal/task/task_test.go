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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var s TaskServiceInterface
var userCollection *mongo.Collection
var taskCollection *mongo.Collection
var ctx context.Context
var cancelCtx context.Context
var cancelFunc context.CancelFunc
var usersIds []primitive.ObjectID
var tasksIds []primitive.ObjectID

func TestMain(m *testing.M) {
	taskCollection = initializer.Db.Collection("tasks")
	userCollection = initializer.Db.Collection("users")
	r := NewTaskRepo(initializer.Db)
	userRepo := auth.NewUserRepo(initializer.Db)
	s = NewTaskService(r, userRepo)
	ctx = context.Background()

	cancelCtx, cancelFunc = context.WithCancel(ctx)
	cancelFunc()

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
		{"test echec mongo", func() primitive.ObjectID {return primitive.NewObjectID()}, "label test", true},
	}

	for _, tt := range tests {
		userId := tt.setup()

		createDto := TaskCreateDto{
			Label: tt.label,
		}

		if tt.name == "test echec mongo" {
			_, err := s.Create(cancelCtx, createDto, userId)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		}else{
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
		{"test echec mongo", usersIds[0], false, true},
		{"test document mal formé", usersIds[0], false, true},
	}

	for _, tt := range tests {
		switch tt.name {
		case "test echec mongo":
			_, err := s.GetAllByUser(cancelCtx, tt.userId)
			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		case "test document mal formé":
			_, err := taskCollection.InsertOne(ctx, bson.M{
				"id_user": tt.userId,
				"label":   1234,
				"done":    false,
			}, options.InsertOne().SetBypassDocumentValidation(true))
			if err != nil {
				t.Fatalf("Erreur lors de l'insertion du document mal formé: %v", err)
			}

			_, err = s.GetAllByUser(ctx, tt.userId)
			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		default:
			task, err := s.GetAllByUser(ctx, tt.userId)
			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
			if (task == nil) == tt.isTask {
				t.Errorf("%s: got task %v, expect task %v", tt.name, task, tt.isTask)
			}
		}
	}
}

func TestUpdateOneDonePropertyByUser(t *testing.T) {
	tests := []struct {
		name   string
		idUser primitive.ObjectID
		idTask primitive.ObjectID
		isErr  bool
	}{
		{"test success", usersIds[0], tasksIds[0], false},
		{"test echec avec task introuvale", usersIds[0], primitive.NewObjectID(), true},
		{"test echec mongo", usersIds[0], tasksIds[0], true},
	}

	for _, tt := range tests {
		switch tt.name {
		case "test echec mongo":
			err := s.UpdateOneDonePropertyByUser(cancelCtx, tt.idUser, tt.idTask)
			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		
		default:
			err := s.UpdateOneDonePropertyByUser(ctx, tt.idUser, tt.idTask)
			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
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
		{"test echec mongo", usersIds[0], tasksIds[0], "label updated", true},
	}

	for _, tt := range tests {
		updateDto := TaskUpdateLabelDto{
			Label: tt.label,
		}
		
		switch tt.name {
		case "test echec mongo":
			err := s.UpdateOneLabelPropertyByUser(cancelCtx, tt.idUser, tt.idTask, updateDto)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		default:
			err := s.UpdateOneLabelPropertyByUser(ctx, tt.idUser, tt.idTask, updateDto)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
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
		{"test echec mongodb", usersIds[0], tasksIds[0], true},
	}

	for _, tt := range tests {
		if tt.name == "test echec mongodb" {
			err := s.DeleteOneByUser(cancelCtx, tt.idUser, tt.idTask)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		}else{
			err := s.DeleteOneByUser(ctx, tt.idUser, tt.idTask)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		}
	}
}

func TestDeleteManyByUser(t *testing.T) {
	setup := func(email string) func() (primitive.ObjectID, primitive.ObjectID) {
		return func() (primitive.ObjectID, primitive.ObjectID) {
			user := &auth.User{
				Email:    email,
				Password: "password",
			}

			result, err := userCollection.InsertOne(ctx, user)
			if err != nil {
				t.Fatalf("Erreur lors de la création de l'utilisateur : %v", err)
			}

			userID := result.InsertedID.(primitive.ObjectID)
			usersIds = append(usersIds, userID)

			task := &Task{
				Label:     "Test task",
				Done:      false,
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
				IdUser:    userID,
			}

			res, err := taskCollection.InsertOne(ctx, task)
			if err != nil {
				t.Fatalf("Erreur lors de la création de la tâche : %v", err)
			}

			taskID := res.InsertedID.(primitive.ObjectID)
			tasksIds = append(tasksIds, taskID)

			return userID, taskID
		}
	}

	tests := []struct {
		name  string
		setup func() (userID primitive.ObjectID, taskID primitive.ObjectID)
		isErr bool
	}{
		{
			name: "test success",
			setup: setup("taskTest3@gmail.com"),
			isErr: false,
		},
		{
			name: "test echec mongodb",
			setup: setup("taskTest4@gmail.com"),
			isErr: true,
		},
	}

	for _, tt := range tests {
		userID, taskID := tt.setup()

		if tt.name == "test echec mongodb" {
			err := s.DeleteManyByUser(cancelCtx, userID, []primitive.ObjectID{taskID})

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		}else{
			err := s.DeleteManyByUser(ctx, userID, []primitive.ObjectID{taskID})

			if (err != nil) != tt.isErr {
				t.Errorf("got error %v, expected error: %v", err, tt.isErr)
			}
		}
	}
}

func TestDeleteById(t *testing.T){
	setup := func(email string) func() (primitive.ObjectID, primitive.ObjectID) {
		return func() (primitive.ObjectID, primitive.ObjectID) {
			user := &auth.User{
				Email:    email,
				Password: "password",
			}

			result, err := userCollection.InsertOne(ctx, user)
			if err != nil {
				t.Fatalf("Erreur lors de la création de l'utilisateur : %v", err)
			}

			userID := result.InsertedID.(primitive.ObjectID)
			usersIds = append(usersIds, userID)

			task := &Task{
				Label:     "Test task",
				Done:      false,
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
				IdUser:    userID,
			}

			res, err := taskCollection.InsertOne(ctx, task)
			if err != nil {
				t.Fatalf("Erreur lors de la création de la tâche : %v", err)
			}

			taskID := res.InsertedID.(primitive.ObjectID)
			tasksIds = append(tasksIds, taskID)

			return userID, taskID
		}
	}

	tests := []struct {
		name string
		setup func() (userID primitive.ObjectID, taskID primitive.ObjectID)
		isErr bool
	}{
		{
			name: "test success",
			setup: setup("taskTest5@gmail.com"),
			isErr: false,
		},
		{
			name: "test echec mongodb",
			setup: setup("taskTest6@gmail.com"),
			isErr: true,
		},
	}

	for _, tt := range tests {
		_, taskID := tt.setup()

		if tt.name == "test echec mongodb" {
			err := s.DeleteById(cancelCtx, []primitive.ObjectID{taskID})

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		}else{
			err := s.DeleteById(ctx, []primitive.ObjectID{taskID})

			if (err != nil) != tt.isErr {
				t.Errorf("got error %v, expected error: %v", err, tt.isErr)
			}
		}
	}
}

func TestDeleteAllTasks(t *testing.T){
	tests := []struct {
		name string
		isErr bool
	}{
		{"test success", false},
		{"test echec mongodb", true},
	}

	for _, tt := range tests {

		if tt.name == "test echec mongodb" {
			err := s.DeleteAllTasks(cancelCtx)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		}else{
			err := s.DeleteAllTasks(ctx)

			if (err != nil) != tt.isErr {
				t.Errorf("got error %v, expected error: %v", err, tt.isErr)
			}
		}
	}
}
