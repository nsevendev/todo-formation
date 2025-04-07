# Réponses JSON standardisées  

Utilise `apiresponse` pour unifier les réponses Json standardisées

## Prérequis

- github.com/gin-gonic/gin

## 📦 Fichier `apiresponse/apiresponse.go`

- Fournit des fonctions pour envoyer des réponses JSON standardisées.

## ✅ Exemple d'utilisation dans un handler

```go
import "app/apiresponse"

func GetUser(c *gin.Context) {
	user := findUserInDB()
	if user == nil {
		api.Error(c, http.StatusNotFound, "Utilisateur non trouvé", nil)
		return
	}
	api.Success(c, user, "Utilisateur trouvé")
}
```

