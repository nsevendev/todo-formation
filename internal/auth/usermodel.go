package auth

import (
	"time"

	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
	Role     string             `bson:"role" json:"role"`
	Username string             `bson:"username" json:"username"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type UserCreateDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"omitempty"`
}

type UserDeleteDto struct {
	Ids []primitive.ObjectID `json:"ids" binding:"required"`
}

type UserLoginDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Ef("Erreur de hashage du mot de passe: %v", err)
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) SetTimeStamps() {
	now := primitive.NewDateTimeFromTime(time.Now())
	if u.CreatedAt == 0 {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}