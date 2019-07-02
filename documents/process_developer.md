# 개발환경 셋팅

CSI는 Go, mongoDB를 사용하여 진행되는 프로젝트입니다.
공동개발을 위한 개발환경 셋팅방법을 다루는 문서입니다.

### 개발환경셋팅
macOS

```bash
$ brew install go
$ brew install mongodb
```

#### 에디터
[MS Visual Code](https://code.visualstudio.com)를 사용하며 툴내부 마켓플레이스에서 Go와 관련된 모든 편리한 기능을 설치, 사용하고 있습니다.
디버그, 실수방지, 개발 속도를 많이 올릴 수 있으니 협업시에는 위 셋팅을 추천합니다.

### 테스트서버
- https://csi.lazypic.org
- 내부 임시 개발서버: 10.0.90.215

#### https 서비스를 위한 인증서 생성

```bash
$ go run /usr/local/go/src/crypto/tls/generate_cert.go -host="csi.lazypic.org" -ca=true
```

#### letsencrypt
macOS에 letsencrypt를 설치, 인증서를 생성합니다.

```bash
$ brew install letsencrypt
$ sudo certbot certonly --standalone -d csi.lazypic.org
```

인증서가 생성되는 경로는 아래와 같습니다.

```
/etc/letsencrypt/live/csi.lazypic.org/fullchain.pem
/etc/letsencrypt/live/csi.lazypic.org/privkey.pem
```

인증서 갱신

```bash
$ certbot renew
```

- reference: https://certbot.eff.org/lets-encrypt/osx-other

### TravisCI
테스트를 위해서 [TravisCI](https://docs.travis-ci.com) 를 사용합니다.

### 배포
항상 바이너리 파일 하나를 지향합니다.

디지털아이디어 아래 경로에 csi3, dilink, dilog, wfs 파일을 배포합니다.

/idea/app/csi