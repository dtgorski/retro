.PHONY: help clean dist-clean build run tidy

VERSION  := $(shell 2>/dev/null git describe --always --tags --dirty)
GO_BUILD := CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)"
GO_TEST  := CGO_ENABLED=1 go test -count=1 -race -v -coverprofile=../coverage.out

help:                   # Displays this list
	@echo; grep -P "^[a-z][a-zA-Z0-9_<> -]+:.*(?=#)" Makefile | sed -E "s/:[^#]*?#?(.*)?/\r\t\t\t\1/" | uniq | sed "s/^/ make /"
	@echo

clean:                  # Removes build/test artifacts
	@2>/dev/null rm ./coverage.html || true
	@2>/dev/null rm -rf ./bin || true

dist-clean: clean       # Removes Go caches and artefacts debris
	@go clean -cache -testcache

build: clean            # Builds executable binary
	@$(GO_BUILD) -o ./bin/ ./cmd/retro

run: build              # Runs the application
	@trap "" TERM INT EXIT; stty -echoctl && ./bin/retro $(ARGS) || true

tidy:                   # Formats source files, cleans go.mod
	@find . -type f -not -path "*/\.*" -name "*.go" | xargs -I{} gofmt -w {}
	@go mod tidy
