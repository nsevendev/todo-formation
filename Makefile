# command
DOCKER = docker
DOCKER_COMP = docker compose

# name container
NAME_CONT_API = todo-f-api
NAME_CONT_FRONT = todo-f-front
NAME_CONT_DB = todo-f-db

# Misc
.DEFAULT_GOAL = help
.PHONY        : help build dev logs sh down dotadd dotef update-test-db

## —— 🎵 🐳 Docker Makefile Todo-formation 🐳 🎵 —————————————————————————————————————————
help: ## Affiche cette aide
	@grep -E '(^[a-zA-Z0-9\./_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

## —— 🎵 🐳 Containers dev 🐳 🎵 —————————————————————————————————————————

dev: ## Lancement des containers en mode dev
	@echo "🚀 Start container dev ---> START"
	$(DOCKER_COMP) up -d
	@echo "✅ Start container dev ---> END OK"

down: ## Arrêt des containers
	@echo "🚀 Close container dev ---> START"
	$(DOCKER_COMP) down
	@echo "✅ Close container dev ---> END OK"

logs: ## Affiche les logs du container specifié c="api", c="front" or c="db" (default is api)
	@$(eval c ?= api)
	@echo "🚀 Affichage des logs du container $(c) ---> START"
	@$(DOCKER) logs -f $(if $(filter $(c),api),$(NAME_CONT_API),$(if $(filter $(c),front),$(NAME_CONT_FRONT),$(if $(filter $(c),db),$(NAME_CONT_DB),$(error "Valeur de c invalide : $(c)"))))

sh: ## Ouvre un shell dans le container specifié c="api", c="front" or c="db" (default is api)
	@$(eval c ?= api)
	@echo "🚀 Ouverture d'un shell dans le container $(c) ---> START"
	@$(DOCKER) exec -it $(if $(filter $(c),api),$(NAME_CONT_API),$(if $(filter $(c),front),$(NAME_CONT_FRONT),$(if $(filter $(c),db),$(NAME_CONT_DB),$(error "Valeur de c invalide : $(c)")))) bash

rebuild: ## Reconstruit uniquement l'image d'un service spécifique (exemple : c="api")
	@if [ -z "$(c)" ]; then \
		echo "❌ Spécifie le service avec c=api, front ou db (exemple : make rebuild c=api)"; \
		exit 1; \
	fi
	@echo "🔨 Reconstruction de l'image du service $(c) ---> START"
	$(DOCKER_COMP) build $(if $(filter $(c),api),$(NAME_CONT_API),$(if $(filter $(c),front),$(NAME_CONT_FRONT),$(if $(filter $(c),db),$(NAME_CONT_DB),$(error "Valeur de c invalide : $(c)"))))
	@echo "✅ Reconstruction de l'image $(c) ---> END OK"