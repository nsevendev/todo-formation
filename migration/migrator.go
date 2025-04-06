package migration

import (
	"context"
	"fmt"
	"time"
	"todof/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migration struct {
	Name string
	Up   func(db *mongo.Database) error
}

type Migrator struct {
	db         *mongo.Database
	migrations []Migration
}

func NewMigrator(db *mongo.Database) *Migrator {
	return &Migrator{
		db:         db,
		migrations: []Migration{},
	}
}

func (m *Migrator) Add(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) ApplyMigrations() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	applied := make(map[string]bool)
	cursor, err := m.db.Collection("migrations").Find(ctx, bson.M{})
	if err != nil {
		logger.Errorf("impossible de lire la collection des migrations : %v", err)
		return fmt.Errorf("impossible de lire la collection des migrations : %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err == nil {
			logger.Infof("Ajout de la migration déjà appliquée: %s", result["name"].(string))
			applied[result["name"].(string)] = true
		}
	}

	for _, migration := range m.migrations {
		if applied[migration.Name] {
			logger.Infof("Migration déjà appliquée: %s", migration.Name)
			continue
		}

		fmt.Println("📦 Migration:", migration.Name)
		if err := migration.Up(m.db); err != nil {
			logger.Errorf("échec de la migration %s : %v", migration.Name, err)
			return fmt.Errorf("échec de la migration %s : %w", migration.Name, err)
		}

		_, err := m.db.Collection("migrations").InsertOne(ctx, bson.M{
			"name":      migration.Name,
			"createdAt": time.Now(),
		})
		if err != nil {
			logger.Errorf("impossible d'enregistrer la migration %s : %v", migration.Name, err)
			return fmt.Errorf("impossible d'enregistrer la migration %s : %w", migration.Name, err)
		}
		logger.Infof("Migration appliquée: %s", migration.Name)
	}

	return nil
}