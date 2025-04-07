# Logger personnalisé

Ce logger ajoute des messages colorés, lisibles, avec le nom du fichier, la fonction, et la ligne d'où il est appelé.  
Il écrit les logs à la fois dans le terminal et dans un fichier par environnement.

---

## ✅ Initialisation

```go
import "app/logger"

// Lit APP_ENV ou utilise "dev" par défaut
logger.Init() 
// ferme proprement le fichier de log si le programme se termine
defer logger.Close()
```

Les fichiers de logs sont créés dans `tmp/log/<env>/log-YYYY-MM-DD.log`.

---

## 📋 Fonctions disponibles

- `logger.Success(msg)` / `Successf(...)`
- `logger.Info(msg)` / `Infof(...)`
- `logger.Warn(msg)` / `Warnf(...)`
- `logger.Error(msg)` / `Errorf(...)`
- `logger.Fatal(msg)` / `Fatalf(...)` → termine le programme

Chaque log affiche :
- Un emoji
- Le niveau (`INFO`, `WARN`, ...)
- Le fichier, la fonction et la ligne d'appel

---

## 🧹 Fermeture

```go
// Ferme proprement le fichier log.
// A appeller apres le logger.Init()
defer logger.Close()
```