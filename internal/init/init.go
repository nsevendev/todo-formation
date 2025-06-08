package init

import (
	"todof/internal/config"
	"todof/internal/job"
	"todof/migration"
	"todof/mod/migratormongodb"

	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

func init() {
	// START GET .ENV
	appEnv := config.Get("APP_ENV")

	// REDIS
	job.Redis(config.Get("REDIS_ADDR"))
	job.StartWorker()

	// LOGGER
	logger.Init(appEnv)

	// DB
	ConnexionDatabase(appEnv)

	// GIN RESPONSE FORMAT
	ginresponse.SetFormatter(&ginresponse.JsonFormatter{})

	// MIGRATION
	migrator := migratormongodb.New(Db)
	migrator.Add(migration.CreateUsersCollection)
	migrator.Add(migration.UpdateSchemaUserCOllection)
	migrator.Add(migration.CreateTasksCollection)
	// ajouter d'autres migrations ici ...
	if err := migrator.Apply(); err != nil {
		logger.Ff("Erreur lors de l'application des migrations : %v", err)
	}
}
