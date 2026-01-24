help: ## Print this message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

#GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
NAME := webhook-go
VERSION := $(shell git describe --tags | cut -d"." -f1)
#DATE := $(shell date +'%y.%m.%d-%H:%M:%S')
#SHA := $(shell git rev-parse HEAD)

test: ## Run go tests
	@go test ./...

binary: ## Build a local binary
	@goreleaser build --single-target --clean
	@mkdir -p bin
	@cp dist/$(NAME)_$(GOOS)_$(GOARCH)_$(VERSION)/$(NAME) bin/

run: ## Run webhook-go
	@go run main.go --config ./build/webhook.yml

clean: ## Clean up build
	@echo "Cleaning Go environment..."
	@go clean
	@echo "Cleaning build directory..."
	@rm -rf dist/
	@echo "Cleaning local bin directory..."
	@rm -rf bin/

compile: ## Build for all supported OSes
	goreleaser build --snapshot --clean
snapshot: ## Build artifacts without releasing
	goreleaser release --snapshot --clean
release: ## Build release for all supported OSes
	goreleaser release
