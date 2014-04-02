#!/bin/bash
export GOPATH=`pwd`
go install github.com/zpeters/speedtest
bin/speedtest $@
