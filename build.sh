#!/bin/sh
APP="csi3"

# assets 폴더의 모든 에셋을 빌드전에 assets_vfsdata.go 파일로 생성한다.
go run assets/asset_generate.go
#rice embed-go

# OS별 필드
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.TEMPLATEPATH=template" -o ./bin/linux/${APP} camera.go check.go cmdfunc.go crypto.go csi3.go db_item.go db_project.go db_setellite.go db_user.go http_api.go http_cmd.go http_item.go http_project.go http_setellite.go http_user.go http.go item.go log.go network.go parser.go project.go searchop.go setellite.go setellitefunc.go templatefunc.go usage.go user.go timefunc.go assets_vfsdata.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.TEMPLATEPATH=template" -o ./bin/darwin/${APP} camera.go check.go cmdfunc.go crypto.go csi3.go db_item.go db_project.go db_setellite.go db_user.go http_api.go http_cmd.go http_item.go http_project.go http_setellite.go http_user.go http.go item.go log.go network.go parser.go project.go searchop.go setellite.go setellitefunc.go templatefunc.go usage.go user.go timefunc.go assets_vfsdata.go
#GOOS=windows GOARCH=amd64 go build -ldflags "-X main.TEMPLATEPATH=template" -o ./bin/windows/${APP}.exe camera.go check.go cmdfunc.go crypto.go csi3.go db_item.go db_project.go db_setellite.go db_user.go http_api.go http_cmd.go http_item.go http_project.go http_setellite.go http_user.go http.go item.go log.go network.go parser.go project.go searchop.go setellite.go setellitefunc.go templatefunc.go usage.go user.go timefunc.go assets_vfsdata.go

GOOS=linux GOARCH=amd64 go build -ldflags "-X main.TEMPLATEPATH=template -X main.DBIP=10.0.90.251 -X main.DILOG=http://10.0.90.251:8080 -X main.WFS=http://10.0.98.20:8081 -X main.THUMBPATH=/lustre/INHouse/Tool/csidata/thumbnail -X main.MEDIAPATH=/lustre/INHouse/Tool/csi3/media -X main.TEMPLATEPATH=/lustre/INHouse/Tool/csi3/template" -o ./bin/linux_di/${APP} camera.go check.go cmdfunc.go crypto.go csi3.go db_item.go db_project.go db_setellite.go db_user.go http_api.go http_cmd.go http_item.go http_project.go http_setellite.go http_user.go http.go item.go log.go network.go parser.go project.go searchop.go setellite.go setellitefunc.go templatefunc.go usage.go user.go timefunc.go assets_vfsdata.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.TEMPLATEPATH=template -X main.DBIP=10.0.90.251 -X main.DILOG=http://10.0.90.251:8080 -X main.WFS=http://10.0.98.20:8081 -X main.THUMBPATH=/lustre/INHouse/Tool/csidata/thumbnail -X main.MEDIAPATH=/lustre/INHouse/Tool/csi3/media -X main.TEMPLATEPATH=/lustre/INHouse/Tool/csi3/template" -o ./bin/darwin_di/${APP} camera.go check.go cmdfunc.go crypto.go csi3.go db_item.go db_project.go db_setellite.go db_user.go http_api.go http_cmd.go http_item.go http_project.go http_setellite.go http_user.go http.go item.go log.go network.go parser.go project.go searchop.go setellite.go setellitefunc.go templatefunc.go usage.go user.go timefunc.go assets_vfsdata.go
#GOOS=windows GOARCH=amd64 go build -ldflags "-X main.TEMPLATEPATH=template -X main.DBIP=10.0.90.251 -X main.DILOG=http://10.0.90.251:8080 -X main.WFS=http://10.0.98.20:8081 -X main.THUMBPATH=/lustre/INHouse/Tool/csidata/thumbnail -X main.MEDIAPATH=/lustre/INHouse/Tool/csi3/media -X main.TEMPLATEPATH=/lustre/INHouse/Tool/csi3/template" -o ./bin/windows_di/${APP}.exe camera.go check.go cmdfunc.go crypto.go csi3.go db_item.go db_project.go db_setellite.go db_user.go http_api.go http_cmd.go http_item.go http_project.go http_setellite.go http_user.go http.go item.go log.go network.go parser.go project.go searchop.go setellite.go setellitefunc.go templatefunc.go usage.go user.go timefunc.go assets_vfsdata.go

# Github Release에 업로드 하기위해 압축
cd ./bin/linux/ && cp -r ../../assets/template . && mkdir thumbnail && tar -zcvf ../${APP}_linux_x86-64.tgz . && cd -
cd ./bin/darwin/ && cp -r ../../assets/template . && mkdir thumbnail && tar -zcvf ../${APP}_darwin_x86-64.tgz . && cd -
#cd ./bin/windows/ && cp -r ../../assets/template . && mkdir thumbnail && tar -zcvf ../${APP}_windows_x86-64.tgz . && cd -

cd ./bin/linux_di/ && cp -r ../../assets/template . && tar -zcvf ../${APP}_linux_di_x86-64.tgz . && cd -
cd ./bin/darwin_di/ && cp -r ../../assets/template . && tar -zcvf ../${APP}_darwin_di_x86-64.tgz . && cd -
#cd ./bin/windows_di/ && cp -r ../../assets/template . && tar -zcvf ../${APP}_windows_di_x86-64.tgz . && cd -

# 삭제
rm -rf ./bin/linux
rm -rf ./bin/darwin
#rm -rf ./bin/windows

rm -rf ./bin/linux_di
rm -rf ./bin/darwin_di
#rm -rf ./bin/windows_di