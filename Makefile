-include .env

build:
	@go build -o bin/typical ./typical/cli
	@go build -o bin/app ./app/cli