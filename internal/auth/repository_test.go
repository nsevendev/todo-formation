package auth

import (
	"context"
	"testing"
	initializer "todof/internal/init"
)

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