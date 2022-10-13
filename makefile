# Common variables
VERSION := 0.0.1
BUILD_INFO := Manual build 
FRONTEND_DIR := ./frontend
FRONTEND_HOST_DIR := ./frontend-host
BACKEND_DIR := ./backend

# Most likely want to override these when calling `make image`
IMAGE_REG ?= ghcr.io
IMAGE_REPO ?= benc-uk/nano-realms
IMAGE_TAG ?= latest
IMAGE_PREFIX := $(IMAGE_REG)/$(IMAGE_REPO)

# Things you don't want to change
REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
# Tools
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint              # Remove if not using Go
AIR_PATH := $(REPO_DIR)/bin/air                           # Remove if not using Go

.PHONY: help image push build run lint lint-fix
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
	@echo "Not implemented yet!"; exit 1
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	cd $(SRC_DIR); $(GOLINT_PATH) run --modules-download-mode=mod *.go
	cd $(SRC_DIR); npm run lint

lint-fix: ## 📝 Lint & format, attempts to fix errors & modify code
	@figlet $@ || true
	@echo "Not implemented yet!"; exit 1
	@$(GOLINT_PATH) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh
	cd $(SRC_DIR); golangci-lint run --modules-download-mode=mod *.go --fix
	cd $(SPA_DIR); npm run lint-fix

image: ## 📦 Build container image from Dockerfile
	@figlet $@ || true
	docker build --file ./build/Dockerfile \
	--build-arg BUILD_INFO="$(BUILD_INFO)" \
	--build-arg VERSION="$(VERSION)" \
	--tag $(IMAGE_PREFIX):$(IMAGE_TAG) . 

push: ## 📤 Push container image to registry
	@figlet $@ || true
	docker push $(IMAGE_PREFIX):$(IMAGE_TAG)

build: ## 🔨 Run a local build without a container
	@figlet $@ || true
	@echo "Not implemented yet!"
	#go build -o __CHANGE_ME__ $(SRC_DIR)/...
	#cd $(SRC_DIR); npm run build

run-backend: ## 🏃 Run backend with hot reload
	@figlet $@ || true
	cd $(BACKEND_DIR); $(AIR_PATH)

run-frontend: ## 🏃 Run frontend with hot reload
	@figlet $@ || true
	cd $(FRONTEND_DIR); npm run serve
