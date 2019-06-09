-include .env

build:
	@go build -o bin/typical ./typical
	@go build -o bin/app ./app