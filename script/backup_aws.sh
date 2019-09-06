#!/bin/bash
mongodump
time=$(date +"%Y-%m-%dT%H%M%S")
filename="csi-$time.tgz"
tar -zcvf ./$filename ./dump
rm -rf ./dump
aws s3 cp ./$filename s3://csidbbackup/$filename --profile lazypic
rm $filename
