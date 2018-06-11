#!/bin/sh
set -e
dir=$(dirname $0)
cd $dir

pkg="github.com/letiantech/pipeline"

rm -rf test.log

runtest(){
    go test tests/*
    go test -cover -covermode=set -coverpkg=`go list ./...` ./tests
    go test -bench="." tests/*
}

runtest | tee $dir/test.log
