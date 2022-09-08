# RestAPI ProjectForPartner

ProjectForPartner Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/projectforpartner/{id}|파트너 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/projectforpartner/{id}"

## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/projectforpartner|새로운 파트너 정보를 추가합니다| Option 참고 |curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"name":"lazypic","phone":"821094117096", "deliverydates":[{"title1":"df","date":"2021-12-13"}]}' "https://openpipeline.io/api/projectforpartner"

- Option: https://github.com/digital-idea/OpenPipelineIO/blob/master/struct_projectforpartner.go

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/projectforpartner/{id}|기존 파트너 정보를 수정합니다|id, name, homepage, address, phone, email, timezone, descriptioon|curl -X PUT -H "Authorization: Basic {TOKEN}“ -d '{"name":"lazypic","phone":"821094117096"}' "https://openpipeline.io/api/projectforpartner/{id}"

## DELETE

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/projectforpartner/{id}|값을 삭제합니다.|id|curl -X DELETE -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/projectforpartner/{id}"

## Option 체크

```bash
curl https://openpipeline.io/api/projectforpartner -v
```

```bash
HTTP/1.1 200 OK
< Access-Control-Allow-Methods: GET,PUT,DELETE,OPTIONS,POST
< Access-Control-Allow-Origin: *
< Date: Tue, 17 May 2022 02:10:41 GMT
< Content-Length: 0
```
