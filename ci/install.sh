#!/bin/bash

sudo gem install coveralls-lcov
go get -u github.com/jandelgado/gcov2lcov
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.45.2
