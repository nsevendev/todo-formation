-include .env

# Redefinir MAKEFILE_LIST pour qu'il ne contienne que le Makefile
MAKEFILE_LIST := Makefile

CONTAINER_APP := todof_$(APP_ENV)
CONTAINER_DB := todof_db_$(APP_ENV)
CONTAINER_REDIS := todof_redis_$(APP_ENV)
COMPOSE_FILES := -f docker/compose.$(APP_ENV).yaml
ENV_FILE := --env-file .env

# Variables
GO_COMMAND_CONTAINER_TEST := docker exec -i -e APP_ENV=test $(CONTAINER_APP) go test

.PHONY: help cm dev devb devbnod ddev tidy addget
.DEFAULT_GOAL = help

## —— HELP ——————————————————————————————————
help: ## Afficher l'aide
	@grep -E '(^[a-zA-Z0-9\./_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

starter: ## Instruction pour installer le projet
	cat doc/lancer-environement.md

## —— CONTAINER ——————————————————————————————————
##
## Attention définisser l'environement avec
## APP_ENV=dev, APP_ENV=prod, APP_ENV=preprod
## dans le .env
##
up: ## Demarre l'environnement
	docker compose $(ENV_FILE) $(COMPOSE_FILES) up ${APP_ENV} db redis -d

upb: ## Demarre l'environnement avec build
	docker compose $(ENV_FILE) $(COMPOSE_FILES) up ${APP_ENV} db redis -d --build

down: ## Arrête les conteneurs
	docker compose $(ENV_FILE) $(COMPOSE_FILES) down ${APP_ENV} db redis

## —— GO ——————————————————————————————————
g: ## Execute une commande go dans le conteneur app - usage: make go c=commande_go
	docker exec -it $(CONTAINER_APP) go $(c)

gg: ## Ajoute une dependance - usage: make gg dep=path_de_la_dependance
	$(MAKE) g c="get $(dep)"

tidy: ## Execute go mod tidy pour nettoyer les dependances
	$(MAKE) g c="mod tidy"

s: ## Ouvre un shell dans le conteneur app
	docker exec -it $(CONTAINER_APP) sh

l: ## Affiche les logs du conteneur app
	docker logs -f $(CONTAINER_APP)

## —— MIGRATION ——————————————————————————————————
cm: ## Crée un fichier pour la migration - usage: make cm file=nom_du_fichier
	@$(MAKE) g c="run mod/migratormongodb/bin/createfilemigration.go $(file)"

## —— TEST ——————————————————————————————————
t: ## Execute tous les tests
	$(GO_COMMAND_CONTAINER_TEST) ./...

tf: ## Execute les tests d'un path - usage: make tf file=path_relatif_du_fichier
	$(GO_COMMAND_CONTAINER_TEST) $(file)

tv: ## Execute tous les tests avec verbose
	$(GO_COMMAND_CONTAINER_TEST) -v -cover ./...

tvf: ## Execute les tests d'un path - usage: make tf file=path_relatif_du_fichier
	$(GO_COMMAND_CONTAINER_TEST) -v -cover $(file)

tcf: ## Execute les tests cover par fichier - usage: make tcf file=path_du_fichier
	$(GO_COMMAND_CONTAINER_TEST) -coverprofile=coverage/coverage.out ./$(file)
	docker exec -it $(CONTAINER_APP) go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	open coverage/coverage.html

## —— DATABASE ——————————————————————————————————
sdb: ## Ouvre un shell dans le conteneur database
	docker exec -it $(CONTAINER_DB) bash

ldb: ## Affiche les logs du conteneur database
	docker logs -f $(CONTAINER_DB)

## —— REDIS ——————————————————————————————————
sredis: ## Ouvre un shell dans le conteneur redis
	docker exec -it $(CONTAINER_REDIS) bash

lredis: ## Affiche les logs du conteneur redis
	docker logs -f $(CONTAINER_REDIS)

## —— SWAGGER ——————————————————————————————————
swag: ## Execute une commande swagger dans le conteneur app - usage: make swag c=commande_go
	docker exec -it $(CONTAINER_APP) swag $(c)