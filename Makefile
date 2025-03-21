# Définition des couleurs
GREEN := $(shell echo "\033[32m")
YELLOW := $(shell echo "\033[33m")
CYAN := $(shell echo "\033[36m")
RESET := $(shell echo "\033[0m")

POSTGRES_USER := todo_formation
POSTGRES_DB := todo_formation

APP := todo-formation-go
APP_DB := todo-formation-db
D := docker
DC := docker compose

help:
	@echo "$(CYAN)Commandes disponibles :$(RESET)"
	@echo "$(CYAN) Container Dev --------------------------- $(RESET)"
	@echo "$(GREEN)  build  $(RESET): $(YELLOW)Construit l'image Docker de l'application.$(RESET)"
	@echo "$(GREEN)  dev    $(RESET): $(YELLOW)Démarre les services en arrière-plan (app et db).$(RESET)"
	@echo "$(GREEN)  stop   $(RESET): $(YELLOW)Arrête et supprime les conteneurs.$(RESET)"
	@echo "$(CYAN) In Container --------------------------- $(RESET)"
	@echo "$(GREEN)  logs   $(RESET): $(YELLOW)Affiche les logs des services en temps réel.$(RESET)"
	@echo "$(GREEN)  sh     $(RESET): $(YELLOW)Accède au shell du conteneur de l'application.$(RESET)"
	@echo "$(GREEN)  sh-db  $(RESET): $(YELLOW)Accède au shell du conteneur de la base de données.$(RESET)"
	@echo "$(GREEN)  psql   $(RESET): $(YELLOW)Ouvre l'interface psql connectée à la base de données.$(RESET)"

build:
	@echo "$(CYAN)Construction de l'image Docker de l'application...$(RESET)"
	@$(DC) build

dev:
	@echo "$(CYAN)Démarrage des services en arrière-plan...$(RESET)"
	@$(DC) up -d

stop:
	@echo "$(CYAN)Arrêt et suppression des conteneurs...$(RESET)"
	@$(DC) down

logs:
	@echo "$(CYAN)Affichage des logs des services en temps réel...$(RESET)"
	@$(D) logs -f $(APP)

sh:
	@echo "$(CYAN)Accès au shell du conteneur de l'application...$(RESET)"
	@$(D) exec -it $(APP) bash

sh-db:
	@echo "$(CYAN)Accès au shell du conteneur de la base de données...$(RESET)"
	@$(D) exec -it $(APP_DB) bash

psql:
	@echo "$(CYAN)Ouverture de l'interface psql connectée à la base de données...$(RESET)"
	@$(DC) exec -it $(APP_DB) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)
