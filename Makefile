help: ## Print this message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

NAME := webhook-go
#VERSION := $(shell date +'v%y%m.%-d.%-H%M%S')
#DATE := $(shell date +'%y.%m.%d-%H:%M:%S')
#SHA := $(shell git rev-parse HEAD)

test: ## Run go tests
	@go test ./...

build: ## Build a local binary
	@go build -o bin/${NAME} .

run: ## Run webhook-go
	@cp webhook.yml.example webhook.yml
	@go run main.go

clean: ## Clean up build
	@echo "Cleaning Go environment..."
	@go clean
	@echo "Cleaning build directory..."
	@rm -rf build/
	@echo "Cleaning local bin directory..."
	@rm -rf bin/

compile: clean windows linux rpi rpisf ## Build for all supported OSes

windows: ## Build a Windows 64-bit binary
	@echo "Building Windows binary"
	GOOS=windows GOARCH=amd64 go build -o build/windows/webhook-go-windows-x84_64.exe .
	@zip build/webhook-go-windows-x84_64.zip build/windows/webhook-go-windows-x84_64.exe

linux: ## Build a Linux 64-bit binary
	@echo "Building Linux AMD64 binary"
	GOOS=linux GOARCH=amd64 go build -o build/linux/webhook-go-linux-x86_64 .
	@tar czf build/webhook-go-linux-x86_64.tar.gz build/linux/webhook-go-linux-x86_64

rpi: ## Build a Raspberry Pi Linux binary
	@echo "Building Linux ARM binary"
	GOOS=linux GOARCH=arm go build -o build/rpi/webhook-go-rpi .
	@tar cJf build/webhook-go-rpi.tar.xz build/rpi/webhook-go-rpi

rpisf: ## Build a Raspberry Pi 64-bit binary
	@echo "Building Linux ARM64 binary"
	GOOS=linux GOARCH=arm64 go build -o build/rpi/webhook-go-rpi64 .
	@tar cJf build/webhook-go-rpi64.tar.xz build/rpi/webhook-go-rpi64

