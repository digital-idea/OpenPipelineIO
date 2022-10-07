# RestAPI FullCalendar Resource

FullCalendar Resource Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/fcresource/{id}|FullCalendar Resource 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/fcresource/{id}"

## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/fcresource|FullCalendar Resource 정보를 추가합니다| Option 참고 |curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"title":"레이어이름","eventcolor":"#FF1100", "textcolor":"#FFFFFF"}' "https://openpipeline.io/api/fcresource"

- Option: https://github.com/digital-idea/OpenPipelineIO/blob/master/struct_fcevent.go

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/fcresource/{id}|FullCalendar Resource 정보를 수정합니다| 자료구조 참고 |curl -X PUT -H "Authorization: Basic {TOKEN}“ -d '{"title":"레이어이름","eventcolor":"#FF1100", "textcolor":"#FFFFFF"}' "https://openpipeline.io/api/fcresource/{id}"

## DELETE

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/fcresource/{id}|FullCalendar Resource 정보를 삭제합니다.|id|curl -X DELETE -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/fcresource/{id}"

## Option 체크

```bash
curl https://openpipeline.io/api/fcresource -v
```

```bash
HTTP/1.1 200 OK
< Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS
< Access-Control-Allow-Origin: *
< Date: Tue, 17 May 2022 02:10:41 GMT
< Content-Length: 0
```
