#!/bin/bash

PROJECT_NAME=TIX-SESSION-GO
PROJECT_PATH=github.com/tiket/TIX-SESSION-GO # relative to $GOPATH
PROJECT_REPO=git@github.com:tiket/TIX-SESSION-GO.git

# setup git
rm -rf .git
git init
git remote add origin $PROJECT_REPO

# clean up
mv Project_README.md README.md
rm LICENSE.md
echo $PROJECT_NAME >> .gitignore

# move the project to correct path
mkdir -p $GOPATH/src/$PROJECT_PATH
cp -r . $GOPATH/src/$PROJECT_PATH
cd $GOPATH/src/$PROJECT_PATH
rm -rf $GOPATH/src/imantung/typical