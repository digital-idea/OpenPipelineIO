# RestAPI Item
restAPI의 장점은 웹서비스의 URI를 이용하기 때문에 네트워크만 연결되어있으면 OS, 디바이스 제약없이 사용할 수 있습니다.
또한 python같은 언어를 이용해서 api를 제작하더라도 OS별로 코드가 달라지는 일이 없습니다.

이 문서는 기본 restAPI옵션, 파이썬을 이용해서 RestAPI를 사용하는 방법을 다룹니다.
파이프라인에 사용될 확률이 높은 코드라서, 에러처리까지 코드로 다루었습니다.

# RestAPI for Item

## Get

| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/item | 아이템 가지고 오기 | project, id | `$ curl -X GET "http://172.30.1.50/api/item?project=TEMP&id=SS_0020_org"` |
| /api/search | 검색 | project, searchword, sortkey | `$ curl -d "project=TEMP&searchword=SS_0020&sortkey=id" http://192.168.31.172/api/search` |
| /api/deadline2d | 2D마감일 리스트 | project | `$ curl -d "project=TEMP" http://192.168.31.172/api/deadline2d` |
| /api/deadline3d | 3D마감일 리스트 | project | `$ curl -d "project=TEMP" http://192.168.31.172/api/deadline3d` |

## Post

| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/rmitem | 아이템 삭제 | project, name | `$ curl -d "project=circle&name=SS_0010" http://127.0.0.1/api/rmitem` |
| /api/settaskstatus | 상태수정 | project, name, task, status | `$ curl -d "project=circle&name=SS_0010&task=comp&status=wip" http://127.0.0.1/api/setstatus` |
| /api/setassigntask | Assign 설정,해제 | project, name, task, status | `$ curl -d "project=TEMP&name=SS_0030&task=mg&status=true" http://192.168.31.172/api/setassigntask` |
| /api/settaskuser | 사용자수정 | project, name, task, user | `$ curl -d "project=TEMP&name=mamma&task=light&user=김한웅" http://192.168.219.104/api/settaskuser` |
| /api/settaskstartdate | 시작일 | project, name, task, date | `$ curl -d "project=TEMP&name=RR_0010&task=comp&date=0506" http://192.168.31.172/api/settaskstartdate` |
| /api/settaskpredate | 1차마감일 | project, name, task, date | `$ curl -d "project=TEMP&name=RR_0010&task=comp&date=0506" http://192.168.31.172/api/settaskpredate` |
| /api/settaskdate | 2차마감일 | project, name, task, date | `$ curl -d "project=TEMP&name=RR_0010&task=comp&date=0506" http://192.168.31.172/api/settaskdate` |
| /api/settaskmov | mov등록 | project, name, task, mov | `$ curl -d "project=TEMP&name=RR_0010&task=comp&mov=/show/test/test.mov" http://192.168.31.172/api/settaskmov` |
| /api/setshottype | shottype 변경 | project, name, type | `$ curl -d "project=TEMP&name=SS_0030&shottype=3d" http://192.168.0.11/api/setshottype` |
| /api/setthummov | 썸네일mov변경 | project, name, path | `$ curl -d "project=TEMP&name=SS_0030&path=/show/test.mov" http://192.168.0.11/api/setthummov` |
| /api/setassettype | assettype 변경 | project, name, type | `$ curl -d "project=TEMP&name=mamma&type=prop" http://192.168.0.11/api/setassettype` |
| /api/setoutputname | 아웃풋이름 등록 | project, name, outputname | `$ curl -d "project=TEMP&name=SS_0010&outputname=S101_010_010" http://192.168.31.172/api/setoutputname` |
| /api/setrnum | 롤넘버 등록 | project, name, rnum | `$ curl -d "project=TEMP&name=SS_0010&rnum=A0001" http://192.168.31.172/api/setrnum` |
| /api/setdeadline2d | 2D마감일 등록 | project, name, date | `$ curl -d "project=TEMP&name=SS_0010&date=0712" http://192.168.31.172/api/setdeadline2d` |
| /api/setdeadline3d | 3D마감일 등록 | project, name, date | `$ curl -d "project=TEMP&name=SS_0010&date=0712" http://192.168.31.172/api/setdeadline3d` |
| /api/setscantimecodein | 스캔 타임코드IN 등록 | project, name, timecode | `$ curl -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" http://192.168.31.172/api/setscantimecodein` |
| /api/setscantimecodeout | 스캔 타임코드OUT 등록 | project, name, timecode | `$ curl -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" http://192.168.31.172/api/setscantimecodeout` |
| /api/setscanin | scan in frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=654231" http://127.0.0.1/api/setscanin` |
| /api/setscanout | scan out frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=654331" http://127.0.0.1/api/setscanout` |
| /api/setscanframe | scan frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=100" http://127.0.0.1/api/setscanframe` |
| /api/setjusttimecodein | JUST 타임코드IN 등록 | project, name, timecode | `$ curl -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" http://192.168.31.172/api/setjusttimecodein` |
| /api/setjusttimecodeout | JUST 타임코드OUT 등록 | project, name, timecode | `$ curl -d "project=TEMP&name=SS_0010&timecode=01:00:01:21" http://192.168.31.172/api/setjusttimecodeout` |
| /api/setfinver | 최종데이터 버전 등록 | project, name, version | `$ curl -d "project=TEMP&name=SS_0010&version=1" http://192.168.31.172/api/setfinver` |
| /api/setfindate | 최종데이터 아웃풋 날짜 등록 | project, name, date | `$ curl -d "project=TEMP&name=SS_0010&date=0711" http://192.168.31.172/api/setfindate` |
| /api/setplatein | plate in frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1001" http://192.168.31.172/api/setplatein` |
| /api/setplateout | plate out frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1130" http://192.168.31.172/api/setplateout` |
| /api/setjustin | just in frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1003" http://192.168.31.172/api/setjustin` |
| /api/setjustout | just out frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1130" http://192.168.31.172/api/setjustout` |
| /api/sethandlein | handle in frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1003" http://192.168.31.172/api/sethandlein` |
| /api/sethandleout | handle out frame 등록 | project, name, frame | `$ curl -d "project=TEMP&name=SS_0010&frame=1130" http://192.168.31.172/api/sethandleout` |
| /api/addtag | tag 추가 | project, name, tag | `$ curl -d "project=TEMP&name=SS_0010&tag=테스트" http://192.168.219.104/api/addtag` |
| /api/rmtag | tags 삭제 | project, name, tag | `$ curl -d "project=TEMP&name=SS_0020&tag=태그3" http://192.168.219.114/api/rmtag` |
| /api/setnote | 작업내용 변경 | project, name, text, (userid) | `$ curl -d "project=TEMP&name=SS_0020&text=바람이 휘날린다" http://192.168.219.104/api/setnote` |
| /api/addcomment | 수정사항 추가 | project, name, text, (userid) | `$ curl -d "project=TEMP&name=SS_0020&text=1003프레임 나무제거" http://192.168.219.104/api/addcomment` |
| /api/rmcomment | 수정사항 삭제 | project, name, text, (userid) | `$ curl -d "project=TEMP&name=SS_0020&text=1003프레임 나무제거" http://192.168.219.104/api/rmcomment` |
| /api/addsource | 링크소스 추가 | project, name, title, path, (userid) | `$ curl -d "project=TEMP&name=SS_0020&title=source1&path=/show/src1/test.mov" http://192.168.31.172/api/addsource` |
| /api/rmsource | 링크소스 삭제 | project, name, title, (userid) | `$ curl -d "project=TEMP&name=SS_0020&title=sourcename" http://192.168.31.172/api/rmsource` |
| /api/setcameraprojection | 카메라 프로젝션여부 | project, name, projection | `$ curl -d "project=TEMP&name=SS_0020&projection=true" http://10.0.90.251/api/setcameraprojection` |
| /api/setcamerapubtask | 카메라 Pub Task설정 | project, name, task | `$ curl -d "project=TEMP&name=SS_0020&task=mm" http://10.0.90.251/api/setcamerapubtask` mm,layout,ani 만 task 등록가능|
| /api/setcamerapubpath | 카메라 Pub Path설정 | project, name, path | `$ curl -d "project=TEMP&name=SS_0020&path=/show/test/cam/pubpath" http://10.0.90.251/api/setcamerapubpath`|
| /api/setdeadline2d | 2D 마감일 설정 | project, name, date | `$ curl -d "project=TEMP&name=SS_0020&date=2019-09-05" http://10.0.90.251/api/setdeadline2d`|
| /api/setdeadline3d | 3D 마감일 설정 | project, name, date | `$ curl -d "project=TEMP&name=SS_0020&date=2019-09-05" http://10.0.90.251/api/setdeadline3d`|
| /api/setretimeplate | Retime Plate 경로설정 | project, name, path | `$ curl -d "project=TEMP&name=SS_0020&path=/show/retime" http://10.0.90.251/api/setretimeplate`|
| /api/settasklevel | Task 레벨설정 | project, name, task, level | `$ curl -d "project=TEMP&name=SS_0020&task=comp&level=1" http://10.0.90.251/api/settasklevel`|
| /api/setplatesize | Plate Size 설정 | project, name, size, (userid) | `$ curl -d "project=TEMP&name=SS_0020&size=2048x1152" http://10.0.90.251/api/setplatesize`|
| /api/setundistortionsize | Undistortion Size 설정 | project, name, size, (userid) | `$ curl -d "project=TEMP&name=SS_0020&size=2048x1152" http://10.0.90.251/api/setundistortionsize`|
| /api/rendersize | Reder Size 설정 | project, name, size, (userid) | `$ curl -d "project=TEMP&name=SS_0020&size=2048x1152" http://10.0.90.251/api/setrendersize`|


