#!/bin/bash

#export GOPATH="$GOPATH:$SRCDIR"
export GO111MODULE="on"

go build -o build/wcserver main.go
