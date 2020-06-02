#!/bin/bash
/usr/bin/mongodump
time=$(date +"%Y-%m-%dT%H%M%S")
filename="csi-$time.tgz"
tar -zcvf ./$filename ./dump
rm -rf ./dump
cp ./$filename /idea/app/csi/dump
rm -f: $filename