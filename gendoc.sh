#!/bin/sh
godoc -index -html pipeline > doc/index.html
godoc -http=:6060 -goroot=./
