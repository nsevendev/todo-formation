-include .env

# Redefinir MAKEFILE_LIST pour qu'il ne contienne que le Makefile
MAKEFILE_LIST := Makefile

ifeq ($(APP_ENV),dev)
  CONTAINER_NAME := todof-go
  CONTAINER_NAME_DB := todof-db
  COMPOSE_FILES := -f docker/compose.yaml -f docker/compose.override.yaml
else ifeq ($(APP_ENV),preprod)
  CONTAINER_NAME := todof-go-preprod
  CONTAINER_NAME_DB := todof-db-preprod
  COMPOSE_FILES := -f docker/compose.preprod.yaml
else ifeq ($(APP_ENV),prod)
  CONTAINER_NAME := todof-go-prod
  CONTAINER_NAME_DB := todof-db-prod
  COMPOSE_FILES := -f docker/compose.prod.yaml
endif

# Variables
GO_COMMAND_CONTAINER_TEST := docker exec -i -e APP_ENV=test $(CONTAINER_NAME) go
GO_COMMAND_CONTAINER := docker exec -i $(CONTAINER_NAME) go
SWAG_COMMAND_CONTAINER := docker exec -i $(CONTAINER_NAME) swag
BASH_CONTAINER := docker exec -it $(CONTAINER_NAME) sh
BASH_CONTAINER_DB := docker exec -it $(CONTAINER_NAME_DB) sh

.PHONY: help cm dev devb devbnod ddev tidy addget
.DEFAULT_GOAL = help

## —— 🐳 ALL 🐳 ——————————————————————————————————
help: ## Afficher l'aide
	@grep -E '(^[a-zA-Z0-9\./_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

starter: ## Instruction pour installer le projet
	cat doc/lancer-environement.md

## —— 🐳 CONTAINER 🐳 ——————————————————————————————————

## Attention définisser l'environement avec APP_ENV=dev, APP_ENV=prod, APP_ENV=preprod
## dans le .env

up: ## Demarre l'environnement
	docker compose --env-file .env $(COMPOSE_FILES) up -d

upb: ## Demarre l'environnement avec build
	docker compose --env-file .env $(COMPOSE_FILES) up -d --build

upbnod: ## Demarre l'environnement sans mode detache et avec build
	docker compose --env-file .env $(COMPOSE_FILES) up --build

down: ## Arrête les conteneurs
	docker compose --env-file .env $(COMPOSE_FILES) down

## —— 🐳 TOOl 🐳 ——————————————————————————————————

cm: ## Crée un fichier pour la migration - usage: make cm file=nom_du_fichier
	$(GO_COMMAND_CONTAINER) run mod/migratormongodb/bin/createfilemigration.go $(file)

tidy: ## Execute go mod tidy pour nettoyer les dependances
	$(GO_COMMAND_CONTAINER) mod tidy

gg: ## Ajoute une dependance - usage: make gg dep=path_de_la_dependance
	$(GO_COMMAND_CONTAINER) get $(dep)

s: ## Ouvre un shell dans le conteneur app
	$(BASH_CONTAINER)

sdb: ## Ouvre un shell dans le conteneur database
	$(BASH_CONTAINER_DB)

l: ## Affiche les logs du conteneur app
	docker logs -f $(CONTAINER_NAME)

ldb: ## Affiche les logs du conteneur database
	docker logs -f $(CONTAINER_NAME_DB)

swag: ## Genere la doc swagger
	$(SWAG_COMMAND_CONTAINER) init -o doc -g main.go app/controller internal doc

t: ## Execute les tests
	$(GO_COMMAND_CONTAINER_TEST) test -v -cover ./...