# Lancement de l'environnement

## Prérequis

- `docker`, `docker compose`
- `make`
- `traefik nseven` qui tourne
- copier coller `.env.dist` en `.env`, renseigner les variables d'environnement
- définissez la variable `APP_ENV` avec `dev`, `prod`, `preprod` en fonction de l'environnement souhaité  
  (par defaut `dev`)

## Demarrage/Arrêt des containers

- Lancer la commande `make up` pour demarrer les containers
- Arrêter les containers avec `make down`  
