# RestAPI Item
restAPI의 장점은 웹서비스의 URI를 이용하기 때문에 네트워크만 연결되어있으면 OS, 디바이스 제약없이 사용할 수 있습니다.
또한 Python 같은 언어를 이용해서 사내 API를 작성하더라도 OS별 코드가 서로 달라지는 상황이 없습니다.

이 문서는 기본 restAPI옵션을 설명하고 파이썬을 이용해서 RestAPI를 사용하는 방법을 다룹니다.
파이프라인에 사용될 확률이 높은 코드라서, 일부 에러처리까지 코드로 다루었습니다.

# RestAPI for Item

## Get

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/item | 아이템 가지고 오기 | project, id | `$ curl "https://csi.lazypic.org/api/item?project=TEMP&id=SS_0020_org"` |
| /api/shot | 샷 정보 가지고 오기 | project, name | `$ curl "https://csi.lazypic.org/api/shot?project=TEMP&name=SS_0010"` |
| /api/shots | 샷 리스트를 가지고 오기 | project, seq | `$ curl "https://csi.lazypic.org/api/shots?project=TEMP&seq=SS"` |
| /api/allshots | 전체 샷 리스트를 가지고 오기 | project | `$ curl "https://csi.lazypic.org/api/allshots?project=TEMP"` |
| /api/asset | 에셋 정보 가지고 오기 | project, name | `$ curl "https://csi.lazypic.org/api/asset/asset?project=TEMP&name=stone01"` |
| /api/assets | 에셋 리스트를 가지고 오기 | project | `$ curl "https://csi.lazypic.org/api/assets?project=TEMP"` |
| /api/usetypes | 샷의 usetype 리스트 가지고오기 | project, name | `$ curl "https://csi.lazypic.org/api/usetypes?project=TEMP&name=SS_0010"` |

## Post

