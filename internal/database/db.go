package database

import (
	"fmt"
	"todo_formation/internal/logger"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func buildDSN() string {
	host := os.Getenv("DB_HOST") 
	port := os.Getenv("DB_PORT") 
	user := os.Getenv("DB_USER") 
	password := os.Getenv("DB_PASSWORD") 
	name := os.Getenv("DB_NAME") 
	zone := os.Getenv("TIMEZONE")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		host, user, password, name, port, zone,
	)
}

func Connect() {
	dsn := buildDSN()
	maxRetries := 10
	var err error

	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: true,
		})

		if err == nil {
			break
		}

		logger.Warnf("Tentative %d/%d - Connexion DB %v échouée : %v", i+1, maxRetries, os.Getenv("DB_NAME"), err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Fatalf("Connexion à la base de données %v échouée après %d tentatives : %v", os.Getenv("DB_NAME"), maxRetries, err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Fatalf("Erreur lors de la récupération de la pool : %v", err)
	}
 
	var dbName string
	err = sqlDB.QueryRow("SELECT current_database()").Scan(&dbName)
	if err != nil {
		logger.Errorf("Erreur lors de la récupération du nom de la base : %v", err)
		return
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Successf("Connexion à la base de données %v réussie", dbName)
}
