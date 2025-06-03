package testsetup

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CleanCollections(ctx context.Context, db *mongo.Database, collections ...string) {
	for _, col := range collections {
		_, err := db.Collection(col).DeleteMany(ctx, bson.M{})
		if err != nil {
			fmt.Printf("Failed to clean collection '%s': %v\n", col, err)
		} else {
			fmt.Printf("Collection '%s' has been cleaned successfully.\n", col)
		}
	}
}
