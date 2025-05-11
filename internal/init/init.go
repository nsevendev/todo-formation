package init

import (
	"todof/internal/config"
	"todof/migration"
	"todof/mod/migratormongodb"

	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

func init() {
	appEnv := config.Get("APP_ENV")
	logger.Init(appEnv)
	ConnexionDatabase(appEnv)
	ginresponse.SetFormatter(&ginresponse.JsonFormatter{})

	migrator := migratormongodb.New(Db)
	migrator.Add(migration.CreateUsersCollection)
	migrator.Add(migration.UpdateSchemaUserCOllection)
	migrator.Add(migration.CreateTasksCollection)
	// ajouter d'autres migrations ici ...
	if err := migrator.Apply(); err != nil {
		logger.Ff("Erreur lors de l'application des migrations : %v", err)
	}
}
