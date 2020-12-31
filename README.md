# CSI

![travisCI](https://secure.travis-ci.org/digital-idea/csi3.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/digital-idea/csi3)](https://goreportcard.com/report/github.com/digital-idea/csi3)

![screenshot](figures/screenshot.png)

CSI(Creation Status Integrator)는 콘텐츠 제작을 위한 프로젝트 매니징 웹 어플리케이션(WAS) 입니다.

- 속도, 검색어 방식, 교육의 최소화, 단일파일 배포를 중점으로 개발되고 있습니다.
- 내부, 외부 서버에 설치가 가능합니다.
- 사용자별 토큰키, 암호화키, 직급별 접근권한 사용이 가능합니다.

### 다운로드
- [Linux 64bit](https://github.com/digital-idea/csi3/releases/download/v3.1.6/csi3_linux_x86-64.tgz)
- [macOS 64bit](https://github.com/digital-idea/csi3/releases/download/v3.1.6/csi3_darwin_x86-64.tgz)
- [Linux 64bit for Digitalidea](https://github.com/digital-idea/csi3/releases/download/v3.1.6/csi3_linux_di_x86-64.tgz): 회사가 필요한 인수가 자동으로 설정되어 있습니다.


### Roadmap
- 진행중: 전문적인 리뷰툴 추가(현재 웹프로토콜 + RV), with 슈퍼바이저,실장님
- [ ] 구조변경: 사용자별 샷 갯수, 리스트가 나오도록 처리 -> items 컬렉션을 통합하기
    - [ ] 통계툴을 내부에서 일괄 처리하도록 통합(기존 Statistics 서비스 통합)
- [ ] 파트너 관리툴 추가
- [ ] 소프트웨어 등록, 환경변수 관리 -> JWT토큰 방식의 웹프로토콜 연동하기.
- [ ] 웹 스캔 툴: ACES2065-1(또는 사용자 설정) > ACEScg
- [ ] 개인(아티스트) 정보 페이지
- [ ] 매니저(팀장) 정보 페이지
- [ ] 매니저(실장) 정보 페이지
- [ ] 장비관리 툴

### 데이터베이스(mongoDB) 설치, 실행

CentOS
- [CentOS7에서 mongoDB 설정](https://github.com/cgiseminar/curriculum/blob/master/docs/install_mongodb.md)

macOS
```bash
$ brew uninstall mongodb
$ brew tap mongodb/brew
$ brew install mongodb-community
$ brew services start mongodb-community
```

Windows
- https://fastdl.mongodb.org/win32/mongodb-win32-x86_64-2008plus-ssl-4.0.10-signed.msi
- Download: https://www.mongodb.com/download-center/community?jmp=docs
- Setup: https://docs.mongodb.com/manual/tutorial/install-mongodb-on-windows-unattended/

### CSI 실행하기
CSI를 실행하기 전에 우선적으로 thumbnail 폴더가 필요합니다.
이 경로는 CSI를 운용하면서 생성되는 썸네일 이미지, 사용자 프로필 사진이 저장되는 경로로 사용됩니다.
썸네일 경로가 이미 존재한다면, 해당경로를 서비스 시작시 `-thumbpath` 인수를 이용해서 설정할 수 있습니다.

```bash
$ mkdir thumbnail // 프로그램 시작전, thumbnail 경로가 없다면 생성해주세요.
$ sudo csi3 -http :80
```

> 여러분이 macOS를 사용한다면 기본적으로 80포트는 아파치 서버가 사용중일 수 있습니다. 80포트에 실행되는 아파치 서버를 종료하기 위해서 `$ sudo apachectl stop` 를 터미널에 입력해주세요.

CSI는 [wfs-웹파일시스템](https://github.com/digital-idea/wfs), [dilog-로그서버](https://github.com/digital-idea/dilog), [dilink-웹프로토콜](https://github.com/digital-idea/dilink)과 같이 연동됩니다. 아래 서비스 실행 및 프로토콜 설치도 같이 진행하면 더욱 강력한 CSI를 사용할 수 있습니다.

```bash
$ dilog -http :8080
$ wfs -http :8081
```

### CentOS 방화벽 설정
다른 컴퓨터에서 접근하기 위해서는 해당 포트를 방화벽 해제합니다.

```
# firewall-cmd --zone=public --add-port=80/tcp --permanent
success
# firewall-cmd --reload
```

### 터미널 명령어 / CommandLine
CSI는 터미널에서 간단하게 관리를 할 수 있습니다.
관리를 위해 필요한 명령어 메뉴얼입니다.

- [Project](documents/project.md)
- [Item](documents/item.md): Asset, Shot
- [User](documents/user.md)
- [Daily](documents/daily.md)
- [Organization](documents/organization.md)

### RestAPI
CSI는 RestAPI가 설계되어 있습니다.
Python, Go, Java, Javascript, node.JS, C++, C, C# 등 수많은 언어에서 CSI의 상태를 변경할 수 있습니다.

- [Project](documents/rest_project.md)
- [Item](documents/rest_item.md): Asset, Shot
- [User](documents/rest_user.md)
- [Organization](documents/rest_organization.md)
- [Tasksetting](documents/rest_tasksetting.md)
- [Status](documents/rest_status.md)
- [Review](documents/rest_review.md)

### 썸네일 경로
위에서 생성된 thumbnail 폴더는 아래 구조를 띄고 있습니다.
썸네일은 사내 다른 응용프로그램에서도 사용될 수 있기 때문에 경로구조를 표기해둡니다.

- 썸네일주소 : `thumbnail/{projectname}/{id}.jpg`
- 사용자이미지 : `thumbnail/user/{id}.jpg`

### 프로젝트 Process
- [디자인 프로세스](documents/process_designer.md)
- [개발 프로세스](documents/process_developer.md)
- [Onset Setellite](documents/setellite.md)
- [DB관리](documents/dbbackup.md)

### Developer
- CSI서버: https://csi.lazypic.org
- Log서버: http://csi.lazypic.org:8080
- WFS서버: http://csi.lazypic.org:8081
- 회사 전용 빌드문의: hello@lazypic.org
- Maintainer: Jason / jason@lazypic.org
- Committer: Bailey / bailey@lazypic.org
- Contributors:
- 체험계정 ID/PW: guest
    - guest 계정은 모든 메뉴가 보이지 않습니다.
    - Guest 계정은 일부 기능만 테스트 가능한 모드입니다.
    - 만약 많은 기능을 테스트하고 싶다면 가입한 ID와 함께 권한변경 요청메일을 hello@lazypic.org로 보내주세요.

### Infomation
- [CSI의 역사](documents/history.md)
- License: BSD 3-Clause License

### License
- CSI: BSD 3-Clause License
- [JScolor](http://jscolor.com/download/): GNU GPL license v3
- [Dropzone](https://www.dropzonejs.com): MIT License
- [JQuery](https://jquery.org/license/): MIT license
- [VFS](https://github.com/blang/vfs): MIT license
- [HttpFS](https://github.com/shurcooL/httpfs): MIT license
- [VFSgen](https://github.com/shurcooL/vfsgen): MIT license
- [Excelize](https://github.com/360EntSecGroup-Skylar/excelize): BSD 3-Clause License
- [Slack go webhook](https://github.com/ashwanthkumar/slack-go-webhook): Apache License, Version 2.0
- [Captcha](https://github.com/dchest/captcha): Apache License, Version 2.0
- [Mgo](https://github.com/go-mgo/mgo): https://github.com/go-mgo/mgo/blob/v2-unstable/LICENSE
- [JWT go](https://github.com/dgrijalva/jwt-go): MIT license
- [OpenColorIO](https://github.com/AcademySoftwareFoundation/OpenColorIO): BSD 3-Clause License
- [alfg/mp4](https://github.com/alfg/mp4): MIT license
- [amarburg/go-quicktime](https://github.com/amarburg/go-quicktime): MIT license
