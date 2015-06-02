#!/usr/bin/env bash

export GOPATH=~/go
source ~/src/golang-crosscompile/crosscompile.bash
BRANCH=`git describe --tags` 

echo "Compiling..."
etc/cross-compile.sh

echo "Uploading..."
ssh zachpeters@thehelpfulhacker.net "mkdir ~/media.thehelpfulhacker.net/speedtest/$BRANCH"
scp -v bin/*$BRANCH* zachpeters@thehelpfulhacker.net:~/media.thehelpfulhacker.net/speedtest/$BRANCH/
