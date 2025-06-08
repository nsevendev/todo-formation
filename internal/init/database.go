package init

import (
	"context"
	"github.com/nsevenpack/env/env"
	"time"

	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database
var CDb *mongo.Client

func ConnexionDatabase(environement string) {
	logger.I("Connexion à la base de données ...")
	logger.If("Environement : %v", environement)

	dbName := env.Get("DB_NAME")
	if environement == "test" {
		dbName = dbName + "_" + environement
	}

	uri := env.Get("DB_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		logger.Ef("une erreur est survenue à la connexion à la base de données, uri : %v", uri)
		logger.Ff("Erreur de connexion à la base de données: %v", err)
	}

	Db = client.Database(dbName)
	CDb = client.Database(dbName).Client()

	res := CDb.Ping(ctx, nil)
	if res != nil {
		logger.Ff("Ping échoué sur la base '%s': %v", dbName, res.Error())
	}

	logger.If("URI de la base de données: %v", uri)
	logger.Sf("Connexion à la base de données %v réussie", dbName)
}
