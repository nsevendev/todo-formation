# 📚 Documentation - MigratorMongoDB (MongoDB Custom Migration System)

## Prérequis

- go.mongodb.org/mongo-driver
- utilise le package `logger` de nseven

## 🔧 Objectif
Un système de migrations simple pour MongoDB, écrit en Go, permettant de :
- Définir des fichiers de migration versionnés
- Appliquer chaque migration une seule fois
- Ajouter des validations et des index

---

## 📁 Structure des dossiers

```
/migratormongodb        → Contient le cœur du système de migration (Migrator)
/migratormongodb/bin    → Script pour générer un fichier de migration vide
/migration              → Contient toutes tes migrations .go (propre au projet mais dois exister)
```

---

## 🚀 1. Générer un fichier de migration

Utilise la commande suivante :

```bash
# go run <path_to_migrator> <name_file>
go run migratormongodb/bin/createfilemigration.go create_users_collection
```

Cela va créer un fichier du type :

```
migration/20250406203522_create_users_collection.go
```

---

## ✍️ 2. Écrire une migration

Dans le fichier généré, importe le package `migratormongodb` et définis une `Migration` :

```go
package migration

import (
	"context"
	"app/migratormongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CreateUsersCollection = migratormongodb.Migration{
    // par exemple le nom du fichier (pas obligatoire)
	Name: "20250406203522_create_users_collection",
	Up: func(db *mongo.Database) error {
		// ... votre code de migration ici
	},
}
```

---

## 🧠 3. Enregistrer les migrations dans l'app

Dans ton `init.go` (ou autre fichier d'initialisation) :

```go
import (
	"app/migratormongodb"
	"app/migration"
)

func init() {
    // creation du helper migrator
	migrator := migratormongodb.New(Db)

    // ajouter une migration
	migrator.Add(migration.CreateUsersCollection)
    // vous pouvez ajouter plusieurs migrations ici ....

    // application des migrations
	if err := migrator.Apply(); err != nil {
		fmt.Fatalf("Erreur lors des migrations : %v", err)
	}
}
```

> ✅ Chaque migration ne sera exécutée **qu'une seule fois**.  
> Leur nom est enregistré dans la collection MongoDB `migrations`.  
> de votre base de données.

---

## 🔄 4. Ajouter une nouvelle migration

À chaque changement (ex : nouveau champ, nouvel index) :

1. Crée un nouveau fichier :
   ```bash
   go run migratormongodb/bin/createfilemigration.go add_username_field
   ```
2. Écris ton code de migration dedans
3. Ajoute-la avec `.Add()` dans l’init
4. Redémarre l’application pour appliquer la migration

---

## 📌 Bonnes pratiques

- Le nom de la migration **doit correspondre exactement** au nom du fichier (sans `.go`)
- Regroupe toutes les migrations dans un seul endroit (`migration/`)
- Tu peux chaîner plusieurs `.Add(...)` pour en exécuter plusieurs

---