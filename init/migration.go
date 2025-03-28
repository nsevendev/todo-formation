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
	
	if err := goose.Up(db, "./migrations"); err != nil {
		logger.Fatalf("Impossible d'exécuter les migrations: %v", err)
	}

	logger.Success("Migrations terminées.")
	
	return nil
}