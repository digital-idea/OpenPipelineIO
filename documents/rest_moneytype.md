# RestAPI Moneytype

Moneytype Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/moneytype/{id}|자금타입 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic {TOKEN}" "https://csi.lazypic.com/api/moneytype/{id}"

## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/moneytype|자금타입 정보를 추가합니다| Option 참고 |curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"order":1,"name":"firstestimate", "description":"최초견적"}' "https://csi.lazypic.com/api/moneytype"

- Option: https://github.com/digital-idea/csi3/blob/master/struct_moneytype.go

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/moneytype/{id}|자금타입 정보를 수정합니다| 자료구조 참고 |curl -X PUT -H "Authorization: Basic {TOKEN}“ -d '{"order":1,"name":"firstestimate", "description":"최초견적"}' "https://csi.lazypic.com/api/moneytype/{id}"

## DELETE

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/moneytype/{id}| 자금타입 값을 삭제합니다.|id|curl -X DELETE -H "Authorization: Basic {TOKEN}" "https://csi.lazypic.com/api/moneytype/{id}"

## Option 체크

```bash
curl https://csi.lazypic.com/api/moneytype -v
```

```bash
HTTP/1.1 200 OK
< Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS
< Access-Control-Allow-Origin: *
< Date: Tue, 17 May 2022 02:10:41 GMT
< Content-Length: 0
```