| URI | description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/search | 검색 | project, searchword, sortkey | `$ curl -d "project=TEMP&searchword=SS_0020&sortkey=id" https://csi.lazypic.org/api/search` |
| /api/deadline2d | 2D마감일 리스트 | project | `$ curl -d "project=TEMP" https://csi.lazypic.org/api/deadline2d` |
| /api/deadline3d | 3D마감일 리스트 | project | `$ curl "https://csi.lazypic.org/api/deadline3d?project=TEMP"` |
| /api/rmitemid | 아이템 삭제 | project, id | `$ curl -d "project=circle&id=SS_0010_org" https://csi.lazypic.org/api/rmitemid` |
| /api/settaskstatus | 상태수정 | project, name, task, status | `$ curl -d "project=circle&name=SS_0010&task=comp&status=wip" https://csi.lazypic.org/api/setstatus` |
| /api/setassigntask | Assign 설정,해제 | project, name, task, status | `$ curl -d "project=TEMP&name=SS_0030&task=mg&status=true" https://csi.lazypic.org/api/setassigntask` |
| /api/settaskuser | 사용자수정 | project, name, task, user | `$ curl -d "project=TEMP&name=mamma&task=light&user=김한웅" https://192.168.219.104/api/settaskuser` |
| /api/settaskstartdate | 시작일 | project, name, task, date | `$ curl -d "project=TEMP&name=RR_0010&task=comp&date=0506" https://csi.lazypic.org/api/settaskstartdate` |
| /api/settaskpredate | 1차마감일 | project, name, task, date | `$ curl -d "project=TEMP&name=RR_0010&task=comp&date=0506" https://csi.lazypic.org/api/settaskpredate` |
| /api/settaskdate | 2차마감일 | project, name, task, date | `$ curl -d "project=TEMP&name=RR_0010&task=comp&date=0506" https://csi.lazypic.org/api/settaskdate` |
| /api2/settaskmov | mov등록 | project, name, task, mov | `$ curl -d "project=TEMP&name=RR_0010&task=comp&mov=/show/test/test.mov" https://csi.lazypic.org/api2/settaskmov` |
| /api/setshottype | shottype 변경 | project, name, type | `$ curl -d "project=TEMP&name=SS_0030&shottype=3d" https://csi.lazypic.org/api/setshottype` |
| /api/setusetype | usetype 변경 | project, id, type | `$ curl -d "project=TEMP&id=SS_0030_org&type=org1" https://csi.lazypic.org/api/setusetype` |
| /api/setthummov | 썸네일mov변경 | project, name, path, (userid) | `$ curl -d "project=TEMP&name=SS_0030&path=/show/thumbnail.mov" https://csi.lazypic.org/api/setthummov` |
| /api/setbeforemov | 썸네일mov변경 | project, name, path, (userid) | `$ curl -d "project=TEMP&name=SS_0030&path=/show/before.mov" https://csi.lazypic.org/api/setbeforemov` |
| /api/setaftermov | 썸네일mov변경 | project, name, path, (userid) | `$ curl -d "project=TEMP&name=SS_0030&path=/show/after.mov" https://csi.lazypic.org/api/setaftermov` |
| /api/setassettype | assettype 변경 | project, name, type | `$ curl -d "project=TEMP&name=mamma&type=prop" https://csi.lazypic.org/api/setassettype` |
| /api/setoutputname | 아웃풋이름 등록 | project, name, outputname | `$ curl -d "project=TEMP&name=SS_0010&outputname=S101_010_010" https://csi.lazypic.org/api/setoutputname` |
| /api/setrnum | 롤넘버 등록 | project, name, rnum | `$ curl -d "project=TEMP&name=SS_0010&rnum=A0001" https://csi.lazypic.org/api/setrnum` |
| /api/setdeadline2d | 2D마감일 등록 | project, name, date | `$ curl -d "project=TEMP&name=SS_0010&date=0712" https://csi.lazypic.org/api/setdeadline2d` |
| /api/setdeadline3d | 3D마감일 등록 | project, name, date | `$ curl -d "project=TEMP&name=SS_0010&date=0712" https://csi.lazypic.org/api/setdeadline3d` |
| /api/setscantimecodein | 스캔 타임코드IN 등록 | project, name, timecode | `$ curl -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" https://csi.lazypic.org/api/setscantimecodein` |
| /api/setscantimecodeout | 스캔 타임코드OUT 등록 | project, name, timecode | `$ curl -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" https://csi.lazypic.org/api/setscantimecodeout` |
| /api/setscanin | scan in frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=654231" https://csi.lazypic.org/api/setscanin` |
| /api/setscanout | scan out frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=654331" https://csi.lazypic.org/api/setscanout` |
| /api/setscanframe | scan frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=100" https://csi.lazypic.org/api/setscanframe` |
| /api/setjusttimecodein | JUST 타임코드IN 등록 | project, name, timecode | `$ curl -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" https://csi.lazypic.org/api/setjusttimecodein`|
| /api/setjusttimecodeout | JUST 타임코드OUT 등록 | project, name, timecode | `$ curl -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" https://csi.lazypic.org/api/setjusttimecodeout` |
| /api/setfinver | 최종데이터 버전 등록 | project, name, version | `$ curl -d "project=TEMP&name=SS_0010&version=1" https://csi.lazypic.org/api/setfinver` |
| /api/setfindate | 최종데이터 아웃풋 날짜 등록 | project, name, date | `$ curl -d "project=TEMP&name=SS_0010&date=0711" https://csi.lazypic.org/api/setfindate` |
| /api/setplatein | plate in frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1001" https://csi.lazypic.org/api/setplatein` |
| /api/setplateout | plate out frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1130" https://csi.lazypic.org/api/setplateout` |
| /api/setjustin | just in frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1003" https://csi.lazypic.org/api/setjustin` |
| /api/setjustout | just out frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1130" https://csi.lazypic.org/api/setjustout` |
| /api/sethandlein | handle in frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1003" https://csi.lazypic.org/api/sethandlein` |
| /api/sethandleout | handle out frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1130" https://csi.lazypic.org/api/sethandleout` |
| /api/addtag | tag 추가 | project, id, tag | `$ curl -d "project=TEMP&id=SS_0010_org&tag=테스트" https://192.168.219.104/api/addtag` |
| /api/rmtag | tags 삭제 | project, id, tag | `$ curl -d "project=TEMP&id=SS_0020_org&tag=태그3" https://192.168.219.114/api/rmtag` |
| /api/setnote | 작업내용 변경 | project, name, text, (userid) | `$ curl -d "project=TEMP&name=SS_0020&text=바람이 휘날린다" https://192.168.219.104/api/setnote` |
| /api/addcomment | 수정사항 추가 | project, name, text, (userid) | `$ curl -d "project=TEMP&name=SS_0020&text=1003프레임 나무제거" https://192.168.219.104/api/addcomment` |
| /api/rmcomment | 수정사항 삭제 | project, name, text, (userid) | `$ curl -d "project=TEMP&name=SS_0020&text=1003프레임 나무제거" https://192.168.219.104/api/rmcomment` |
| /api/addsource | 링크소스 추가 | project, name, title, path, (userid) | `$ curl -d "project=TEMP&name=SS_0020&title=source1&path=/show/src1/test.mov" https://csi.lazypic.org/api/addsource` |
| /api/rmsource | 링크소스 삭제 | project, name, title, (userid) | `$ curl -d "project=TEMP&name=SS_0020&title=sourcename" https://csi.lazypic.org/api/rmsource` |
| /api/setcameraprojection | 카메라 프로젝션여부 | project, name, projection, (userid) | `$ curl -d "project=TEMP&name=SS_0020&projection=true" https://csi.lazypic.org/api/setcameraprojection` |
| /api/setcamerapubtask | 카메라 Pub Task설정 | project, name, task, (userid) | `$ curl -d "project=TEMP&name=SS_0020&task=mm" https://csi.lazypic.org/api/setcamerapubtask` mm,layout,ani 만 task 등록가능|
| /api/setcamerapubpath | 카메라 Pub Path설정 | project, name, path, (userid) | `$ curl -d "project=TEMP&name=SS_0020&path=/show/test/cam/pubpath" https://csi.lazypic.org/api/setcamerapubpath`|
| /api/setdeadline2d | 2D 마감일 설정 | project, name, date | `$ curl -d "project=TEMP&name=SS_0020&date=2019-09-05" https://csi.lazypic.org/api/setdeadline2d`|
| /api/setdeadline3d | 3D 마감일 설정 | project, name, date | `$ curl -d "project=TEMP&name=SS_0020&date=2019-09-05" https://csi.lazypic.org/api/setdeadline3d`|
| /api/setretimeplate | Retime Plate 경로설정 | project, name, path | `$ curl -d "project=TEMP&name=SS_0020&path=/show/retime" https://csi.lazypic.org/api/setretimeplate`|
| /api/settasklevel | Task 레벨설정 | project, name, task, level | `$ curl -d "project=TEMP&name=SS_0020&task=comp&level=1" https://csi.lazypic.org/api/settasklevel`|
| /api/setplatesize | Plate Size 설정 | project, name, size, (userid) | `$ curl -d "project=TEMP&name=SS_0020&size=2048x1152" https://csi.lazypic.org/api/setplatesize`|
| /api/setundistortionsize | Undistortion Size 설정 | project, name, size, (userid) | `$ curl -d "project=TEMP&name=SS_0020&size=2048x1152" https://csi.lazypic.org/api/setundistortionsize`|
| /api/rendersize | Reder Size 설정 | project, name, size, (userid) | `$ curl -d "project=TEMP&name=SS_0020&size=2048x1152" https://csi.lazypic.org/api/setrendersize`|
| /api/setobjectid | ObjectID 설정 | project, name, in, out, (userid) | `$ curl -d "project=TEMP&name=SS_0020&in=100&out=200&userid=khw7096" https://csi.lazypic.org/api/setobjectid`|
| /api/setociocc | OCIO .cc 경로설정 | project, name, path, (userid) | `$ curl -d "project=TEMP&name=SS_0020&path=/show/color.cc" https://csi.lazypic.org/api/setociocc`|
| /api/setrollmedia | Settelite Rollmedia 설정 | project, name, rollmedia, (userid) | `$ curl -d "project=TEMP&name=SS_0020&rollmedia=A001_C0003_DFGE" https://csi.lazypic.org/api/setrollmedia`|
| /api/setscanname | Scanname 설정 | project, id, scanname | `$ curl -d "project=TEMP&id=SS_0020_org&scanname=A001_C0003_DFGE" https://csi.lazypic.org/api/scanname`|
| /api/task | Task 정보를 가지고 온다. | project, name, task | `$ curl -d "project=TEMP&name=SS_0020&task=comp" https://csi.lazypic.org/api/task`|
| /api/shottype | Shottype 정보를 가지고 온다. | project, name | `$ curl -d "project=TEMP&name=SS_0020" https://csi.lazypic.org/api/shottype`|



