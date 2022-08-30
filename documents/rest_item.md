# RestAPI Item

restAPI의 장점은 웹서비스의 URI를 이용하기 때문에 네트워크만 연결되어있으면 OS, 디바이스 제약없이 사용할 수 있습니다.
또한 Python 같은 언어를 이용해서 중간 관리를 위한 API를 작성하더라도 OS별로 코드가 서로 달라지는 상황이 없습니다.

이 문서는 기본 restAPI옵션을 설명하고 파이썬을 이용해서 RestAPI를 사용하는 방법을 다룹니다.
파이프라인에 사용될 확률이 높은 코드라서, 일부 에러처리까지 코드로 다루었습니다.

# RestAPI for Item(Shot, Asset)

## Get

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api2/item | 아이템 가지고 오기 | project, name or id, (type) | `$ curl -H "Authorization: Basic <Token>" "https://openpipeline.io/api2/item?project=TEMP&name=SS_0020&type=org"` |
| /api3/items | 아이템을 검색하고 가지고 오기 | project, searchword, status상태 | `$ curl -H "Authorization: Basic <Token>" "https://openpipeline.io/api3/items?project=TEMP&searchword=SS&wip=true"` 또는 `$ curl -H "Authorization: Basic <Token>" "https://openpipeline.io/api3/items?project=TEMP&searchword=task:mm+user:jason&assign=true&wip=true"` |
| /api3/items | 아이템을 검색하고 가지고 오기(유연한 Status) | project, searchword, searchbartemplate, truestatus | `$ curl -H "Authorization: Basic <Token>" "https://openpipeline.io/api3/items?project=TEMP&searchword=SS&searchbartemplate=searchbarV2&truestatus=assing,wip"` |
| /api/shot | 샷 정보 가지고 오기 | project, name | `$ curl -H "Authorization: Basic <Token>" "https://openpipeline.io/api/shot?project=TEMP&name=SS_0010"` |
| /api2/shots | 샷 리스트를 가지고 오기 | project, seq | `$ curl -H "Authorization: Basic <Token>" "https://openpipeline.io/api/shots?project=TEMP&seq=SS"` |
| /api/allshots | 전체 샷 리스트를 가지고 오기 | project | `$ curl -H "Authorization: Basic <Token>" "https://openpipeline.io/api/allshots?project=TEMP"` |
| /api/asset | 에셋 정보 가지고 오기 | project, name | `$ curl -H "Authorization: Basic <Token>" "https://openpipeline.io/api/asset/asset?project=TEMP&name=stone01"` |
| /api/assets | 에셋 리스트를 가지고 오기 | project | `$ curl -H "Authorization: Basic <Token>" "https://openpipeline.io/api/assets?project=TEMP"` |
| /api/usetypes | 샷의 usetype 리스트 가지고오기 | project, name | `$ curl -H "Authorization: Basic <Token>" "https://openpipeline.io/api/usetypes?project=TEMP&name=SS_0010"` |
| /api/publishkeys | 존재하는 Publish Key 를 가지고 온다 | | `$ curl -X GET -H "Authorization: Basic <Token>" https://openpipeline.io/api/publishkeys` |

## Post

