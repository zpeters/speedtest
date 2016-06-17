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
	git submodule update --init --recursive
	go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest

clean:
	rm bin/speedtest*

test:
	go test $(shell glide nv)

coverage:
	go test -cover
	go test ./internal/... -cover

cross:
	echo "Building darwin-amd64..."
	GOOS="darwin" GOARCH="amd64" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-mac-amd64-${VERSION}

	echo "Building windows-386..."
	GOOS="windows" GOARCH="386" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-32-${VERSION}.exe

	echo "Building windows-amd64..."
	GOOS="windows" GOARCH="amd64" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-64-${VERSION}.exe

	echo "Building freebsd-386..."
	GOOS="freebsd" GOARCH="386" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-freebsd-386-${VERSION}

	echo "Building linux-arm..."
	GOOS="linux" GOARCH="arm" go  build -o bin/speedtest-linux-arm-${VERSION}

	echo "Building linux-386..."
	GOOS="linux" GOARCH="386" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-linux-386-${VERSION}

	echo "Building linux-amd64..."
	GOOS="linux" GOARCH="amd64" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-linux-amd64-${VERSION}

localweb: cross
	scp -v -P 22 bin/speed* root@zachpeters.org:/var/www/html/files/speedtest/

deploy: cross
	echo "Uploading..."
	ssh thehelpfulhacker.net "mkdir -p ~/media.thehelpfulhacker.net/speedtest/${VERSION}"
	scp -v bin/*${VERSION}* thehelpfulhacker.net:~/media.thehelpfulhacker.net/speedtest/${VERSION}/
