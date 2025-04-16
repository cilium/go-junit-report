VERSION=$(shell git describe --match="v*")
REVISION=$(shell git rev-parse HEAD)
TIMESTAMP=$(shell date +%FT%T)

.PHONY: all
all: test build

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build ./cmd/go-junit-report

.PHONY: release
release: test
	goreleaser release --snapshot --clean

.PHONY: clean
clean:
	rm -rf dist/

.PHONY: clean release test
