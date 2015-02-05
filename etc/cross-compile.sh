#!/usr/bin/env bash

export GOPATH=~/go
source ~/src/golang-crosscompile/crosscompile.bash

echo darwin-amd64
go-darwin-amd64 build -o bin/speedtest-mac-amd64 

echo windows-386
go-windows-386 build -o bin/speedtest-32.exe

echo windows-amd64
go-windows-amd64 build -o bin/speedtest-64.exe

echo freebsd-amd64
go-freebsd-amd64 build -o bin/speedtest-freebsd-amd64 

echo linux-arm
go-linux-arm build -o bin/speedtest-linux-arm 

echo linux-amd64
go-linux-amd64 build -o bin/speedtest-linux-amd64
