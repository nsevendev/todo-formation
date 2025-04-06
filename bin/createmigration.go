package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ Utilisation : go run bin/migration-create.go <nom_migration>")
		os.Exit(1)
	}
 
	name := strings.Join(os.Args[1:], "_")
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("migration/%s_%s.go", timestamp, name)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("❌ Erreur lors de la création du fichier : %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Printf("✅ Fichier de migration créé : %s\n", filename)
}