#### 샷정보 가지고오기. Python2.7x
- TEMP 프로젝트 OPN_0010 샷 정보를 가지고 오기(암호화 토큰키 사용)
- 일반적으로 샷은 org(일반), left(입체) 타입을 가지게 됩니다.

```python
#!/usr/bin/python
#coding:utf-8
import urllib2
import json
try:
    request = urllib2.Request("https://csi.lazypic.org/api/shot?project=TEMP&name=OPN_0010")
    key = "JDJhJDEwJHBBREluL0JuRTdNa3NSb3RKZERUbWVMd0V6OVB1TndnUGJzd2k0RlBZcmEzQTBSczkueHZH"
    request.add_header("Authorization", "Basic %s" % key)
    result = urllib2.urlopen(request)
    data = json.load(result)
except:
    print("RestAPI에 연결할 수 없습니다.")
    # 이후 에러처리 할 것
if "error" in data:
    print(data["error"])
print(data["data"])
```

#### 샷,에셋정보(Item) 가지고오기. Python3.7.x(블랜더)
토큰키를 이용한 GET
```python
import requests
endpoint = "https://csi.lazypic.org/api/item?project=TEMP&id=SS_0020_org"
auth = {'Authorization': 'Basic JDJhJDEwJHBBREluL0JuRTdNa3NSb3RKZERUbWVMd0V6OVB1TndnUGJzd2k0RlBZcmEzQTBSczkueHZH'}
r = requests.get(url=endpoint, headers=auth)
print(r.json())
```

