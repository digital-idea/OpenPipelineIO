#!/bin/sh
APP="csi3"

# assets 폴더의 모든 에셋을 빌드전에 assets_vfsdata.go 파일로 생성한다.
go run assets/asset_generate.go

# OS별 기본빌드
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.MINVER=`date -u +.%Y%m%d.%H%M%S`" -o ./bin/linux/${APP} *.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.MINVER=`date -u +.%Y%m%d.%H%M%S`" -o ./bin/darwin/${APP} *.go
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.MINVER=`date -u +.%Y%m%d.%H%M%S`" -o ./bin/windows/${APP}.exe *.go

# 디지털아이디어 빌드
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.DNS=csi.idea.co.kr -X main.MINVER=`date -u +.%Y%m%d.%H%M%S` -X main.COMPANY=digitalidea -X main.MAILDNS=idea.co.kr -X main.DBIP=10.0.90.251 -X main.DILOG=http://10.0.90.251:8080 -X main.WFS=http://10.0.98.20:8081 -X main.THUMBPATH=/lustre/INHouse/Tool/csidata/thumbnail" -o ./bin/linux_di/${APP} *.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.DNS=csi.idea.co.kr -X main.MINVER=`date -u +.%Y%m%d.%H%M%S` -X main.COMPANY=digitalidea -X main.MAILDNS=idea.co.kr -X main.DBIP=10.0.90.251 -X main.DILOG=http://10.0.90.251:8080 -X main.WFS=http://10.0.98.20:8081 -X main.THUMBPATH=/lustre/INHouse/Tool/csidata/thumbnail" -o ./bin/darwin_di/${APP} *.go
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.DNS=csi.idea.co.kr -X main.MINVER=`date -u +.%Y%m%d.%H%M%S` -X main.COMPANY=digitalidea -X main.MAILDNS=idea.co.kr -X main.DBIP=10.0.90.251 -X main.DILOG=http://10.0.90.251:8080 -X main.WFS=http://10.0.98.20:8081 -X main.THUMBPATH=/lustre/INHouse/Tool/csidata/thumbnail" -o ./bin/windows_di/${APP}.exe *.go

# Github Release에 업로드 하기위해 압축
cd ./bin/linux/ && mkdir thumbnail && tar -zcvf ../${APP}_linux_x86-64.tgz . && cd -
cd ./bin/darwin/ && mkdir thumbnail && tar -zcvf ../${APP}_darwin_x86-64.tgz . && cd -
cd ./bin/windows/ && mkdir thumbnail && tar -zcvf ../${APP}_windows_x86-64.tgz . && cd -

cd ./bin/linux_di/ && mkdir thumbnail && tar -zcvf ../${APP}_linux_di_x86-64.tgz . && cd -
cd ./bin/darwin_di/ && mkdir thumbnail && tar -zcvf ../${APP}_darwin_di_x86-64.tgz . && cd -
cd ./bin/windows_di/ && mkdir thumbnail && tar -zcvf ../${APP}_windows_di_x86-64.tgz . && cd -

# 삭제
rm -rf ./bin/linux
rm -rf ./bin/darwin
rm -rf ./bin/windows

rm -rf ./bin/linux_di
rm -rf ./bin/darwin_di
rm -rf ./bin/windows_di
