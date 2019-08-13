# RestAPI Project
restAPI의 장점은 URL을 이용해서 정보를 처리할 수 있습니다.
CSI 접속시 IP를 이용하면 VPN환경에서도 restAPI를 사용할 수 있습니다.

Python, Go, Java, C++, node.JS 언어를 이용해서 restAPI를 사용할 수 있습니다.

이 문서는 파이썬을 이용해서 RestAPI를 프로젝트정보를 가지고 오는 방법을 다룹니다.
파이프라인에 사용될 확률이 높은 코드라서, 에러처리까지 코드로 다루었습니다.

#### 프로젝트리스트 가지고오기
- 작업중인 프로젝트 리스트가지고오기

```
#coding:utf-8
import json
import urllib2

restURL = "http://10.0.90.251/api/projects" # 기본적으로 현재 작업중인 프로젝트를 가지고옵니다.(pre + post + backup상태)
restURL = "http://10.0.90.251/api/projects?status=pre" # Preproduction 상태를 가진 프로젝트를 가지고 옵니다.
try:
	projects = json.load(urllib2.urlopen(restURL))
except:
	print("RestAPI에 연결할 수 없습니다.")
	# 에러처리

if projects["error"]:
	print(projects["error"])
	# 에러처리
	# sys.exit(projects["error"]) #일반 Cmd 툴일때는 옆 처럼 에러처리를 해준다.
	# 그래픽스 툴이면 에러박스를 띄운다.
print(projects["data"])
```

##### 프로젝트 상태는 아래와 같습니다.
- unknown : 아무상태도 아님
- pre : 프리프로덕션단계(컨셉,프리비즈,계약단계)의 프로젝트
- post : 진행중인 프로젝트
- backup : 백업해야할 프로젝트
- lawsuit : 소송중인 프로젝트
- layover : 중단된 프로젝트
- archive : 백업완료된 프로젝트

#### 특정 프로젝트의 프로젝트정보(ProjectInfo)값을 가지고 오기
- 군함도 프로젝트정보중 EmailHead를 가지고오는 방법

```
#coding:utf-8
import json
import urllib2

restURL = "http://10.0.90.251/api/project?id=gunhamdo" # 군함도 프로젝트 자료구조를 가지고 옵니다.
try:
	data = json.load(urllib2.urlopen(restURL))
except:
	print("RestAPI에 연결할 수 없습니다.")
	# 에러처리

if "error" in data:
	print(data["error"])
	# 에러처리
print(data["mailhead"]) # 프로젝트의 MailHead를 가지고오는 방법
```

#### 프로젝트 담당자 가지고오기
```
#coding:utf-8
import json
import urllib2

restURL = "http://10.0.90.251/api/project?id=gunhamdo" # 군함도 프로젝트 자료구조를 가지고 옵니다.
try:
	data = json.load(urllib2.urlopen(restURL))
except:
	print("RestAPI에 연결할 수 없습니다.")
	# 에러처리

if "error" in data:
	print(data["error"])
	# 에러처리
print(data["super"])   # 슈퍼바이저
print(data["cgsuper"]) # CG슈퍼바이저
print(data["pd"])      # PD
print(data["pm"])      # PM
print(data["pa"])      # PA
```