사용자 토큰기를 `~/.csi/token` 파일에 저장해 두었다면 아래 형태로 코드를 작성, POST 할 수 있습니다.
```python
import os
import requests
from pathlib import Path
home = str(Path.home())
f = open(os.path.join(home, ".csi", "token"), "r")
token = f.read(80)
data = {'project':'TEMP', 'name':'SS_0010', 'text':'test'}
endpoint = "https://csi.lazypic.org/api/setoutputname"
auth = {'Authorization': 'Basic ' + token}
r = requests.post(url=endpoint, data=data, headers=auth)
print(r.json())
```

#### 샷,에셋(Item) 검색하기
- 군함도에서 "SS"문자열이 들어간 아이템을 검색하는 예제이다.
- 상태는 작업중인 아이템을 대상으로 한다.

```
#coding:utf-8
import json
import urllib2

restURL = "https://csi.lazypic.org/api2/items?project=gunhamdo&searchword=SS&wip=true"
try:
	data = json.load(urllib2.urlopen(restURL))
except:
	print("RestAPI에 연결할 수 없습니다.")
	# 에러처리
if "error" in data:
	print(data["error"])
	# 에러처리
print(data["data"])
```

- url encoding 방법

```
#coding:utf-8
import json
import urllib
import urllib2

values = {}
values["project"] = "yeomryeok"
values["searchword"] = "shot"
values["wip"] = "true"

url = "https://csi.lazypic.org/api2/items"
query = urllib.urlencode(values)
restURL = url + "?" + query
try:
	data = json.load(urllib2.urlopen(restURL))
except:
	print("RestAPI에 연결할 수 없습니다.")
	# 에러처리
if "error" in data:
	print(data["error"])
	# 에러처리
print(data["data"])
```
- 서치키워드에 다중 문자열 검색시 공백이 들어가면 공백을 + 사용

```
restURL = "https://csi.lazypic.org/api2/items?project=TEMP&searchword=comp+박지섭&wip=true"
```

- url endcode 방법

```
values = {}
values["project"] = "TEMP"
values["searchword"] = "comp+김한웅"
values["wip"] = "true"

url = "https://csi.lazypic.org/api2/items"
query = urllib.urlencode(values)
restURL = url + "?" + query
```

#### 샷,에셋에 대한 이름(Name) 검색하기
- 샷, 에셋에 대하여 문자열이 포함된 이름(Name)으로 검색할 수 있다.
- adventure 프로젝트의 "R0VFX" 이름을 가진 샷, 에셋 정보를  검색하는 예제이다.

