#!/bin/sh
set -e
dir=$(dirname $0)
cd $dir

rm -rf test.log

runtest(){
    go test -cover -bench="."
}

runtest | tee $dir/test.log
