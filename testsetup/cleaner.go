package testsetup

import (
	"context"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CleanCollections(ctx context.Context, db *mongo.Database, collections ...string) {
	for _, col := range collections {
		_, err := db.Collection(col).DeleteMany(ctx, bson.M{})
		if err != nil {
			logger.Ff("Erreur de nettoyage pour %s: %v", col, err)
		}
	}
}
