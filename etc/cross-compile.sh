#!/usr/bin/env bash

export GOPATH=~/go

BRANCH=`git describe --tags` 

echo "Building darwin-amd64..."
GOOS=darwin GOARCH=amd64 go build -o bin/speedtest-mac-amd64-$BRANCH

echo "Building windows-386..."
GOOS=windows GOARCH=386 go build -o bin/speedtest-32-$BRANCH.exe

echo "Building windows-amd64..."
GOOS=windows GORCH=amd64 go build -o bin/speedtest-64-$BRANCH.exe

echo "Building freebsd-amd64..."
GOOS=freebsd GOARCH=amd64 go build -o bin/speedtest-freebsd-amd64-$BRANCH

echo "Building linux-arm..."
GOOS=linux GOARCH=arm go build -o bin/speedtest-linux-arm-$BRANCH

echo "Building linux-amd64..."
GOOS=linux GOARCH=amd64 go build -o bin/speedtest-linux-amd64-$BRANCH
