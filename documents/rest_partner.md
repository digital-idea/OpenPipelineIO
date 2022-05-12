# RestAPI Partner

Partner Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/partner/|파트너 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.org/api/partner?id={id}"


## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/partner/|새로운 파트너 정보를 추가합니다| Option 참고 |curl -X POST -H 'Authorization: Basic <TOKEN>' -d '{"name":"lazypic","phone":"821094117096", "deliverydates":[{"title1":"df","date":"2021-12-13"}]}' "https://csi.lazypic.org/api/partner"

- Option: https://github.com/digital-idea/csi3/blob/master/partner.go

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/partner/|기존 파트너 정보를 수정합니다|id, name, homepage, address, phone, email, timezone, descriptioon|curl -X PUT -H "Authorization: Basic <TOKEN>“ -d '{"name":"lazypic","phone":"821094117096"}' "https://csi.lazypic.org/api/partner"

## DELETE

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/partner/|파트너 정보를 가져옵니다|id|curl -X DELETE -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.org/api/partner?id={id}"