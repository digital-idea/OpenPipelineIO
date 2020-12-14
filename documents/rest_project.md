# RestAPI Project
restAPI의 장점은 URL을 이용해서 정보를 처리할 수 있습니다.
CSI 접속시 IP를 이용하면 VPN환경에서도 restAPI를 사용할 수 있습니다.

Python, Go, Java, C++, node.JS 언어를 이용해서 restAPI를 사용할 수 있습니다.

이 문서는 파이썬을 이용해서 RestAPI를 프로젝트정보를 가지고 오는 방법을 다룹니다.
파이프라인에 사용될 확률이 높은 코드라서, 에러처리까지 코드로 다루었습니다.

# RestAPI for Item

## Get

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/project | 프로젝트 정보를 가지고 옵니다. | project | `$ curl -H "Authorization: Basic <Token>" "https://csi.lazypic.org/api/project?id=TEMP"` |
| /api/projects | 프로젝트 상태를 입력하고 프로젝트 정보를 가지고 옵니다. | status | `$ curl -H "Authorization: Basic <Token>" "https://csi.lazypic.org/api/projects?status=post"` |
| /api/projecttags | 프로젝트에 사용중인 tags 가지고 옵니다. | project | `$ curl -H "Authorization: Basic <Token>" "https://csi.lazypic.org/api/projecttags?project=TEMP"` |
| /api/projectassettags | 프로젝트에 사용중인 asssettags 가지고 옵니다. | project | `$ curl -H "Authorization: Basic <Token>" "https://csi.lazypic.org/api/projectassettags?project=TEMP"` |

## Post

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/addproject | 프로젝트를 생성합니다. | id | `$ curl -X POST -d "id=TEMP" -H "Authorization: Basic <Token>" "https://csi.lazypic.org/api/addproject"` |

#### 프로젝트리스트 가지고오기
- 작업중인 프로젝트 리스트가지고오기

```python
#!/usr/bin/python
#coding:utf8
import urllib2
import json

endPoint = "http://10.0.90.251/api/projects" # 기본적으로 현재 작업중인 프로젝트를 가지고옵니다.(pre + post + backup상태)
# 특정상태의 프로젝트만 가지고 오고 싶다면 status 인수를 사용해주세요.
# endPoint = "http://10.0.90.251/api/projects?status=pre" # Preproduction 상태를 가진 프로젝트를 가지고 옵니다.
request = urllib2.Request(endPoint)
request.add_header("Authorization", "Basic <Token>")
result = urllib2.urlopen(request)
data = json.load(result)
print(data)
```

##### 프로젝트의 상태는 아래상태로 정의됩니다.
- test : 테스트 단계
- pre : 프리프로덕션 단계(컨셉,프리비즈,계약단계)의 프로젝트
- post : 진행중인 프로젝트
- backup : 백업해야할 프로젝트
- lawsuit : 소송중인 프로젝트
- layover : 중단된 프로젝트
- archive : 백업완료된 프로젝트


#### 프로젝트 정보 및 담당자 가지고오기
```python
#!/usr/bin/python
#coding:utf8
import urllib2
import json

request = urllib2.Request("https://csi.lazypic.org/api/project?id=TEMP") # TEMP 프로젝트 자료구조를 가지고 옵니다.
request.add_header("Authorization", "Basic <Token>")
result = urllib2.urlopen(request)
data = json.load(result)
print(data)
print(data["mailhead"]) # 프로젝트의 MailHead를 구하는 방법
print(data["super"])   # 슈퍼바이저
print(data["cgsuper"]) # CG슈퍼바이저
print(data["pd"])      # PD
print(data["pm"])      # PM
print(data["pa"])      # PA
```

#### 프로젝트에 사용중인 에셋태그 가지고오기
```python
#!/usr/bin/python
#coding:utf8
import urllib2
import json

request = urllib2.Request("https://csi.lazypic.org/api/projectassettags?project=TEMP")
request.add_header("Authorization", "Basic <Token>")
result = urllib2.urlopen(request)
data = json.load(result)
print(data)
```