```
- https://csi.lazypic.org/api/searchname?project=adventure&name=R0VFX
- https://csi.lazypic.org/api/searchname?project=adventure&name=R0VFX_sh001
```
```
#coding:utf-8
import json
import urllib2

restURL = "https://csi.lazypic.org/api/searchname?project=adventure&name=R0VFX_sh001"
try:
	data = json.load(urllib2.urlopen(restURL))
except:
	print("RestAPI에 연결할 수 없습니다.")
	# 에러처리
if "error" in data:
	print(data["error"])
	# 에러처리
print(data["data"])
```

#### 시퀀스 리스트 검색하기
- 프로젝트에 대한 시퀀스 리스트를 검색할 수 있다.
- 'mkk3' 프로젝트의 시퀀스 리스트를 검색하는 예제이다.

```
- https://csi.lazypic.org/api/seqs?project=mkk3
```
```
#coding:utf-8
import json
import urllib2

restURL = "https://csi.lazypic.org/api/seqs?project=mkk3"
try:
	data = json.load(urllib2.urlopen(restURL))
except:
	print("RestAPI에 연결할 수 없습니다.")
	# 에러처리
if "error" in data:
	print(data["error"])
	# 에러처리
print(data["data"])
```

#### 샷 리스트 검색하기
- 해당 프로젝트의 시퀀스 문자열이 포함된 샷 리스트를 검색할 수 있다.
- 'mkk3'프로젝트의 "BNS"시퀀스에 대한 샷 리스트를 검색하는 예제이다.

```
- https://csi.lazypic.org/api/shots?project=mkk3&seq=BNS
- https://csi.lazypic.org/api/shots?project=mkk3&seq=BNS_00
```
```
#coding:utf-8
import json
import urllib2

restURL = "https://csi.lazypic.org/api/shots?project=mkk3&seq=BNS"
try:
	data = json.load(urllib2.urlopen(restURL))
except:
	print("RestAPI에 연결할 수 없습니다.")
	# 에러처리
if "error" in data:
	print(data["error"])
	# 에러처리
print(data["data"])
```

#### curl을 이용해서 샷 mov등록하기
터미널에서 Curl명령을 이용하여 mov를 등록할 수 있습니다.

circle 프로젝트 SS_0010 샷에 light 테스크에 /show/test.mov 등록하기.

```bash
$ curl -d "project=circle&name=SS_0010&task=light&mov=/show/test.mov" https://csi.lazypic.org/api2/settaskmov
```

circle 프로젝트 mamma 에셋 fur 테스크에 /show/fur.mov 등록하기.

```bash
$ curl -d "project=circle&name=mamma&task=fur&mov=/show/fur.mov" https://csi.lazypic.org/api2/settaskmov
```

#### python에서 샷 mov등록하기
- 아래는 파이썬에서 urllib2.request를 이용하여 mov를 등록하는 예제이다.

```python
#coding:utf8
import json
import urllib2
data = "project=TEMP&name=AS_0010&task=comp&mov=test.mov"
request = urllib2.Request("https://csi.lazypic.org/api/settaskmov", data)
key = "JDJhJDEwJHBBREluL0JuRTdNa3NSb3RKZERUbWVMd0V6OVB1TndnUGJzd2k0RlBZcmEzQTBSczkueHZH"
request.add_header("Authorization", "Basic %s" % key)
data = json.load(urllib2.urlopen(request))
print(data)
```

#### curl을 사용한 restAPI 셋팅
- 아래 예제부터는 따로 코드를 작성하지 않고 curl 예제만 다룬다.

- Task에 대한 시작일 설정
```
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&startdate=0502" https://csi.lazypic.org/api/setstartdate
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&startdate=2018-05-02" https://csi.lazypic.org/api/setstartdate
# + 문자는 web에서 %2B 이다. 터미널에서 curl로 테스트시 + 대신 %2B를 넣어주어야 한다.
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&startdate=2018-05-02T14:45:34%2B09:00" https://csi.lazypic.org/api/setstartdate
```

- Task에 대한 1차 마감일 설정
```
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&predate=0605" https://csi.lazypic.org/api/setpredate
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&predate=2018-06-05" https://csi.lazypic.org/api/setpredate
# + 문자는 web에서 %2B 이다. 터미널에서 curl로 테스트시 + 대신 %2B를 넣어주어야 한다.
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&predate=2018-06-05T14:45:34%2B09:00" https://csi.lazypic.org/api/setpredate
```