## Post(Legacy)
| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/addlink | 링크소스 추가 | project, name, text | `$ curl -d "project=TEMP&name=SS_0020&text=/show/src1/test.mov" http://192.168.31.172/api/addlink` |
| /api/rmlink | 링크소스 삭제 | project, name, text | `$ curl -d "project=TEMP&name=SS_0020&text=/show/src1/test.mov" http://192.168.31.172/api/rmlink` |
| /api/addnote | 작업내용 추가 | project, name, text | `$ curl -d "project=TEMP&name=SS_0020&text=바람이 휘날린다" http://192.168.219.104/api/addnote` |
| /api/rmnote | 작업내용 삭제 | project, name, text | `$ curl -d "project=TEMP&name=SS_0020&text=삭제내용" http://192.168.219.104/api/rmnote` |
| /api/setnotes | 현장,작업내용 교체 | project, name, text | `$ curl -d "project=TEMP&name=SS_0020&text=첫번째줄. 두번째줄.세번째줄." http://192.168.219.104/api/setnotes` |
| /api/setcomments | 수정사항 교체 | project, name, text | `$ curl -d "project=TEMP&name=SS_0020&text=첫번째줄. 두번째줄.세번째줄." http://192.168.31.172/api/setcomments` |
| /api/settags | tags 변경 | project, name, tags | `$ curl -d "project=TEMP&name=SS_0010&tags=태그1,태그2" http://192.168.219.114/api/settags` |
| /api/setlinks | 링크소스 교체 | project, name, text | `$ curl -d "project=TEMP&name=SS_0020&text=/show/src1,/show/src2" http://192.168.31.172/api/setlinks` |

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

restURL = "http://10.0.90.251/api2/items?project=gunhamdo&searchword=SS&wip=true"
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

url = "http://10.0.90.251/api2/items"
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
restURL = "http://10.0.90.251/api2/items?project=TEMP&searchword=comp+박지섭&wip=true"
```

- url endcode 방법

```
values = {}
values["project"] = "TEMP"
values["searchword"] = "comp+김한웅"
values["wip"] = "true"

url = "http://10.0.90.251/api2/items"
query = urllib.urlencode(values)
restURL = url + "?" + query
```

#### 샷,에셋에 대한 이름(Name) 검색하기
- 샷, 에셋에 대하여 문자열이 포함된 이름(Name)으로 검색할 수 있다.
- adventure 프로젝트의 "R0VFX" 이름을 가진 샷, 에셋 정보를  검색하는 예제이다.

```
- http://10.0.90.251/api/searchname?project=adventure&name=R0VFX
- http://10.0.90.251/api/searchname?project=adventure&name=R0VFX_sh001
```
```
#coding:utf-8
import json
import urllib2

restURL = "http://10.0.90.251/api/searchname?project=adventure&name=R0VFX_sh001"
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
- http://10.0.90.251/api/seqs?project=mkk3
```
```
#coding:utf-8
import json
import urllib2

restURL = "http://10.0.90.251/api/seqs?project=mkk3"
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
- http://10.0.90.251/api/shots?project=mkk3&seq=BNS
- http://10.0.90.251/api/shots?project=mkk3&seq=BNS_00
```
```
#coding:utf-8
import json
import urllib2

restURL = "http://10.0.90.251/api/shots?project=mkk3&seq=BNS"
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
$ curl -d "project=circle&name=SS_0010&task=light&mov=/show/test.mov" http://127.0.0.1/api/setmov
```

circle 프로젝트 mamma 에셋 fur 테스크에 /show/fur.mov 등록하기.

```bash
$ curl -d "project=circle&name=mamma&task=fur&mov=/show/fur.mov" http://127.0.0.1/api/setmov
```

#### python에서 샷 mov등록하기
- 파이썬은 curl 대신 urllib2.request를 이용해서 post 할 수 있다.
- 아래는 파이썬에서 urllib2.request를 이용하여 mov를 등록하는 예제이다.

```
def Setmov(project, shot, task, mov):
    """
	CSI에 mov를 등록한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&shot=&s&task=%s&mov=%s" % (project, shot, task, mov)
	try:
		request = urllib2.Request("http://10.0.90.251/api/setmov", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err
```

- csi3.py에 API로 활용할 수 있도록 구현이 되어있다.

```
import csi3
csi3.Setmov("TEMP", "SS_0010", "mm", "SS_0010.mov")
```

#### curl을 사용한 restAPI 셋팅
- 아래 예제부터는 따로 코드를 작성하지 않고 curl 예제만 다룬다.

- Task에 대한 시작일 설정
```
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&startdate=0502" http://10.0.90.251/api/setstartdate
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&startdate=2018-05-02" http://10.0.90.251/api/setstartdate
# + 문자는 web에서 %2B 이다. 터미널에서 curl로 테스트시 + 대신 %2B를 넣어주어야 한다.
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&startdate=2018-05-02T14:45:34%2B09:00" http://10.0.90.251/api/setstartdate
```

- Task에 대한 1차 마감일 설정
```
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&predate=0605" http://10.0.90.251/api/setpredate
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&predate=2018-06-05" http://10.0.90.251/api/setpredate
# + 문자는 web에서 %2B 이다. 터미널에서 curl로 테스트시 + 대신 %2B를 넣어주어야 한다.
curl -X POST -d "project=TEMP&name=SS_0011&task=fx&predate=2018-06-05T14:45:34%2B09:00" http://10.0.90.251/api/setpredate
```