| URI | description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/search | 검색 | project, searchword, sortkey | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&searchword=SS_0020&sortkey=id" https://openpipeline.io/api/search` |
| /api/deadline2d | 2D마감일 리스트 | project | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP" https://openpipeline.io/api/deadline2d` |
| /api/deadline3d | 3D마감일 리스트 | project | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP" https://openpipeline.io/api/deadline3d` |
| /api/rmitemid | 아이템 삭제 | project, id | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=circle&id=SS_0010_org" https://openpipeline.io/api/rmitemid` |
| /api/settaskstatus | 상태수정 | project, name, task, status | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=circle&name=SS_0010&task=comp&status=wip" https://openpipeline.io/api/settaskstatus` |
| /api2/settaskstatus | 상태수정 | project, name, task, status | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=circle&name=SS_0010&task=comp&status=wip" https://openpipeline.io/api2/settaskstatus` |
| /api2/settaskuser | 사용자설정 | project, id, task, user | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&task=comp&user=d10191(김한웅,개발팀)" https://openpipeline.io/api2/settaskuser` |
| /api/settaskstartdate | 시작일 | project, name, task, date | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=RR_0010&task=comp&date=0506" https://openpipeline.io/api/settaskstartdate` |
| /api/settaskpredate | 1차마감일 | project, name, task, date | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=RR_0010&task=comp&date=0506" https://openpipeline.io/api/settaskpredate` |
| /api/settaskdate | 2차마감일 | project, name, task, date | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=RR_0010&task=comp&date=0506" https://openpipeline.io/api/settaskdate` |
| /api2/settaskmov | mov등록 | project, name, task, mov | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=RR_0010&task=comp&mov=/show/test/test.mov" https://openpipeline.io/api2/settaskmov` |
| /api/setshottype | shottype 변경 | project, name, type | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0030&shottype=3d" https://openpipeline.io/api/setshottype` |
| /api/setusetype | usetype 변경 | project, id, type | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0030_org&type=org1" https://openpipeline.io/api/setusetype` |
| /api2/setthummov | 썸네일mov변경 | project, name, path | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0030&path=/show/thumbnail.mov" https://openpipeline.io/api2/setthummov` |
| /api/setbeforemov | 썸네일 이전 mov 등록 | project, name, path, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0030&path=/show/before.mov" https://openpipeline.io/api/setbeforemov` |
| /api/setaftermov | 썸네일 이후 mov 등록 | project, name, path, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0030&path=/show/after.mov" https://openpipeline.io/api/setaftermov` |
| /api/seteditmov | 편집본 mov 등록 | project, id, path | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0030_org&path=/show/edit.mov" https://openpipeline.io/api/seteditmov` |
| /api/setassettype | assettype 변경 | project, name, type | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=mamma&type=prop" https://openpipeline.io/api/setassettype` |
| /api/setoutputname | 아웃풋이름 등록 | project, name, outputname | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&outputname=S101_010_010" https://openpipeline.io/api/setoutputname` |
| /api2/setrnum | 롤넘버 등록 | project, id, rnum | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&rnum=A0001" https://openpipeline.io/api2/setrnum` |
| /api/setdeadline2d | 2D마감일 등록 | project, name, date | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&date=0712" https://openpipeline.io/api/setdeadline2d` |
| /api/setdeadline3d | 3D마감일 등록 | project, name, date | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&date=0712" https://openpipeline.io/api/setdeadline3d` |
| /api/setscantimecodein | 스캔 타임코드IN 등록 | project, name, timecode | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" https://openpipeline.io/api/setscantimecodein` |
| /api/setscantimecodeout | 스캔 타임코드OUT 등록 | project, name, timecode | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" https://openpipeline.io/api/setscantimecodeout` |
| /api/setscanin | scan in frame 등록 | project, name, frame | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&frame=654231" https://openpipeline.io/api/setscanin` |
| /api/setscanout | scan out frame 등록 | project, name, frame | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&frame=654331" https://openpipeline.io/api/setscanout` |
| /api/setscanframe | scan frame 등록 | project, name, frame | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&frame=100" https://openpipeline.io/api/setscanframe` |
| /api/setjusttimecodein | JUST 타임코드IN 등록 | project, name, timecode | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" https://openpipeline.io/api/setjusttimecodein`|
| /api/setjusttimecodeout | JUST 타임코드OUT 등록 | project, name, timecode | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" https://openpipeline.io/api/setjusttimecodeout` |
| /api/setfinver | 최종데이터 버전 등록 | project, name, version | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&version=1" https://openpipeline.io/api/setfinver` |
| /api/setfindate | 최종데이터 아웃풋 날짜 등록 | project, name, date | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&date=0711" https://openpipeline.io/api/setfindate` |
| /api/setplatein | plate in frame 등록 | project, name, frame | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&frame=1001" https://openpipeline.io/api/setplatein` |
| /api/setplateout | plate out frame 등록 | project, name, frame | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&frame=1130" https://openpipeline.io/api/setplateout` |
| /api/setjustin | just in frame 등록 | project, name, frame | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&frame=1003" https://openpipeline.io/api/setjustin` |
| /api/setjustout | just out frame 등록 | project, name, frame | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&frame=1130" https://openpipeline.io/api/setjustout` |
| /api/sethandlein | handle in frame 등록 | project, name, frame | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&frame=1003" https://openpipeline.io/api/sethandlein` |
| /api/sethandleout | handle out frame 등록 | project, name, frame | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&frame=1130" https://openpipeline.io/api/sethandleout` |
| /api/addtag | tag 추가 | project, id, tag | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&tag=테스트" https://openpipeline.io/api/addtag` |
| /api/addassettag | assettag 추가 | project, id, assettag | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&assettag=테스트" https://openpipeline.io/api/addassettag` |
| /api/rmtag | tags 삭제 | project, id, tag, (iscontain) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0020_org&tag=태그3&iscontain=true" https://openpipeline.io/api/rmtag` |
| /api/rmassettag | assettag 삭제 | project, id, assettag, (iscontain) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0020_org&assettag=태그3&iscontain=true" https://openpipeline.io/api/rmassettag` |
| /api/setnote | 작업내용 변경 | project, name, text, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&text=바람이 휘날린다" https://openpipeline.io/api/setnote` |
| /api/addcomment | 수정사항 추가 | project, name, text, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&text=1003프레임 나무제거" https://openpipeline.io/api/addcomment` |
| /api/rmcomment | 수정사항 삭제 | project, name, date | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&date=2021-01-26T11:38:53%2B09:00" https://openpipeline.io/api/rmcomment` |
| /api/addsource | 링크소스 추가 | project, name, title, path, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&title=source1&path=/show/src1/test.mov" https://openpipeline.io/api/addsource` |
| /api/rmsource | 링크소스 삭제 | project, name, title, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&title=sourcename" https://openpipeline.io/api/rmsource` |
| /api/setcameraprojection | 카메라 프로젝션여부 | project, id, projection | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0020_org&projection=true" https://openpipeline.io/api/setcameraprojection` |
| /api/setcamerapubtask | 카메라 Pub Task설정 | project, id, task | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0020_org&task=mm" https://openpipeline.io/api/setcamerapubtask` mm,layout,ani 만 task 등록가능|
| /api/setcamerapubpath | 카메라 Pub Path설정 | project, id, path | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0020_org&path=/show/test/cam/pubpath" https://openpipeline.io/api/setcamerapubpath`|
| /api/setcameralensmm | 카메라 렌즈mm 설정 | project, id, lensmm | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0020_org&lensmm=45" https://openpipeline.io/api/setcameralensmm`|
| /api/setdeadline2d | 2D 마감일 설정 | project, name, date | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&date=2019-09-05" https://openpipeline.io/api/setdeadline2d`|
| /api/setdeadline3d | 3D 마감일 설정 | project, name, date | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&date=2019-09-05" https://openpipeline.io/api/setdeadline3d`|
| /api/setretimeplate | Retime Plate 경로설정 | project, name, path | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&path=/show/retime" https://openpipeline.io/api/setretimeplate`|
| /api/settasklevel | Task 레벨설정 | project, name, task, level | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&task=comp&level=1" https://openpipeline.io/api/settasklevel`|
| /api/setplatesize | Plate Size 설정 | project, name, size, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&size=2048x1152" https://openpipeline.io/api/setplatesize`|
| /api/setundistortionsize | Undistortion Size 설정 | project, name, size, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&size=2048x1152" https://openpipeline.io/api/setundistortionsize`|
| /api2/setrendersize | Reder Size 설정 | project, id, size | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0020_org&size=2048x1152" https://openpipeline.io/api2/setrendersize`|
| /api/setoverscanratio | Overscan Ratio 비율 설정 | project, id, ratio | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0020_org&ratio=1.1" https://openpipeline.io/api/setoverscanratio`|
| /api/setobjectid | ObjectID 설정 | project, name, in, out, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&in=100&out=200&userid=khw7096" https://openpipeline.io/api/setobjectid`|
| /api/setociocc | OCIO .cc 경로설정 | project, name, path, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&path=/show/color.cc" https://openpipeline.io/api/setociocc`|
| /api/setrollmedia | Settelite Rollmedia 설정 | project, name, rollmedia, (userid) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&rollmedia=A001_C0003_DFGE" https://openpipeline.io/api/setrollmedia`|
| /api/setscanname | Scanname 설정 | project, id, scanname | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0020_org&scanname=A001_C0003_DFGE" https://openpipeline.io/api/scanname`|
| /api/task | Task 정보를 가지고 온다. | project, name, task | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020&task=comp" https://openpipeline.io/api/task`|
| /api/shottype | Shottype 정보를 가지고 온다. | project, name | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0020" https://openpipeline.io/api/shottype`|
| /api/statusnum | project status 갯수를 가지고 온다. | project | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP" https://openpipeline.io/api/statusnum`|
| /api/taskstatusnum | task status 갯수를 가지고 온다. | project, task | `$ curl -H "Authorization: Basic <Token>" -d "project=TEMP&task=comp" https://openpipeline.io/api/taskstatusnum`|
| /api/taskanduserstatusnum | task, user status 갯수를 가지고 온다. | project, task, user | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&task=comp&user=jason" https://openpipeline.io/api/taskanduserstatusnum`|
| /api/userstatusnum | user status 갯수를 가지고 온다. | project, user | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&user=jason" https://openpipeline.io/api/userstatusnum`|
| /api/addpublish | Publish를 추가한다. status는 usethis, working, notuse 로 설정할 수 있다. | project, name, task, key, (secondarykey), path, status, (mainversion), (subversion), (subject), (filetype), (kindofusd), (createtime), (isoutput), (authornamekor), (outputdatapath) | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&name=SS_0010&task=comp&key=nuke&secondarykey=pub&path=/path/file.nk&mainversion=1&subversion=1&subject=roto&filetype=.usd&kindofusd=component&status=usethis&isoutput=true&authornamekor=김한웅&outputdatapath=toclient.mov" https://openpipeline.io/api/addpublish`|
| /api/setpublishstatus | Publish status를 변경한다. status는 usethis, working, notuse 로 설정할 수 있다. | project, id, task, key, path, status, createtime | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&task=comp&key=pub&path=/path/path&status=usethis&createtime=2020-05-21T09:00:00%2B09:00" https://openpipeline.io/api/setpublishstatus`|
| /api/rmpublishkey | Publish Key를 삭제한다. | project, id, task, key | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&task=comp&key=pub" https://openpipeline.io/api/rmpublishkey`|
| /api/rmpublish | Publish 를 삭제한다. | project, id, task, key, createtime, path | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&task=comp&key=pub&createtime=2020-05-21T09:00:00%2B09:00&path=/show/test" https://openpipeline.io/api/rmpublish`|
| /api/setseq | seq를 설정한다. | project, id, seq | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&seq=SS" https://openpipeline.io/api/setseq`|
| /api/setseason | season를 설정한다. | project, id, season | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&season=S01" https://openpipeline.io/api/setseason`|
| /api/setepisode | episode를 설정한다. | project, id, episode | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&episode=E01" https://openpipeline.io/api/setepisode`|
| /api/setnetflixid | netflixid를 설정한다. | project, id, netflixid | `$ curl -X POST -H "Authorization: Basic <Token>" -d "project=TEMP&id=SS_0010_org&netflixid=123456" https://openpipeline.io/api/setnetflixid`|
| /api/uploadthumbnail | 썸네일 업로드 | project, name, (type) | `$ curl -X POST -H "Authorization: Basic <Token>" -F project=TEMP -F name=SS_0010 -F "image=@/path/thumbnail.png" "https://openpipeline.io/api/uploadthumbnail"`|
| /api/item | 샷 추가 | json | `$ curl -X POST -H "Authorization: Basic <Token>" -d '{"project":"TEMP","id":"SS_0010_org","name":"SS_0010","type":"org","status":"4","statusv2":"assign","tasks":{"mm":{"title":"mm","status":"4","statusv2":"assign"}}}' "https://openpipeline.io/api/item"`|



