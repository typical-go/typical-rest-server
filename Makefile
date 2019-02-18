-include .env

PROJECT_NAME := $(shell basename "$(PWD)")
BINARY := $(PROJECT_NAME)
MOCK_TARGET := test/mock

## install: Install missing dependencies.
install: go-dep

## build: Build the binary.
build: go-dep go-build

## clean: Clean build files. Runs `go clean` internally.
clean:
	@rm -rf vendor
	@-$(MAKE) go-clean

## mockgen: Generate mock class
mockgen:
	@./mockgen.sh $(MOCK_TARGET)

## test: Running test
test: go-dep go-test

## test-detail: Show test detail
test-detail:
	@-$(MAKE) test
	@go tool cover -html=cover.out

go-build:
	@echo "  >  Building binary..."
	@go build -o $(BINARY)

go-dep:
	@echo "  >  Checking if there is any missing dependencies..."
	@go get github.com/golang/dep/cmd/dep
	@go install github.com/golang/dep/cmd/dep
	@$(GOPATH)/bin/dep ensure

go-test:
	@go test ./config ./app/controller ./app/repository  -coverprofile cover.out

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
