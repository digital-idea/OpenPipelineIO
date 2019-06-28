# RestAPI User

id가 `woong`인 유저의 정보를 가지고 오기

```bash
$ curl http://csi.lazypic.org/api/user?id=woong
```

회사에서 팀장 사용자 데이터를 가지고 오는 방법
```bash
curl http://192.168.219.104/api/users?searchword=팀장
```

합성팀, 1팀, 팀장 사용자 데이터를 가지고 오는 방법
```bash
curl http://192.168.219.104/api/users?searchword=합성팀,1팀,팀장
```

개발팀,1팀 사용자 데이터를 가지고 오는 방법
```bash
curl http://192.168.219.104/api/users?searchword=개발팀,1팀
```

#### 인증을 통한 restAPI 사용방법

`-authmode` 모드가 활성화 되어있지 않다면 restAPI은 공개되어 있고, 아무나 사용할 수 있습니다.
csi서비스를 실행할 때 `-authmode` 옵션을 붙히면 보안모드가 작동됩니다.
이 때부터는 curl 명령을 사용할 때 `-H "Authorization: Basic [TOKEN KEY]"` 문자를 사용해야 데이터를 가지고 올 수 있습니다.
이 Token 키는 사용자의 가입, 패스워드변경, 패스워드초기화시 자동으로 변경됩니다.

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

#### JWT에 사용되는 환경변수
CSI_JWT_TOKEN_KEY 로 세션 암호화에 사용될 문자를 환경변수로 잡아주세요.
서버로 사용될 컴퓨터에서 아래 파일을 편집하면 됩니다.
보안에 문제가 될 이슈가 있다면 가끔 주기적으로 바꾸어주세요. session 암호화에 사용되기 때문에
사용자는 로그인만 다시 해주면 됩니다.

macOS라면 ~/.profile, centOS 라면 ~/.bashrc 파일 입니다.

```bash
export CSI_JWT_TOKEN_KEY="암호화,복호화에 사용될 문자"
```