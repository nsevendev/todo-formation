# Middleware  

## Prérequis

- gin
- apiresponse (package nseven)
- logger (package nseven)

## Liste des middlewares

- `routenotfound.go` : Middleware pour gérer les erreurs 404  
Cela permet d'avoir une gestion centralisée et cohérente des routes non trouvées  
avec une réponse JSON propre et un log structuré.  
```go
s := gin.Default()
s.Use(middleware.RouteNotFound())
```
