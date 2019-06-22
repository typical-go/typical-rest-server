-include .env

PROJECT_NAME := $(shell basename "$(PWD)")

build:
	@go build -o bin/typical ./cmd/typical
	@go build -o bin/$(PROJECT_NAME) ./cmd/app