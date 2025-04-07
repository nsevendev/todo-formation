package init

import (
	"todof/migration"
	"todof/mod/logger"
	"todof/mod/migratormongodb"
)

func init() {
	initEnv()
	logger.Init()
	defer logger.Close()
	ConnexionDatabase()

	migrator := migratormongodb.New(Db)
	migrator.Add(migration.CreateUsersCollection)
	// ajouter d'autres migrations ici ...
	if err := migrator.Apply(); err != nil {
		logger.Fatalf("Erreur lors de l'application des migrations : %v", err)
	}
}