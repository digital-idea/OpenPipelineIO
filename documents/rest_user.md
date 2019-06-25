# RestAPI User

id가 `woong`인 유저의 정보를 가지고 오기

```bash
$ curl http://csi.lazypic.org/api/user?id=woong
```

회사에서 팀장 사용자 데이터를 가지고 오는 방법
```bash
curl http://192.168.219.104/api/users?searchword="팀장"
```

합성팀, 1팀, 팀장 사용자 데이터를 가지고 오는 방법
```bash
curl http://192.168.219.104/api/users?searchword="합성팀,1팀,팀장"
```

개발팀,1팀 사용자 데이터를 가지고 오는 방법
```bash
curl http://192.168.219.104/api/users?searchword="개발팀,1팀"
```

#### 인증을 통한 restAPI 사용방법
curl 명령을 사용할 때 `-H "Authorization: Basic [TOKEN KEY]"` 문자를 사용합니다.

```bash
curl -H "Authorization: Basic JDJhJDEwJHBBREluL0JuRTdNa3NSb3RKZERUbWVMd0V6OVB1TndnUGJzd2k0RlBZcmEzQTBSczkueHZH" http://192.168.219.101/api/user?id=khw7096
```