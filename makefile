# Common variables
VERSION := 0.0.1
BUILD_INFO := Manual build
FRONTEND_DIR := ./frontend
FRONTEND_HOST_DIR := ./frontend-host
BACKEND_DIR := ./backend

# Most likely want to override these when calling `make image`
IMAGE_REG ?= ghcr.io
IMAGE_NAME ?= benc-uk/nano-realms
IMAGE_TAG ?= latest

# Used only when deploying to container apps
AZURE_REGION ?= uksouth
AZURE_RESGRP ?= nano-realms
AZURE_BASE_NAME ?= nanorealms

# Things you don't want to change
REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
# Tools
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint
AIR_PATH := $(REPO_DIR)/bin/air

.EXPORT_ALL_VARIABLES:
.PHONY: help images push build run lint lint-fix deploy world-build world-show clean
.DEFAULT_GOAL := help

help: ## 💬 This help message :)
	@figlet $@ || true
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-tools: ## 🔮 Install dev tools into project bin directory
	@figlet $@ || true
	@$(GOLINT_PATH) > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin/
	@$(AIR_PATH) -v > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh

lint: ## 🔍 Lint & format check only, sets exit code on error for CI
	@figlet $@ || true
	$(GOLINT_PATH) run --modules-download-mode=mod
	cd $(FRONTEND_DIR); npm run lint

lint-fix: ## 📝 Lint & format, attempts to fix errors & modify code
	@figlet $@ || true
	$(GOLINT_PATH) run --modules-download-mode=mod --fix
	cd $(FRONTEND_DIR); npm run lint-fix

image-database: ## 📦 Build container image for the database
	@figlet $@ || true
	docker compose -f build/compose.yaml build database
image-frontend: ## 📦 Build container image for the frontend
	@figlet $@ || true
	docker compose -f build/compose.yaml build database
image-backend: ## 📦 Build container image for the backend
	@figlet $@ || true
	docker compose -f build/compose.yaml build database
images: ## 📦 Build all container images
	@figlet $@ || true
	docker compose -f build/compose.yaml build

push: ## 📤 Push container images to registry
	@figlet $@ || true
	docker compose -f build/compose.yaml push

build: ## 🔨 Run a local build without a container
	@figlet $@ || true
	go build -o bin/backend \
	  -ldflags "-X main.version=$(VERSION) -X 'main.buildInfo=$(BUILD_INFO)'" \
	  nano-realms/backend 
	go build -o bin/frontend-host nano-realms/frontend-host
	cd $(FRONTEND_DIR); npm run build

run-backend: ## 🏃 Run backend with hot reload
	@figlet $@ || true
	$(AIR_PATH)

run-frontend: ## 🏃 Run frontend with hot reload
	@figlet $@ || true
	cd $(FRONTEND_DIR); npm run serve

run-db: ## 🔁 Run Neo4J database
	@figlet $@ || true
	docker run -p 7474:7474 -p 7687:7687 \
	--volume=$(REPO_DIR)/data:/data \
	--env=NEO4J_AUTH=none neo4j

clean: ## 🧹 Clean up the repo
	@figlet $@ || true
	rm -rf bin
	sudo rm -rf data
	rm -rf tmp

world-build: ## 🌎 (Re)build the world graph
	@figlet $@ || true
	@./world/run-script.sh world/realm.cypher

world-show: ## 📃 Dump a report of the world state
	@figlet $@ || true
	@./world/run-script.sh world/dump.cypher

deploy: ## 🚀 Deploy to Azure Container Apps
	@figlet $@ || true
	@./deploy/container-apps/deploy.sh