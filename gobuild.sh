#!/bin/bash

SRCDIR=$(pwd)

export GOPATH="$GOPATH:$SRCDIR"
export GO111MODULE="off"

go build -o build/wcserver src/main.go
