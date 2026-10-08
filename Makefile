all: build

fmt:
	goimports -w .

build: 
	go build ./...

test:
	go test -race ./...

lint:
	golangci-lint run ./...
	golangci-lint fmt --diff

.PHONY: all fmt build test lint
