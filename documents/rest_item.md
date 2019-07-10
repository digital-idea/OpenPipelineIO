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

## Post

| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/rmitem | 아이템 삭제 | project, name, type | `$ curl -d "project=circle&name=SS_0010&type=org" http://127.0.0.1/api/rmitem` |
| /api/settaskstatus | 상태수정 | project, name, task, status | `$ curl -d "project=circle&name=SS_0010&task=comp&status=wip" http://127.0.0.1/api/setstatus` |
| /api/settaskuser | 사용자수정 | project, name, task, user | `$ curl -d "project=TEMP&name=mamma&task=light&user=김한웅" http://192.168.219.104/api/settaskuser` |
| /api/settaskstartdate | 시작일 | project, name, task, date | `$ curl -d "project=TEMP&name=RR_0010&task=comp&date=0506" http://192.168.31.172/api/settaskstartdate` |
| /api/settaskpredate | 1차마감일 | project, name, task, date | `$ curl -d "project=TEMP&name=RR_0010&task=comp&date=0506" http://192.168.31.172/api/settaskpredate` |
| /api/settaskdate | 2차마감일 | project, name, task, date | `$ curl -d "project=TEMP&name=RR_0010&task=comp&date=0506" http://192.168.31.172/api/settaskdate` |
| /api/settaskmov | mov등록 | project, name, task, mov | `$ curl -d "project=TEMP&name=RR_0010&task=comp&mov=/show/test/test.mov" http://192.168.31.172/api/settaskmov` |
| /api/setshottype | shottype 변경 | project, name, type | `$ curl -d "project=TEMP&name=SS_0030&type=3d" http://192.168.0.11/api/setshottype` |
| /api/setthummov | 썸네일mov변경 | project, name, path | `$ curl -d "project=TEMP&name=SS_0030&path=/show/test.mov" http://192.168.0.11/api/setthummov` |


#### 샷,에셋정보(Item) 가지고오기. Python2.7x
- 군함도 S001_0001_org 샷 정보를 가지고 오기

```python
#coding:utf-8
import json
import urllib2

restURL = "http://10.0.90.251/api/item?project=gunhamdo&id=S001_0001_org" # CSI에서 군함도, S001_0001_org 샷 정보를 가지고 옵니다.
try:
	data = json.load(urllib2.urlopen(restURL))
except:
	print("RestAPI에 연결할 수 없습니다.")
	# 에러처리
if "error" in data:
	print(data["error"])
	# 에러처리
print(data)
```

#### 샷,에셋정보(Item) 가지고오기. Python3.7.x(블랜더)
```python
import requests
r = requests.get("http://172.30.1.50/api/item?project=TEMP&id=SS_0020_org")
print(r.json())
```

토큰키를 이용한 리퀘스트
```python
import requests
r = requests.get("http://172.30.1.50/api/item?project=TEMP&id=SS_0020_org", headers={'Authorization': 'Basic JDJhJDEwJHBBREluL0JuRTdNa3NSb3RKZERUbWVMd0V6OVB1TndnUGJzd2k0RlBZcmEzQTBSczkueHZH'})
print(r.json())
```

#### Shot 정보 가지고오기.
- 일반적으로 샷은 org타입이다.
- 입체의 경우 샷이 left 타입이다.
- 입체일 경우라도 컨버팅샷은 org 타입이다.
- 이 과정에서 입체상황을 고려하지 않고 자동으로 타입을 결정해서 샷 정보를 가지고 오고 싶을땐 아래 코드를 사용한다.

```
#coding:utf-8
import json
import urllib2

restURL = "http://10.0.90.251/api/shot?project=mkk3&name=JYW_0200" # mkk3에서 JYW_0200 샷 정보를 가지고 옵니다.
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
values["searchword"] = "comp+박지섭"
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

- 렌더 사이즈 셋팅
```
$ curl -X POST -d "project=TEMP&name=SS_0010&size=2048x1152" http://10.0.90.251/api/setrendersize
```

- 플레이트 사이즈 셋팅
```
$ curl -X POST -d "project=TEMP&name=SS_0010&size=2048x1152" http://10.0.90.251/api/setplatesize
```

- 렌즈디스토션 사이즈 셋팅
```
$ curl -X POST -d "project=TEMP&name=SS_0010&size=2048x1152" http://10.0.90.251/api/setdistortionsize
```

- 샷 프로젝션 사용여부 셋팅
```
$ curl -X POST -d "project=TEMP&name=SS_0010&projection=true" http://10.0.90.251/api/setcameraprojection
```

- 카메라 퍼블리쉬 테스크 : 사용가능한 테스크명은 "", "mm", "layout", "ani" 이다.
```
$ curl -X POST -d "project=TEMP&name=SS_0010&task=mm" http://10.0.90.251/api/setcamerapubtask
```

- 카메라 Publish 경로셋팅
```
$ curl -X POST -d "project=TEMP&name=SS_0010&path=/show/test/cam/pubpath" http://10.0.90.251/api/setcamerapubpath
```

- 썸네일 mov 등록
```
$ curl -X POST -d "project=TEMP&name=SS_0010&path=/show/thumb/nail.mov" http://10.0.90.251/api/setthummov
```

- Task 상태 변경
```
curl -X POST -d "project=TEMP&name=SS_0011&task=mm&status=omit" http://10.0.90.251/api/setstatus
```

- Just In/Out 값 설정
```
curl -X POST -d "project=TEMP&name=SS_0011&frame=1001" http://10.0.90.251/api/setjustin
curl -X POST -d "project=TEMP&name=SS_0011&frame=1030" http://10.0.90.251/api/setjustout
```

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

