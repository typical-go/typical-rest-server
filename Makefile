-include .env

PROJECT_NAME := $(shell basename "$(PWD)")
BINARY := $(PROJECT_NAME)
BINARY_LINUX := $(PROJECT_NAME)_linux
MOCK_TARGET := test/mock

## install: Install missing dependencies.
install: go-dep

## build: Build the binary.
build: go-dep go-build

## clean: Clean build files. Runs `go clean` internally.
clean:
	@rm -rf vendor
	@-$(MAKE) go-clean

## mock: Generate mock class
mock:
	@./mockgen.sh $(MOCK_TARGET)

go-build:
	@echo "  >  Building binary..."
	@go build -o $(PROJECT_NAME)

go-dep:
	@echo "  >  Checking if there is any missing dependencies..."
	@go get github.com/golang/dep/cmd/dep
	@go install github.com/golang/dep/cmd/dep
	@$(GOPATH)/bin/dep ensure

go-clean:
	@echo "  >  Cleaning build cache"
	@go clean

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
