#!/bin/bash

if [ -z ${1} ]; then
  echo "mock target is missing; check makefile for usage"
  exit 1
fi

go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen

# generate mock for repository
for filename in app/repository/*_repository.go; do
  $GOPATH/bin/mockgen -source=$filename -destination=${1}/$(basename "$filename") -package=${1##*/}
done
