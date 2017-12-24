#!/bin/sh
dir=$(dirname $0)
cd $dir

rm -rf test.log

for t in `ls tests`
do
    go test tests/$t | tee test.log
done
