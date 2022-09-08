#!/bin/bash
/usr/bin/mongodump
time=$(date +"%Y-%m-%dT%H%M%S")
filename="openpipelineio-$time.tgz"
tar -zcvf ./$filename ./dump
rm -rf ./dump
cp ./$filename /idea/app/openpipelineio/dump
rm -f $filename