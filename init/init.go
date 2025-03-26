package init

import (
	"todo_formation/internal/database"
	"todo_formation/internal/logger"
)

func init() {
	initEnv()
	logger.InitFromEnv()
	defer logger.Close()
	database.Connect()
	database.AutoMigrate()
}