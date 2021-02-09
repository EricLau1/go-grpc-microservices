#!/bin/bash

go clean --cache && go test -v -cover microservices/...
go build -o authentication/authsvc authentication/main.go
go build -o api/apisvc api/main.go
