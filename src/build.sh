#!/bin/sh

mkdir -p ./build/
GOOS=darwin GOARCH=amd64 go build -o ./build/goDB-mac *.go
GOOS=linux GOARCH=amd64 go build -o ./build/goDB-linux *.go
