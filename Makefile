-include .env

PROJECT_NAME := $(shell basename "$(PWD)")

typical:
	@go build -o bin/typical ./cmd/typical

t: 
	@bin/typical $(filter-out $@,$(MAKECMDGOALS))

app:
	@go build -o bin/$(PROJECT_NAME) ./cmd/app

.PHONY: typical app 