#!/bin/sh
VERSION=$(git describe --tags --always)


echo "Building darwin-amd64..."
GOOS="darwin" GOARCH="amd64" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-mac-amd64-${VERSION} ./cmd/speedtest

echo "Building windows-386..."
GOOS="windows" GOARCH="386" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-32-${VERSION}.exe ./cmd/speedtest

echo "Building windows-amd64..."
GOOS="windows" GOARCH="amd64" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-64-${VERSION}.exe ./cmd/speedtest

echo "Building freebsd-386..."
GOOS="freebsd" GOARCH="386" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-freebsd-386-${VERSION} ./cmd/speedtest

echo "Building linux-arm..."
GOOS="linux" GOARCH="arm" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-linux-arm-${VERSION} ./cmd/speedtest

echo "Building linux-386..."
GOOS="linux" GOARCH="386" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-linux-386-${VERSION} ./cmd/speedtest

echo "Building linux-amd64..."
GOOS="linux" GOARCH="amd64" go build -ldflags="-X main.Version=${VERSION}" -o bin/speedtest-linux-amd64-${VERSION} ./cmd/speedtest
