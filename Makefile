.PHONY: all
all: help

.PHONY: submodule
submodule: ## Pull submodules
	git submodule update --init --recursive 

.PHONY: searx-up
searx-up: ## Deploy Searx
	docker-compose -f docker/docker-compose.yml up -d

.PHONY: searx-down
searx-down: ## Tear down Searx
	docker-compose -f docker/docker-compose.yml down

.PHONY: build
build: ## Build Emerald Tooling CLI
	@mkdir -p build/bin
	go build -o build/bin/emld-cli cmd/cli/main.go

.PHONY: help
help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'