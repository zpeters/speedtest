#!/usr/bin/env bash

export GOPATH=~/go
source ~/src/golang-crosscompile/crosscompile.bash

BRANCH=`git describe` 

echo "Building darwin-amd64..."
go-darwin-amd64 build -o bin/speedtest-mac-amd64-$BRANCH

echo "Building windows-386..."
go-windows-386 build -o bin/speedtest-32-$BRANCH.exe

echo "Building windows-amd64..."
go-windows-amd64 build -o bin/speedtest-64-$BRANCH.exe

echo "Building freebsd-amd64..."
go-freebsd-amd64 build -o bin/speedtest-freebsd-amd64-$BRANCH

echo "Building linux-arm..."
go-linux-arm build -o bin/speedtest-linux-arm-$BRANCH

echo "Building linux-amd64..."
go-linux-amd64 build -o bin/speedtest-linux-amd64-$BRANCH
