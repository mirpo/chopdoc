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

release-draft:
	goreleaser release --snapshot --draft

compare:
	go run ./chopdoc.go -input ./tests/pg_essay.txt -output ./recursive_60_0_go.jsonl -size 60 -overlap 0 -method recursive
	(cd tests && uv run ./recursive.py --size 60 --overlap 0 --input ./pg_essay.txt --output ../recursive_60_0_py.jsonl)
	./scripts/diff.sh ./recursive_60_0_py.jsonl ./recursive_60_0_go.jsonl
