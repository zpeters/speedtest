#!/usr/local/bin/bash
export GOPATH=`pwd`
source ~/src/golang-crosscompile/crosscompile.bash

echo windows-386
go-windows-386 build -o bin/speedtest-32.exe speedtest

echo windows-amd64
go-windows-amd64 build -o bin/speedtest-64.exe speedtest

echo freebsd-arm
go-freebsd-arm build -o bin/speedtest-freebsd-arm speedtest

echo freebsd-amd64
go-freebsd-amd64 build -o bin/speedtest-freebsd-amd64 speedtest

echo linux-arm
go-linux-arm build -o bin/speedtest-linux-arm speedtest

echo linux-amd64
go-linux-amd64 build -o bin/speedtest-linux-amd64 speedtest