#### URL Encode
`/path/test.%04d.exr` 형태의 데이터를 보내고 싶다면 url-encode를 처리해야합니다.
`%` 문자는 `%25` 값에 해당한다. 일일이 변환할 수 없기 때문에 curl에서는 --data-urlencode 명령어를 사용하면 됩니다.

```bash
$ curl -X POST \
--data-urlencode "project=TEMP" \
--data-urlencode "name=SS_0010" \
--data-urlencode "task=comp" \
--data-urlencode "key=pub" \
--data-urlencode "path=/path/file.%04d.exr" \
--data-urlencode "mainversion=1" \
--data-urlencode "subversion=1" \
--data-urlencode "subject=roto" \
--data-urlencode "kindofusd=component" \
--data-urlencode "status=usethis" \
https://openpipeline.io/api/publish
```

#### 샷정보 가지고오기. Python2.7x
- TEMP 프로젝트 OPN_0010 샷 정보를 가지고 오기(암호화 토큰키 사용)
- 일반적으로 샷은 org(일반), left(입체) 타입을 가지게 됩니다.

```python
#!/usr/bin/python
#coding:utf8
import urllib2
import json

endpoint = "https://openpipeline.io/api/shot?project=TEMP&name=SS_0010"
request = urllib2.Request(endpoint)
request.add_header("Authorization", "Basic <TOKEN-KEY>")
result = urllib2.urlopen(request)
data = json.load(result)
print(data)
```

