# RestAPI Step

Step Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/step/{id} | R&R 정보를 가져옵니다 | id | curl -X GET -H "Authorization: Basic {TOKEN}" "https://csi.lazypic.com/api/step/{id}"

## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/step | R&R 정보를 추가합니다 | Option 참고 | curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"order":1,"name":"견적받기","description":"견젹서를 받는 단계"}' "https://csi.lazypic.com/api/step"

- Option: https://github.com/digital-idea/csi3/blob/master/struct_step.go

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/step/{id}| R&R 정보를 수정합니다| 자료구조 참고 |curl -X PUT -H "Authorization: Basic {TOKEN}“ -d '{"order":1,"name":"견적받기","description":"견젹서를 받는 단계"}' "https://csi.lazypic.com/api/step/{id}"

## DELETE

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/step/{id}| R&R  값을 삭제합니다.|id|curl -X DELETE -H "Authorization: Basic {TOKEN}" "https://csi.lazypic.com/api/step/{id}"

## Option 체크

```bash
curl https://csi.lazypic.com/api/step -v
```

```bash
HTTP/1.1 200 OK
< Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS
< Access-Control-Allow-Origin: *
< Date: Tue, 17 May 2022 02:10:41 GMT
< Content-Length: 0
```
