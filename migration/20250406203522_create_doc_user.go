package migration

import (
	"context"
	"todof/mod/migratormongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CreateUsersCollection = migratormongodb.Migration{
	Name: "20250406203522_create_doc_user",
	Up: func(db *mongo.Database) error {
		ctx := context.Background()

		// schema validation
		validator := bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"email", "password"},
				"properties": bson.M{
					"email": bson.M{
						"bsonType": "string",
						"description": "Email de l'utilisateur",
					},
					"password": bson.M{
						"bsonType": "string",
						"description": "Mot de passe (hashé)",
					},
					"role": bson.M{
						"bsonType": "string",
						"description": "Rôle de l'utilisateur (admin, user, etc.)",
					},
				},
			},
		}

		opts := options.CreateCollection().SetValidator(validator)
		err := db.CreateCollection(ctx, "users", opts)
		if err != nil {
			return err
		}

		// index unique email
		indexModel := mongo.IndexModel{
			Keys: bson.M{"email": 1},
			Options: options.Index().SetUnique(true).SetName("idx_unique_email"),
		}

		_, err = db.Collection("users").Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return err
		}

		return nil
	},
}
