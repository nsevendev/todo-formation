package init

import (
	"todof/logger"
	"todof/migration"
)

func init() {
	initEnv()
	logger.Init()
	defer logger.Close()
	ConnexionDatabase()

	migrator := migration.NewMigrator(Db)
	migrator.Add(migration.CreateUsersCollection)

	if err := migrator.ApplyMigrations(); err != nil {
		logger.Fatalf("Erreur lors de l'application des migrations : %v", err)
	}
}