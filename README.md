# Todo Formation  

- petit api de gestion de todo

## Prérequis

- docker, docker compose
- cli `ns` pour la gestion des containers
- `traefik nseven` qui tourne 

## Installation/Demarrage des containers dev

- copier coller `.env.dist` en `.env`, renseigner les variables d'environnement
- lancer la commande `ns c dev`à l'aide du CLI `ns` l'application est démarrée
ou utiliser les commandes docker traditionnelles (mais n'oubliez pas de specifier le env file)

## Arreter les containers

- utiliser le CLI `ns` pour arreter les containers, commande `ns d <namecontainer>...`  
ou tous autre commande docker

## Commande disponible  

- avec le CLI `ns` vous pouvez executer des commandes specifique au projet  
pour afficher cette liste tapez `ns c list`

## Installation/Demarrage des containers prod