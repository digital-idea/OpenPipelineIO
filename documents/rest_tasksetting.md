# Tasksetting RestAPI
Tasksetting Restapi 입니다.

## Get
| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/shottasksetting | shot tasksetting 정보를 가지고 온다. | `$ curl http://csi.lazypic.org/api/shottasksetting` |
| /api/assettasksetting | asset tasksetting 정보를 가지고 온다. | `$ curl http://csi.lazypic.org/api/assettasksetting` |

## POST
| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/tasksetting | 인수를 입력받고 task에서 사용하는 경로를 반환한다. | project, name, task, type, assettype, os, seq, cut, userid | `$ curl -X POST -d "project=TEMP&seq=SS&cut=0010&task=comp" http://csi.lazypic.org/api/tasksetting` |
| /api/categorytasksettings | 카테고리를 입력받아 해당 task를 반환한다. | category | `$ curl -X POST -d "category=fx" http://csi.lazypic.org/api/categorytasksettings` |

# 파이썬 예제 Python2.7x
fx 카테고리를 가지고 있는 Task 가지고 오기

```python
#!/usr/bin/python
#coding:utf-8
import urllib
import urllib2
import json
try:
    request = urllib2.Request("http://192.168.31.172/api/categorytasksettings")
    key = "JDJhJDEwJHBBREluL0JuRTdNa3NSb3RKZERUbWVMd0V6OVB1TndnUGJzd2k0RlBZcmEzQTBSczkueHZH"
    request.add_header("Authorization", "Basic %s" % key)
    values = {"category":"fx"}
    data = urllib.urlencode(values)
    request.add_data(data)
    result = urllib2.urlopen(request)
    data = json.load(result)
except:
    print("RestAPI에 연결할 수 없습니다.")
    # 이후 에러처리 할 것
print(data)
```