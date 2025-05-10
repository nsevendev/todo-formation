package auth

import "testing"

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
	}
}

