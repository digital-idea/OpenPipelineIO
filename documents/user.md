# 사용자

사용자 관련 명령어

#### 사용자 패스워드 초기화
사용자의 패스워드를 `Welcome2csi!`로 초기화 하고 싶다면 아래 명령어를 사용합니다.
관리자만 처리할 수 있습니다.

```bash
$ sudo csi3 -initpass Welcome2csi! -id [userid]
```

#### 사용자 제거

```bash
$ sudo csi3 -rm user -id [userid]
```

#### 서비스 시작시 사용자 권한 설정
csi3를 실행하면 기본적으로 accesslevel 3 으로 가입됩니다.
만약 웹서비스 시작시 accesslevel 0으로 서비스를 시작하고 싶다면, -initaccesslevel 옵션을 이용해서 웹서버를 실행할 수 있습니다.

```bash
$ sudo csi3 -initaccesslevel 0 -http :80
```