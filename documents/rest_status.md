# Status RestAPI
Status RestAPI 입니다.

## Get
| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/status | status의 모든 정보를 가지고 온다. | `$ curl -X GET -H "Authorization: -H Basic {YourTokenKey}" https://csi.lazypic.org/api/status` |


# 파이썬 예제 Python2.7x
모든 Status 정보를 가지고 오기.

```python
#!/usr/bin/python
#coding:utf-8
import urllib
import urllib2
import json
try:
    request = urllib2.Request("https://csi.lazypic.org/api/status")
    key = "JDJhJDEwJHBBREluL0JuRTdNa3NSb3RKZERXbWVMd0V6OVB1TndnUGJzd2k0RlBZcmEzQTBSczkueHZH"
    request.add_header("Authorization", "Basic %s" % key)
    data = urllib.urlencode(values)
    request.add_data(data)
    result = urllib2.urlopen(request)
    data = json.load(result)
except:
    print("RestAPI에 연결할 수 없습니다.")
    # 이후 에러처리 할 것
print(data)
```