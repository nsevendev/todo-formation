package auth

import (
	"context"
	"errors"

	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	collection *mongo.Collection
}

type userRepoInterface interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	Create(ctx context.Context, user *User) error
}

func NewUserRepo(db *mongo.Database) userRepoInterface {
	return &userRepo{
		collection: db.Collection("users"),
	}
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Wf("Erreur mongo no document: %v", err)
			return nil, nil
		}
		logger.Ef("Erreur à la recuperation du user: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Wf("Erreur mongo no document: %v", err)
			return nil, nil
		}
		logger.Ef("Erreur à la recuperation du user: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(ctx context.Context, user *User) error {
	user.SetTimeStamps()
	
	existingUser, err := r.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		logger.Ef("L'utilisateur avec cet email existe déjà : %s", user.Email)
		return errors.New("impossible de créer votre compte")
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		logger.Ef("Une erreur est survenue au moment de creer le compte : %v", err)
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}
