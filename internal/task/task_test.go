package task

import (
	"context"
	"log"
	"os"
	"testing"
	"todof/internal/auth"
	initializer "todof/internal/init"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var s TaskServiceInterface
var userCollection *mongo.Collection
var ctx context.Context

//service
func TestMain(m *testing.M) {
	taskCollection := initializer.Db.Collection("tasks")
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
				}

				result, err := userCollection.InsertOne(ctx, user)
				if err != nil {
					t.Fatalf("Erreur lors de la création du user: %v", err)
				}

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

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if err == nil && task.Label != tt.label {
			t.Errorf("%s: got label %s, expect label %s", tt.name, task.Label, tt.label)
		}
	}
}