install:
	go mod download

build:
	go build -o chopdoc ./chopdoc.go

test: 
	go test -v ./...

lint:
	golangci-lint run -v
