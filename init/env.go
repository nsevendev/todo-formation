package init

import (
	"log"

	"github.com/joho/godotenv"
)

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ [ERROR] Erreur de chargement du fichier .env : %v", err)
	}

	log.Printf("✅ [SUCCESS] Fichier .env chargé avec succès!")
}