package testsetup

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func CleanCollections(ctx context.Context, db *mongo.Database, collections ...string) {
	for _, col := range collections {
		_, err := db.Collection(col).DeleteMany(ctx, bson.M{})
		if err != nil {
			fmt.Printf("🔴 Failed to clean collection '%s': %v\n", col, err)
		}
	}

	fmt.Printf("Collection cleaning completed for: %v\n", collections)
}

func LogI(t *testing.T, name string) {
	t.Logf("🔵 Test: %s", name)
}
