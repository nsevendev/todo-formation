# ————————————————————————————————————————
# 🎵 🐳 Docker Makefile Todo-formation API 🐳 🎵
# ————————————————————————————————————————

# Variables
DOCKER = docker
DOCKER_COMP = docker compose
NAME_CONT_API = todo-f-api
APP_CONT_API = $(DOCKER_COMP) exec $(NAME_CONT_API) bash
PYTHON = $(DOCKER) exec $(NAME_CONT_API) python
PYTHON_IT = $(DOCKER) exec -it $(NAME_CONT_API) python
PIP = $(DOCKER) exec $(NAME_CONT_API) pip

# Misc
.DEFAULT_GOAL = help
.PHONY        : help check diffsettings runserver test createsuperuser dbshell migrate makemigrations startapp pip-list pip-install pip-upgrade update-requirements

## —— 🎵 🐍 Django Commands 🐍 🎵 —————————————————————————————————————————

help: ## Affiche cette aide
	@grep -E '(^[a-zA-Z0-9\./_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

check: ## Vérifie la configuration Django
	@$(PYTHON) manage.py check

diffsettings: ## Affiche les paramètres Django modifiés
	@$(PYTHON) manage.py diffsettings

runserver: ## Démarre le serveur Django avec les logs en mode verbeux
	@$(PYTHON) manage.py runserver 0.0.0.0:8000 --verbosity 3

test: ## Lance tous les tests Django avec verbose=2
	@$(PYTHON) manage.py test --verbosity 2

test-cibling: ## Lance un test spécifique (exemple : c=tests.test_home.HomePageTest.test_home_page_status_code)
	@if [ -z "$(c)" ]; then \
		echo "❌ Spécifie le chemin du test avec c=nom_du_test (ex : c=tests.test_home.HomePageTest.test_home_page_status_code)"; \
		exit 1; \
	fi
	@$(PYTHON) manage.py test $(c) --verbosity 2

test-app: ## Teste une app spécifique (exemple : c=nom_de_lapp)
	@if [ -z "$(c)" ]; then echo "❌ Spécifie l'app avec c=nom_de_lapp"; exit 1; fi
	@$(PYTHON) manage.py test $(c) --verbosity 2

coverage: ## Lancer les tests avec un rapport de couverture
	@$(PYTHON) -m coverage run --source=. manage.py test
	@$(PYTHON) -m coverage report -m

createsuperuser: ## Crée un super utilisateur pour l'admin Django
	@$(PYTHON) manage.py createsuperuser

dbshell: ## Ouvre une console PostgreSQL via Django
	@$(PYTHON_IT) manage.py dbshell

migrate: ## Applique les migrations à la base de données
	@$(PYTHON) manage.py migrate

makemigrations: ## Génère les migrations Django
	@$(PYTHON) manage.py makemigrations

startapp: ## Crée une nouvelle application Django (utilisation : make startapp c=nom_de_lapp)
	@if [ -z "$(c)" ]; then echo "❌ Erreur : spécifiez c=nom_de_lapp"; exit 1; fi
	@$(PYTHON) manage.py startapp $(c)

## —— 📦 Gestion des Packages 📦 —————————————————————————————————————————

pip-list: ## Liste les packages Python installés
	@$(PIP) list

pip-install: ## Installe un package (utilisation : make pip-install c=nom_du_package)
	@if [ -z "$(c)" ]; then echo "❌ Erreur : spécifiez c=nom_du_package"; exit 1; fi
	@$(PIP) install $(c) --root-user-action=ignore

pip-upgrade: ## Met à jour un package (utilisation : make pip-upgrade c=nom_du_package)
	@if [ -z "$(c)" ]; then echo "❌ Erreur : spécifiez c=nom_du_package"; exit 1; fi
	@$(PIP) install --upgrade $(c)

pip-install-all: ## Installe tous les packages depuis requirements/base.txt
	@$(PIP) install -r requirements/base.txt

update-requirements: ## Met à jour requirements/base.txt avec les packages installés
	@echo "🚀 Mise à jour de requirements/base.txt"
	@$(PIP) freeze > requirements/base.txt
	@echo "✅ Fichier requirements/base.txt mis à jour."

## —— 🚀 Docker & Utilitaires 🚀 —————————————————————————————————————————

logs: ## Affiche les logs du conteneur Django
	@$(DOCKER) logs -f $(NAME_CONT_API)

sh: ## Ouvre un shell bash dans le conteneur Django
	@$(DOCKER) exec -it $(NAME_CONT_API) bash

clean: ## Nettoie les fichiers temporaires (__pycache__, *.pyc)
	@echo "🧹 Nettoyage des fichiers temporaires"
	@$(APP_CONT_API) -c "find . -type d -name '__pycache__' -exec rm -rf {} +"
	@$(APP_CONT_API) -c "find . -type f -name '*.pyc' -delete"
	@echo "✅ Nettoyage terminé."

