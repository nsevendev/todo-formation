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
var ctx context.Context

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
		{"test avec user email existant", "test@gmail.com", "password", "test2", "user", true},
	}

	for _, tt := range tests {
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

func TestFindByEmail(t *testing.T){
	tests := []struct {
		name  string
		email string
		isUser bool
		isErr bool
	}{
		{"test avec un email existant", "test@gmail.com", true, false},
		{"test avec un email inexistant", "unknow@gmail.com", false, false},
	}

	for _, tt := range tests {
		user, err := r.FindByEmail(ctx, tt.email)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (user == nil) == tt.isUser {
			t.Errorf("%s: aucun utilisateur trouvé avec l'email : %v", tt.name, tt.email)
		}

		if user != nil && user.Email != tt.email{
			t.Errorf("%s: got email %v, expect email %v", tt.name, user.Email, tt.email)
		}
	}
}