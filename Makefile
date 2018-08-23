#
#  Makefile for Go
#
SHELL=/usr/bin/env bash
VERSION=$(shell git describe --tags --always)
PACKAGES = $(shell find ./ -type d | grep -v 'vendor' | grep -v '.git' | grep -v 'bin')

.PHONY: list
.PHONY: test-cover-html

default: build

build:
	go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-${VERSION} ./cmd/speedtest

static:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-extldflags \"static\" -s -w" -o bin/speedtest ./cmd/speedtest
	upx bin/speedtest

clean:
	scripts/clean.sh

vet:
	go vet ./cmd/...
	go vet ./internal/...

lint:
	golint ./cmd/...
	golint ./internal/...

fmt:
	gofmt -w ./cmd/speedtest
	gofmt -w ./internal/app
	gofmt -w ./internal/pkg

test:
	go test ./cmd/... ./internal/...

update:
	dep ensure -update

cross:
	scripts/cross-compile.sh
