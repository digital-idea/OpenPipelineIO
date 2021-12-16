# RestAPI Partner
Partner Restapi 입니다.

## Get
| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/partner/|파트너 정보를 가져옵니다|id|curl -H "Authorization: Basic <TOKEN>“ -X GET http://csi.lazypic.org/api/partner?id={id}

## POST
| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/partner/|새로운 파트너 정보를 추가합니다|name, homepage, address, phone, email. timezone, descriptioon|curl -H "Authorization: Basic <TOKEN>“ -X POST -d "name=test" http://csi.lazypic.org/api/partner

## PUT
| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/partner/|기존 파트너 정보를 수정합니다|id, name, homepage, address, phone, email. timezone, descriptioon|curl -H "Authorization: Basic <TOKEN>“ -X PUT -d "id=_id" http://csi.lazypic.org/api/partner