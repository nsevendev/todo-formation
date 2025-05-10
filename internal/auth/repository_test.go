package auth

import (
	"context"
	"log"
	"os"
	"testing"
	initializer "todof/internal/init"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var r UserRepoInterface
var c *mongo.Collection
var ctx context.Context
var ids []primitive.ObjectID

func TestMain(m *testing.M) {
	c = initializer.Db.Collection("users")
	r = NewUserRepo(initializer.Db)
	ctx := context.Background()

	code := m.Run()

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
		{"test avec user admin valid", "admin@gmail.com", "password", "admin", "admin", false},
		{"test avec user avec propriété role vide", "test@gmail.com", "password", "test", "", false},
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

		if user != nil {
			ids = append(ids, user.ID)
		}
	}
}

func TestFindByID(t *testing.T){
	tests := []struct {
		name string
		id primitive.ObjectID
		isUser bool
		isErr bool
	}{
		{"test avec l'id d'un utilisateur existant", ids[0], true, false},
		{"test avec un id inexistant", primitive.NewObjectID(), false, false},
	}

	for _, tt := range tests {
		user, err := r.FindByID(ctx, tt.id)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (user == nil) == tt.isUser {
			t.Errorf("%s: aucun utilisateur trouvé avec l'id : %v", tt.name, tt.id)
		}
	}
}

func TestFindNonAdmin(t *testing.T){
	tests := []struct {
		name string
		isUsers bool
		isErr bool
	}{
		{"test avec un utilisateur non admin trouvé", true, false},
	}

	for _, tt := range tests {
		users, err := r.FindNonAdmin(ctx)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (users == nil) == tt.isUsers{
			t.Errorf("%s: aucun utilisateurs trouvé", tt.name)
		}

		if tt.isUsers && users != nil && len(users) != 1 {
			t.Errorf("%s: got users length : %d, expect users length : 1", tt.name, len(users))
		}
	}
}

func TestDelete(t *testing.T){
	tests := []struct {
		name string
		id primitive.ObjectID
		deletedCount int64
		isErr bool
	}{
		{"test succes du delete", ids[0], 1, false},
	}

	for _, tt := range tests {
		deletedCount, err := r.Delete(ctx, tt.id)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if tt.deletedCount != deletedCount {
			t.Errorf("%s: got deletedCount %v, expect deletedCount %v", tt.name, deletedCount, tt.deletedCount)
		}
	}
}

func TestDeleteMany(t *testing.T){
	tests := []struct {
		name string
		filter interface{}
		deletedCount int64
		isErr bool
	}{
		{"test avec filter valid avec aucun document correspondant", bson.M{"role": "test"}, 0, false},
		{"test avec filter valid et 1 document correspondant", bson.M{}, 1, false},
	}

	for _, tt := range tests {
		deletedCount, err := r.DeleteMany(ctx, tt.filter)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if tt.deletedCount != deletedCount {
			t.Errorf("%s: got deletedCount %v, expect deletedCount %v", tt.name, deletedCount, tt.deletedCount)
		}
	}
}