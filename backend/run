#!/bin/bash
export GOROOT=/opt/homebrew/Cellar/go/1.26.0/libexec
go generate ./...
go run ./cmd/server > server.log 2>&1 &
