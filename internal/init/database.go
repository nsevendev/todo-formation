package init

import (
	"context"
	"os"
	"time"
	"todof/mod/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database
var CDb *mongo.Client

func ConnexionDatabase() {
	logger.Info("Connexion à la base de données ...")

	uri := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		logger.Errorf("une erreur est survenue à la connexion à la base de données, uri : %v", uri)
		logger.Fatalf("Erreur de connexion à la base de données: %v", err)
	}

	Db = client.Database(dbName)
	CDb = client.Database(dbName).Client()
	
	res := CDb.Ping(ctx, nil)
	if res != nil {
		logger.Fatalf("Ping échoué sur la base '%s': %v", dbName, res.Error())
	}

	logger.Infof("URI de la base de données: %v", uri)
	logger.Success("Connexion à la base de données réussie")
}