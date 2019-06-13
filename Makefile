-include .env

build:
	@go build -o bin/typical ./_typical
	@go build -o bin/app ./app