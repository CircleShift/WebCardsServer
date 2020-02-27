#!/bin/bash

SRCDIR=$(pwd)

GOPATH="$GOPATH:$SRCDIR"

go env -w GOPATH=$GOPATH

go build -o build/server src/main.go
