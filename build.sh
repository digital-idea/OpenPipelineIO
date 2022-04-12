#!/bin/sh
APP="csi3"

# assets 폴더의 모든 에셋을 빌드전에 assets_vfsdata.go 파일로 생성한다.
go run assets/asset_generate.go

# OS별 기본빌드
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.SHA1VER=`git rev-parse HEAD` -X main.BUILDTIME=`date -u +%Y-%m-%dT%H:%M:%S`" -o ./bin/linux/${APP} *.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.SHA1VER=`git rev-parse HEAD` -X main.BUILDTIME=`date -u +%Y-%m-%dT%H:%M:%S`" -o ./bin/darwin/${APP} *.go
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.SHA1VER=`git rev-parse HEAD` -X main.DNS=csi.idea.co.kr -X main.BUILDTIME=`date -u +%Y-%m-%dT%H:%M:%S` -X main.COMPANY=digitalidea -X main.MAILDNS=idea.co.kr -X main.DBIP=10.0.90.249 -X main.DILOG=http://10.0.90.251:8080 -X main.WFS=http://10.0.98.20:8081" -o ./bin/linux_di/${APP} *.go

# Github Release에 업로드 하기위해 압축
cd ./bin/linux/ && tar -zcvf ../${APP}_linux_x86-64.tgz . && cd -
cd ./bin/darwin/ && tar -zcvf ../${APP}_darwin_x86-64.tgz . && cd -
cd ./bin/linux_di/ && tar -zcvf ../${APP}_linux_di_x86-64.tgz . && cd -

# 삭제
rm -rf ./bin/linux
rm -rf ./bin/darwin
rm -rf ./bin/linux_di