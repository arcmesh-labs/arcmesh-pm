VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X github.com/arcmesh-labs/arcmesh-pm/cmd.version=$(VERSION)
MODULE  := github.com/arcmesh-labs/arcmesh-pm

.PHONY: build install release clean test test-integration

build:
	go build -ldflags "$(LDFLAGS)" -o bin/apm .

install:
	go build -ldflags "-s -w -X github.com/arcmesh-labs/arcmesh-pm/cmd.version=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)" -o $(HOME)/gopath/bin/apm .

release:
	mkdir -p dist
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/apm-linux-amd64 .
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/apm-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/apm-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/apm-windows-amd64.exe .

clean:
	rm -rf bin/ dist/

test:
	go test ./...

test-integration:
	@bash tests/integration/run.sh
