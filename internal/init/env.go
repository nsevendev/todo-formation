package init

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func initEnv() {
	env := "./.env"
	if os.Getenv("APP_ENV") == "test" {
		env = "../.env"
	}

	err := godotenv.Load(env)
	if err != nil {
		log.Fatalf("❌ [ERROR] Erreur de chargement du fichier .env : %v", err)
	}

	log.Printf("✅ [SUCCESS] Fichier .env chargé avec succès!")
}
