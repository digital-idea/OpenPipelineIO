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

#### Access Level 변경
웹을 통해서 가입되는 모든 유저의 AccessLevel 값은 아티스트를 상징하는 3 값이다.

아래 명령어를 통해서 AccessLevel 을 변경할 수 있다.
```bash
$ sudo csi3 -accesslevel 4 -id khw7096
```

- 0: 권한이 없는 레벨
- 1: 손님
- 2: 클라이언트
- 3: 아티스트
- 4: 팀장, 실무 매니저
- 5: 프로젝트 매니저
- 6: 슈퍼바이저
- 7: IO매니저, 서버권한이슈
- 8: PD, 자금이슈
- 9: 개발자
- 10: 시스템 관리자