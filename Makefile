-include .env

PROJECT_NAME := $(shell basename "$(PWD)")

typical:
	@go build -o bin/typical ./cmd/typical

app:
	@go build -o bin/$(PROJECT_NAME) ./cmd/app

.PHONY: typical app