package auth

import (
	"context"
	"log"
	"os"
	"testing"
	initializer "todof/internal/init"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var r UserRepoInterface
var c *mongo.Collection

func TestMain(m *testing.M) {
	c = initializer.Db.Collection("users")
	r = NewUserRepo(initializer.Db)

	code := m.Run()

	ctx := context.Background()
	if _, err := c.DeleteMany(ctx, bson.M{}); err != nil {
		log.Fatalf("Erreur lors du nettoyage de la collection users : %v", err)
	}

	os.Exit(code)
}

func TestCreate(t *testing.T){
	tests := []struct {
		name	 string
		email    string
        password string
        username string
        role     string
		isErr    bool
	}{
		{"test avec user valid", "test@gmail.com", "password", "test", "user", false},
		{"test avec user email existant", "test2@gmail.com", "password", "test2", "user", true},
	}

	for _, tt := range tests {
		ctx := context.Background()
		r := NewUserRepo(initializer.Db)

		user := &User{
			Email:    tt.email,
			Password: tt.password,
			Username: tt.username,
			Role:     tt.role,
		}

		err := r.Create(ctx, user)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}
	}
}