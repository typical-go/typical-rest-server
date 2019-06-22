-include .env

build:
	@go build -o bin/typical ./typical/cmd
	@go build -o bin/app ./app/cmd