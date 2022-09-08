# RestAPI Pipelinestep

Pipelinestep Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/pipelinestep/{id}| Pipelinestep 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/pipelinestep/{id}"

## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/pipelinestep|새로운 Pipelienstep 정보를 추가합니다| Option 참고 |curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"name":"fx","description":"FX팀"}' "https://openpipeline.io/api/pipelinestep"

- Option: https://github.com/digital-idea/OpenPipelineIO/blob/master/struct_pipelinestep.go

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/pipelinestep/{id}|기존 Pipelinestep 정보를 수정합니다|id, name, description|curl -X PUT -H "Authorization: Basic {TOKEN}“ -d '{"name":"fx","description":"FX팀Task"}' "https://openpipeline.io/api/pipelinestep/{id}"

## DELETE

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/pipelinestep/{id}|값을 삭제합니다.|id|curl -X DELETE -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/pipelinestep/{id}"

## Option 체크

```bash
curl https://openpipeline.io/api/pipelinestep -v
```

```bash
HTTP/1.1 200 OK
< Access-Control-Allow-Methods: GET,PUT,DELETE,OPTIONS,POST
< Date: Tue, 17 May 2022 02:10:41 GMT
< Content-Length: 0
```
