package auth

import (
	"context"
	"log"
	"os"
	"testing"
	"todof/internal/config"
	initializer "todof/internal/init"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var r UserRepoInterface
var s UserServiceInterface
var c *mongo.Collection
var ctx context.Context
var ids []primitive.ObjectID
var users []*User
var tokenString string

func TestMain(m *testing.M) {
	c = initializer.Db.Collection("users")
	r = NewUserRepo(initializer.Db)
	s = NewUserService(r, config.Get("JWT_SECRET"))
	ctx := context.Background()

	if _, err := c.DeleteMany(ctx, bson.M{}); err != nil {
		log.Fatalf("Erreur lors du nettoyage de la collection users : %v", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestRegister(t *testing.T){
	tests := []struct {
		name	 string
		email    string
		password string
		username string
		role     string
		isUser	 bool
		isErr    bool
	}{
		{"test avec user admin valid", "admin@gmail.com", "password", "admin", "admin", true, false},
		{"test avec email existant", "admin@gmail.com", "password", "admin2", "admin", false, true},
		{"test avec propriete role vide", "user@gmail.com", "password", "user", "", true, false},
	}

	for _, tt := range tests {
		userCreateDto := UserCreateDto{
			Email:    tt.email,
			Password: tt.password,
			Username: tt.username,
			Role:     tt.role,
		}

		user, err := s.Register(ctx, userCreateDto)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (user == nil) == tt.isUser{
			t.Errorf("%s: got user %v, expect user %v", tt.name, user, tt.isUser)
		}

		if user != nil {
			ids = append(ids, user.ID)
			users = append(users, user)
		}
	}
}

func TestLogin(t *testing.T){
	tests := []struct {
		name string
		email string
		password string
		isToken bool
		isErr bool
	}{
		{"test success", "admin@gmail.com", "password", true, false},
		{"test success", "fail@gmail.com", "password", false, true},
		{"test success", "admin@gmail.com", "invalid password", false, true},
	}

	for _, tt := range tests {
		userLoginDto := UserLoginDto{
			Email: tt.email,
			Password: tt.password,
		}

		token, err := s.Login(ctx, userLoginDto)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (err == nil) == tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (token == "") == tt.isToken{
			t.Errorf("%s: login fail with email %v and password %v", tt.name, tt.email, tt.password)
		}

		if token != ""{
			tokenString = token
		}
	}
}

func TestValidateToken(t *testing.T){
	tests := []struct {
		name string
		tokenString string
		isTokenClaims bool
		isErr bool
	}{
		{"test success", tokenString, true, false},
		{"test with invalid token", "invalidtokenstring", false, true},
	}

	for _, tt := range tests {
		token, err := s.ValidateToken(tt.tokenString)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (err == nil) == tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (token == nil) == tt.isTokenClaims {
			t.Errorf("%s: got %v, expect token claims", tt.name, token)
		}
	}
}

func TestGetProfilCurrentUser(t *testing.T){
	tests := []struct {
		name string
		id primitive.ObjectID
		isUser bool
		isErr bool
	}{
		{"test success", ids[0], true, false},
	}

	for _, tt := range tests {
		user, err := s.GetProfilCurrentUser(ctx, tt.id)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (err == nil) == tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (user == nil) == tt.isUser {
			t.Errorf("%s: got %v, expect user", tt.name, user)
		}
	}
}

func TestDeleteByAdmin(t *testing.T){
	tests := []struct {
		name string
		deletedCount int
		isErr bool
	}{
		{"test delete success", 2, false},
	}

	for _, tt := range tests {
		deletedCount, err := s.DeleteByAdmin(ctx, ids)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (err == nil) == tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (tt.deletedCount != deletedCount) != tt.isErr{
			t.Errorf("%s: got deletedCount : %v, expect deletedCount : %v", tt.name, deletedCount, tt.deletedCount)
		}
	}
}