#### 샷,에셋정보(Item) 가지고오기. Python3.7.x(블랜더)
토큰키를 이용한 GET
```python
#!/usr/bin/python
#coding:utf8

import requests
endpoint = "https://openpipeline.io/api/item?project=TEMP&id=SS_0020_org"
auth = {'Authorization': 'Basic JDJhJDEwJHBBREluL0JuRTdNa3NSb3RKZERUbWVMd0V6OVB1TndnUGJzd2k0RlBZcmEzQTBSczkueHZH'}
r = requests.get(url=endpoint, headers=auth)
print(r.json())
```

홈디렉토리에 `~/.csi/token` 파일 내부에 CSI 토큰키를 저장해 두었다면 아래 형태로 코드를 작성, POST 할 수 있습니다.

```python
#!/usr/bin/python
#coding:utf8

import os
import requests
from pathlib import Path
home = str(Path.home())
f = open(os.path.join(home, ".csi", "token"), "r")
token = f.read(80)
data = {'project':'TEMP', 'name':'SS_0010', 'text':'test'}
endpoint = "https://openpipeline.io/api/setoutputname"
auth = {'Authorization': 'Basic ' + token}
r = requests.post(url=endpoint, data=data, headers=auth)
print(r.json())
```

#### 샷,에셋(Item) 검색하기
- gunhamdo 프로젝트에서 "SS"문자열이 들어간 아이템을 검색하는 예제이다.
- 상태는 작업중인 아이템을 대상으로 한다.

