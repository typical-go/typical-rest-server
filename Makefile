-include .env

PROJECT_NAME := $(shell basename "$(PWD)")
BINARY := $(PROJECT_NAME)
MOCK_TARGET := test/mock

## all: install missing dependency and build the binary
all: dep test build

## build: Build the binary.
build:
	@echo "  >  Building binary..."
	@go build -o $(BINARY)

## test: Running test
test:
	@echo "  >  Running test..."
	@go test ./config ./app/controller ./app/repository  -coverprofile cover.out

## dep: install missing dependency
dep:
	@echo "  >  Install missing dependencies..."
	@go get github.com/golang/dep/cmd/dep
	@go install github.com/golang/dep/cmd/dep
	@$(GOPATH)/bin/dep ensure

dep-clean:
	@echo "  >  Clean dependencies..."
	@rm -rf vendor

## clean: Clean build files
clean:
	@echo "  >  Clean build files..."
	@rm $(BINARY)
	@-$(MAKE) go-clean

## clean: Clean build files and dependency
clean-all: dep-clean clean

## mock: Generate mock class
mock:
	@echo "  >  Generate mock class..."
	@./mockgen.sh $(MOCK_TARGET)

## test-report: Show test report
test-report:
	@-$(MAKE) test
	@go tool cover -html=cover.out

.PHONY: help test
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
