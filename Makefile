VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)

install:
	go mod download

build:
	CGO_ENABLED=0 go build -ldflags="-w -s -X main.version=${VERSION} -X main.commit=${COMMIT}" -o chopdoc ./chopdoc.go

test: 
	go test -v ./...

lint:
	which golangci-lint
	golangci-lint run --verbose
