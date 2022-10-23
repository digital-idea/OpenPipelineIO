# 개발환경 셋팅
공동개발을 위한 개발환경 셋팅방법을 다루는 문서입니다.

## 개발에 필요한 소프트웨어
- Go 1.16.6 이상 버전을 사용해주세요.
- mongoDB 3.x 이상 버전을 사용해주세요.

## 에디터
[MS Visual Code](https://code.visualstudio.com)를 사용하며 툴내부 마켓플레이스에서 Go와 관련된 모든 편리한 기능을 설치, 사용하고 있습니다.
디버그, 실수방지, 특히 협업시에 코드의 일관성을 위해 사용하고 있습니다.

## 테스트서버
- https://openpipeline.io

## 로드벨런서
만약 웹서버에 많은 부하가 걸리게 되면 로드밸런서 셋팅을 추천합니다. 웹서버를 2개 이상 셋팅하고 사용해주세요.

- HAProxy: http://www.haproxy.org
- Reference
    - 설명: https://d2.naver.com/helloworld/284659
    - 이미지 서버 static 셋팅: http://www.mobileflow.co.kr/main/blog/220821876221
    - 알고리즘: https://chhanz.github.io/linux/2020/11/30/configuration-haproxy-rr/

## 인증서 관리

#### Letsencrypt
macOS에 letsencrypt를 설치, 인증서를 생성합니다.

```bash
$ brew install letsencrypt
$ sudo certbot certonly --standalone -d openpipeline.io
```

일반적으로 인증서가 생성되는 경로는 아래와 같습니다.

```bash
/etc/letsencrypt/live/openpipeline.io/fullchain.pem
/etc/letsencrypt/live/openpipeline.io/privkey.pem
```

#### https 서비스를 위한 인증서 갱신

```bash
$ sudo certbot renew
```

- reference: https://certbot.eff.org/lets-encrypt/osx-other

#### 인증서 경로를 수동으로 설정하여 서비스를 실행하는 방법

[Letsencrypt](https://letsencrypt.org)는 기존 인증서가 존재하는 상태에서 새로운 인증서를 재발급하면 `/etc/letsencrypt/live` 하위경로에 DNS 문자 뒤 숫자가 붙는 형태로 경로가 생성됩니다.
만약 재발급된 새로운 인증서를 사용하여 웹서비스를 운용하고 싶다면, OpenPipelineIO를 실행할 때 certfullchanin, certprivkey 옵션을 이용해서 fullchain.pem, privkey.pem 경로 설정 후 서비스를 운용할 수 있습니다.

```bash
$ sudo openpipelineio -certfullchanin /etc/letsencrypt/live/openpipeline.io-0002/fullchain.pem -certprivkey /etc/letsencrypt/live/openpipeline.io-0002/privkey.pem
```


#### https 서비스를 위한 자가 인증서 생성

자가 인증방식입니다. 접속시 에러가 있지만, https 보안프로토콜을 사용할 수 있습니다.

```bash
$ go run /usr/local/go/src/crypto/tls/generate_cert.go -host="openpipeline.io" -ca=true
```

## 배포

항상 바이너리 파일 하나를 지향합니다.
설치 서버 bin 폴더에 openpipelineio, protocol, dilog, wfs 파일을 배포합니다.

## 자동실행

#### CentOS

실행파일은 서버 /usr/local/OpenPipelineIO/bin/openpipelineio에 배포합니다.

OpenPipelineIO.service 파일을 작성하여 /etc/systemd/system 폴더에 넣어줍니다.

아래 명령어로 실행해 줍니다.

```bash
$ sudo systemctl start OpenPipelineIO.service
```

OpenPipelineIO.service 파일 내용

```
[Service]
User=linux_user
Group=linux_user_group
ExecStart=/usr/local/OpenPipelineIO/bin/openpipelineio
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=OpenPipelineIO

Restart=always
KillMode=process

[Install]
WantedBy=multi-user.target
```

#### macOS
~/Library/LaunchAgents 경로가 존재하는지 체크합니다.

서비스 관리 파일을 하나 생성합니다.

```bash
$ touch io.openpipeline.plist
```

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>KeepAlive</key>
    <true />
    <key>RunAtLoad</key>
    <true/>
    <key>Label</key>
    <string>io.openpipeline</string>
    <key>ProgramArguments</key>
    <array>
      <string>/usr/local/bin</string>
      <string>OpenPipelineIO</string>
      <string>-http=:443</string>
      <string>-dilog=http://openpipeline.io:8080</string>
      <string>-wfs=http://openpipeline.io:8081</string>
      <string>-authmode</string>
      <string>-signupaccesslevel=1</string>
    </array>
  </dict>
</plist>
```

- KeepAlive: 실행상태 유지
- RunAtLoad: Load시 실행되도록 하기

서비스 등록

```bash
$ launchctl load ~/Library/LaunchAgents/io.openpipeline.plist
```

서비스 시작

```bash
$ launchctl start io.openpipeline
```

서비스 종료

KeepAlive가 설정되어 있다면, 프로세스가 종료되더라도 완전히 자동으로 실행되도록 하고 싶지 않다면 unload후 plist 파일을 삭제합니다.

```bash
$ launchctl unload ~/Library/LaunchAgents/io.openpipeline.plist
$ rm ~/Library/LaunchAgents/io.openpipeline.plist
```