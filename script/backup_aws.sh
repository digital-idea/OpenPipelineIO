#!/bin/bash
/usr/local/bin/mongodump
time=$(date +"%Y-%m-%dT%H%M%S")
filename="openpipelineio-$time.tgz"
tar -zcvf ./$filename ./dump
rm -rf ./dump
/usr/local/bin/aws s3 cp ./$filename s3://openpipelineio-dbbackup/$filename --profile lazypic
rm $filename
