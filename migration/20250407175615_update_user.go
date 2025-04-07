package migration

import (
	"context"
	"todof/mod/migratormongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var UpdateSchemaUserCOllection = migratormongodb.Migration{
	Name: "20250407175615_update_user",
	Up: func(db *mongo.Database) error {
		ctx := context.Background()

		// Nouveau schema avec les champs ajoutés
		validator := bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"email", "password"},
				"properties": bson.M{
					"email": bson.M{
						"bsonType":    "string",
						"description": "Email de l'utilisateur",
					},
					"password": bson.M{
						"bsonType":    "string",
						"description": "Mot de passe (hashé)",
					},
					"role": bson.M{
						"bsonType":    "string",
						"description": "Rôle de l'utilisateur (admin, user, etc.)",
					},
					"username": bson.M{
						"bsonType":    "string",
						"description": "Nom d'utilisateur",
					},
					"created_at": bson.M{
						"bsonType":    "date",
						"description": "Date de création du compte",
					},
					"updated_at": bson.M{
						"bsonType":    "date",
						"description": "Date de dernière modification du compte",
					},
				},
			},
		}

		// Met à jour la validation du schema de la collection
		command := bson.D{
			{Key: "collMod", Value: "users"},
			{Key: "validator", Value: validator},
			{Key: "validationLevel", Value: "moderate"},
		}

		if err := db.RunCommand(ctx, command).Err(); err != nil {
			return err
		}

		return nil
	},
}