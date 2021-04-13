#!/bin/bash

set -e
set -x

go get ./...
mkdir -p build
go build -o build/koinos-log-golang *.go
