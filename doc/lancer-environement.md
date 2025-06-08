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
**Pour afficher toutes les commandes "make" tapez "make" dans le terminal**

## Acceder au swagger

- Lancer les container
- Aller sur https://todof.local/swagger/index.html  
  (ATTENTION LE HOST CHANGE EN FONCTION DE L'ENVIRONNEMENT)
  (todof.local EST LE HOST PAR DEFAUT EN DEV)

## Lancement des tests

- creer un fichier `.env.test` à la racine du projet
- ajouter la variable `APP_ENV=test` dans ce fichier
- lancer la commande `make t` pour lancer les tests ou tous autres commande make de tests
