# Status RestAPI

Status 정보 관련 RestAPI 입니다.

## Get

| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/status | status의 모든 정보를 가지고 온다. | . | `$ curl -X GET -H "Authorization: Basic {TOKEN}" https://openpipeline.io/api/status` |
| /api/statusinfo | legacy status의 상태 정보를 가지고 온다. | reverse | `$ curl -X GET https://openpipeline.io/api/statusinfo?reverse=true` |


```bash
curl "https://openpipeline.io/api/statusinfo" -v # 옵션확인
curl "https://openpipeline.io/api/statusinfo?reverse=false" | json_pp -json_opt pretty,canonical # json pretty on macOS
```

## POST

| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/addstatus | status를 추가한다. | id,description,textcolor,bgcolor,bordercolor,order,defaulton,initstatus |`$ curl -X POST -H "Authorization: Basic {TOKEN}" -d "id=ready&description=ready&textcolor=#000000&bgcolor=#BEEF37&bordercolor=#BEEF37&order=3&defaulton=true&initstatus=false" "https://openpipeline.io/api/addstatus"` |

## 파이썬 예제 Python2.7x

모든 Status 정보를 가지고 오기.

```python
#!/usr/bin/python
#coding:utf-8
import urllib
import urllib2
import json
try:
    request = urllib2.Request("https://openpipeline.io/api/status")
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
