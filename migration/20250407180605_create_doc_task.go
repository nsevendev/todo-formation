package migration

import (
	"context"
	"todof/mod/migratormongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CreateTasksCollection = migratormongodb.Migration{
	Name: "20250407180605_create_doc_task",
	Up: func(db *mongo.Database) error {
		ctx := context.Background()

		validator := bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"label", "done", "id_user"},
				"properties": bson.M{
					"label": bson.M{
						"bsonType":    "string",
						"description": "Libellé de la tâche",
					},
					"done": bson.M{
						"bsonType":    "bool",
						"description": "État de la tâche (complétée ou non)",
					},
					"created_at": bson.M{
						"bsonType":    "date",
						"description": "Date de création",
					},
					"updated_at": bson.M{
						"bsonType":    "date",
						"description": "Date de mise à jour",
					},
					"id_user": bson.M{
						"bsonType":    "objectId",
						"description": "ID de l'utilisateur propriétaire",
					},
				},
			},
		}

		opts := options.CreateCollection().SetValidator(validator)
		if err := db.CreateCollection(ctx, "tasks", opts); err != nil {
			return err
		}

		// index id_user
		indexModel := mongo.IndexModel{
			Keys: bson.M{"id_user": 1},
			Options: options.Index().SetName("idx_id_user"),
		}
		if _, err := db.Collection("tasks").Indexes().CreateOne(ctx, indexModel); err != nil {
			return err
		}
		
		return nil
	},
}