```python
#coding:utf8
import json
import urllib2
endpoint = "https://openpipeline.io/api3/items?project=circle&searchword=SS&searchbartemplate=searchbarV2&truestatus=assign,wip"
data = json.load(urllib2.urlopen(endpoint))
print(data)
```

- url encoding 방법

```python
#coding:utf8
import json
import urllib
import urllib2

values = {}
values["project"] = "TEMP"
values["searchword"] = "SS"
values["truestatus"] = "assign,wip"
values["searchbartemplate"] = "searchbarV2"

endpoint = "https://openpipeline.io/api3/items"
query = urllib.urlencode(values)
endpoint = endpoint + "?" + query
data = json.load(urllib2.urlopen(endpoint))
print(data)
```

- 서치키워드에 다중 문자열 검색시 공백이 들어가면 공백을 + 사용

```
endpoint = "https://openpipeline.io/api3/items?project=TEMP&searchword=comp+배서영&searchbartemplate=searchbarV2&truestatus=wip"
```


#### 샷,에셋에 대한 이름(Name) 검색하기
- 샷, 에셋에 대하여 문자열이 포함된 이름(Name)으로 검색할 수 있다.
- adventure 프로젝트의 "R0VFX" 이름을 가진 샷, 에셋 정보를  검색하는 예제이다.

