test:
	go test ./... -v

build:
	go build -o bin/ ./...

run:
	go run main.go

.PHONY: test build run
.DEFAULT_GOAL := build