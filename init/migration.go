package init

import (
	"database/sql"
	"os"
	"todof/internal/logger"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func initMigration() error {
	db, err := sql.Open("postgres", os.Getenv("NSC_MIGRATION_DB_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	logger.Info("Exécution des migrations...")

	_, errGoose := goose.CollectMigrations("./migrations", 0, goose.MaxVersion)
	logger.Infof("%v", errGoose)
	if errGoose != nil && strings.Contains(errGoose.Error(), "no migration files found") {
		logger.Infof("Error => %v", errGoose)
		logger.Warn("Aucune migration trouvée. Rien à faire.")
		return nil
	}

	if errGoose != nil {
		logger.Fatalf("Erreur en collectant les migrations: %v", errGoose)
	}
	
	if err := goose.Up(db, "./migrations"); err != nil {
		logger.Fatalf("Impossible d'exécuter les migrations: %v", err)
	}

	logger.Success("Migrations terminées.")
	
	return nil
}