```python
#coding:utf8
import json
import urllib2

endpoint = "https://openpipeline.io/api/searchname?project=adventure&name=R0VFX"
data = json.load(urllib2.urlopen(endpoint))
print(data)
```

#### 시퀀스 리스트 검색하기
- 프로젝트에 대한 시퀀스 리스트를 검색할 수 있다.
- 'mkk3' 프로젝트의 시퀀스 리스트를 검색하는 예제이다.

```python
#coding:utf8
import json
import urllib2

endpoint = "https://openpipeline.io/api/seqs?project=mkk3"
data = json.load(urllib2.urlopen(endpoint))
print(data)
```

#### 샷 리스트 검색하기
- 해당 프로젝트의 시퀀스 문자열이 포함된 샷 리스트를 검색할 수 있다.
- 'mkk3'프로젝트의 "BNS"시퀀스에 대한 샷 리스트를 검색하는 예제이다.

```python
#coding:utf8
import json
import urllib2

endpoint = "https://openpipeline.io/api2/shots?project=mkk3&seq=BNS"
data = json.load(urllib2.urlopen(endpoint))
print(data)
```

#### curl 명령을 이용하여 mov등록하기

circle 프로젝트 SS_0010 샷에 light 테스크에 /show/test.mov 등록하기.

```bash
$ curl -d "project=circle&name=SS_0010&task=light&mov=/show/test.mov" https://openpipeline.io/api2/settaskmov
```

circle 프로젝트 mamma 에셋 fur 테스크에 /show/fur.mov 등록하기.

```bash
$ curl -d "project=circle&name=mamma&task=fur&mov=/show/fur.mov" https://openpipeline.io/api2/settaskmov
```

#### python에서 샷 mov등록하기
- 아래는 파이썬에서 urllib2.request를 이용하여 mov를 등록하는 예제이다.

```python
#coding:utf8
import json
import urllib2
data = "project=TEMP&name=AS_0010&task=comp&mov=test.mov"
request = urllib2.Request("https://openpipeline.io/api/settaskmov", data)
request.add_header("Authorization", "Basic JDJhJDEwJHBBREluL0JuRTdNa3NSb3RKZERUbWVMd0V6OVB1TndnUGJzd2k0RlBZcmEzQTBSczkueHZH")
data = json.load(urllib2.urlopen(request))
print(data)
```

#### curl을 사용한 날짜 지정

