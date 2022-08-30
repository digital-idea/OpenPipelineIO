# RestAPI Money

Money Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/money/{id}|자금 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/money/{id}"

## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/money|자금 정보를 추가합니다| Option 참고 |curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"moneytypeid":"firstestimate","project":"test", "sender":"lazypic","recipient":"client","amount":30000000,"date":"2022-07-07","monetaryunit":"KRW","issuanceelectronictaxinvoice":false}' "https://openpipeline.io/api/money"

- Option: https://github.com/digital-idea/csi3/blob/master/struct_money.go

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/money/{id}|자금 정보를 수정합니다| 자료구조 참고 |curl -X PUT -H "Authorization: Basic {TOKEN}“ -d '{"moneytypeid":"firstestimate","project":"test", "sender":"lazypic","recipient":"client","amount":"30000000","date":"2022-07-07","monestaryunit":"KRW","issuanceelectronictaxinvoice":false}' "https://openpipeline.io/api/money/{id}"

## DELETE

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/money/{id}| 자금 값을 삭제합니다.|id|curl -X DELETE -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/money/{id}"

## Option 체크

```bash
curl https://openpipeline.io/api/money -v
```

```bash
HTTP/1.1 200 OK
< Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS
< Access-Control-Allow-Origin: *
< Date: Tue, 17 May 2022 02:10:41 GMT
< Content-Length: 0
```
