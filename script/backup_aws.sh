#!/bin/bash
/usr/local/bin/mongodump
time=$(date +"%Y-%m-%dT%H%M%S")
filename="csi-$time.tgz"
tar -zcvf ./$filename ./dump
rm -rf ./dump
/usr/local/bin/aws s3 cp ./$filename s3://csidbbackup/$filename --profile lazypic
rm $filename
