#!/bin/sh
cat speedtest.go | sed "s/var VERSION = .*/var VERSION = '`git describe`'/g" > speedtest.go
