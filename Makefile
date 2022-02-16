help: ## Print this message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
NAME := webhook-go
#VERSION := $(shell date +'v%y%m.%-d.%-H%M%S')
#DATE := $(shell date +'%y.%m.%d-%H:%M:%S')
#SHA := $(shell git rev-parse HEAD)

test: ## Run go tests
	@go test ./...

build: ## Build a local binary
	@goreleaser build --single-target --snapshot --rm-dist
	@mkdir -p bin
	@cp dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME) bin/

run: ## Run webhook-go
	@cp webhook.yml.example webhook.yml
	@go run main.go

clean: ## Clean up build
	@echo "Cleaning Go environment..."
	@go clean
	@echo "Cleaning build directory..."
	@rm -rf dist/
	@echo "Cleaning local bin directory..."
	@rm -rf bin/

snapshot: ## Build artifacts without releasing
	goreleaser release --snapshot --rm-dist

compile: release ## Alias for release
release: ## Build for all supported OSes
	goreleaser release
