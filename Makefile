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
	#go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-${VERSION} ./cmd/speedtest
	go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest ./cmd/speedtest

clean:
	rm -f bin/speedtest*
