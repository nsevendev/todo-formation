package task

import (
	"context"
	"os"
	"testing"
	"todof/internal/auth"
	initializer "todof/internal/init"

	"go.mongodb.org/mongo-driver/mongo"
)

var r taskRepoInterface
var userRepo auth.UserRepoInterface
var s TaskServiceInterface
var c *mongo.Collection
var ctx context.Context

//service
func TestMain(m *testing.M) {
	c = initializer.Db.Collection("tasks")
	r = NewTaskRepo(initializer.Db)
	userRepo = auth.NewUserRepo(initializer.Db)
	s = NewTaskService(r, userRepo)
	ctx := context.Background()

	code := m.Run()

	os.Exit(code)
}