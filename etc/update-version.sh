#!/bin/sh
cat speedtest.go | sed "s/var VERSION = .*/var VERSION = \"`git describe --tags`\"/g" > speedtest.go
