package init

import (
	"todof/internal/database"
	"todof/internal/logger"
)

func init() {
	initEnv()
	logger.Init()
	defer logger.Close()
	database.InitConnect()
	initMigration()
}