# Todo Formation  

> Application de todo pour des exercices de formation

## Prérequis

- Attention il y a deux Makefile, un à la racine du projet et un dans le dossier `api/`. Les commandes sont à lancer depuis la racine de chaque dossier respectif.  
ils ont chacun leur propre utilité, le Makefile à la racine du projet est pour les commandes globales et le Makefile dans le dossier `api/` est pour les commandes spécifiques à l'api. (python ou django)  

- Tous les makefiles ont un commande help pour afficher les commandes disponibles.  
```bash
$ make help
# ou tout simplement
$ make
```
**Donc executer ces commandes pour vous afficher l'aide des commandes disponible**  

- ATTENTION: le repository ìnfo-traefik doit être lancé, avant de lancer ce projet  

## Demarrer l'application en mode dev

- Demarrer les containers
```bash
$ make dev
```

- Arrêter les containers
```bash
$ make down
```

- Une fois les containers démarrés, vous pouvez accéder à l'application sur l'url suivante:
`todo-cda.api.local` pour l'api et `todo-cda.front.local` pour le front  
- La bdd postgresql à une liaison interne au container api via le réseau docker todo-f