- Task 작업 시작일 설정
```bash
$ curl -X POST -d "project=TEMP&name=SS_0011&task=fx&startdate=0502" https://openpipeline.io/api/setstartdate
$ curl -X POST -d "project=TEMP&name=SS_0011&task=fx&startdate=2018-05-02" https://openpipeline.io/api/setstartdate
# + 문자는 web에서 %2B 이다. 터미널에서 curl로 테스트시 + 대신 %2B를 넣어주어야 한다.
$ curl -X POST -d "project=TEMP&name=SS_0011&task=fx&startdate=2018-05-02T14:45:34%2B09:00" https://openpipeline.io/api/setstartdate
```

- Task 1차 마감일 설정
```bash
$ curl -X POST -d "project=TEMP&name=SS_0011&task=fx&predate=0605" https://openpipeline.io/api/setpredate
$ curl -X POST -d "project=TEMP&name=SS_0011&task=fx&predate=2018-06-05" https://openpipeline.io/api/setpredate
# + 문자는 web에서 %2B 이다. 터미널에서 curl로 테스트시 + 대신 %2B를 넣어주어야 한다.
$ curl -X POST -d "project=TEMP&name=SS_0011&task=fx&predate=2018-06-05T14:45:34%2B09:00" https://openpipeline.io/api/setpredate
```

#### Python2.7.x 에서 publish값을 notuse로 바꾸는 예제

```python
#!/usr/bin/python
#coding:utf8
import urllib
import urllib2
import json

url = 'https://openpipeline.io/api/setpublishstatus'
header = {"Authorization": "Basic <Token>"}
data = { 'project' : 'TEMP',
		'id' : 'SS_0010_org',
		'task' : 'comp',
		'key' : 'maya',
		'status' : 'notuse',
		'createtime' : '2020-11-19T11:35:27+09:00',
		'path' : '/show/test.nk',
		}
req = urllib2.Request(url=url, data=urllib.urlencode(data), headers=header)
f = urllib2.urlopen(req)
print(json.load(f))
```

#### Python2.7.x 에서 publish 삭제 예제

```python
#!/usr/bin/python
#coding:utf8
import urllib
import urllib2
import json

url = 'https://openpipeline.io/api/rmpublish'
header = {"Authorization": "Basic <Token>"}
data = { 'project' : 'TEMP',
		'id' : 'SS_0010_org',
		'task' : 'comp',
		'key' : 'maya',
		'createtime' : '2020-11-19T11:35:27+09:00',
		'path' : '/show/test.nk',
		}
req = urllib2.Request(url=url, data=urllib.urlencode(data), headers=header)
f = urllib2.urlopen(req)
print(json.load(f))
```

#### Python3 에서 publish 삭제 예제
- python 3.4 이후 버전부터는 requests 모듈이 자동으로 내장되어 있습니다.
- python 3.4 이전 버전은 requests 모듈을 수동으로 설치해주어야 함

```bash
$ pip install requests
```

```python
#!/usr/bin/python
#coding:utf8
import requests
import json

url = 'https://openpipeline.io/api/rmpublish'
data = { 'project' : 'TEMP',
		'id' : 'SS_0010_org',
		'task' : 'comp',
		'key' : 'maya',
		'createtime' : '2020-09-09T11:32:24+09:00',
		'path' : '/show/project/path',
		}
header = {"Authorization": "Basic <Token>"}
result = requests.post(url=url, data=data, headers=header)
print(result.json())
```

# Python2.7에서 task publish key를 제거하는 코드
```python
#!/usr/bin/python
#coding:utf8
import urllib
import urllib2
import json

url = 'https://openpipeline.io/api/rmpublish'
data = { 'project' : 'TEMP',
	'id' : 'SS_0010_org',
	'task' : 'ani',
	'key' : 'cam',
	'createtime' : '2021-11-18T03:41:58Z',
	'path' : '/path/cam.mb',
}
data = urllib.urlencode(data)
request = urllib2.Request(url)
request.add_header("Authorization", "Basic <TOKEN>")
request.add_data(data)
response = urllib2.urlopen(request)
print(response.read())
