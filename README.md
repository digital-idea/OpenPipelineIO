# CSI3

![travisCI](https://secure.travis-ci.org/digital-idea/csi3.png)

![screenshot](figures/screenshot.png)

콘텐츠 제작을 위한 프로젝트 매니징 웹 어플리케이션 입니다.

- 속도, 간결함(검색어 Based), 교육의 최소화, 단일파일 배포를 중점으로 개발되고 있습니다.
- 내부, 외부 서버 모두 설치가 가능합니다.
- 사용자별 세션키, 암호화키를 사용합니다. 직책별 페이지 접근레벨을 설정할 수 있습니다.

### 다운로드
- [Linux 64bit](https://github.com/digital-idea/csi3/releases/download/v3.0.6/csi3_linux_x86-64.tgz)
- [Linux 64bit for Digitalidea](https://github.com/digital-idea/csi3/releases/download/v3.0.6/csi3_linux_di_x86-64.tgz): 회사가 필요한 인수가 자동으로 설정되어 있습니다.
- [macOS 64bit](https://github.com/digital-idea/csi3/releases/download/v3.0.6/csi3_darwin_x86-64.tgz)
- [macOS 64bit for Digitalidea](https://github.com/digital-idea/csi3/releases/download/v3.0.6/csi3_darwin_di_x86-64.tgz): 회사가 필요한 인수가 자동으로 설정되어 있습니다.
- [Windows 64bit](https://github.com/digital-idea/csi3/releases/download/v3.0.6/csi3_windows_x86-64.tgz)
- [Windows 64bit for Digitalidea](https://github.com/digital-idea/csi3/releases/download/v3.0.6/csi3_windows_di_x86-64.tgz): 회사가 필요한 인수가 자동으로 설정되어 있습니다.

### Roadmap
- [ ] PM입력기 웹 입력기로 대체
- [ ] Multi Task 기능추가. (참고: CSI의 내부구조가 완전히 달라진다. 유저 자료구조 안정화이후 진행 시작)
- [ ] 3D 파이프라인에 추가적으로 필요한 자료구조 및 API 생성(샷트레킹, 버전)
- [ ] 통계툴을 내부에서 처리하도록 변경.(기존 Statistics 페이지 통합)
- [ ] 현 작업내용,현장정보,수정사항이 "날짜;툴;작성자;내용" 형태인것을 맵으로 바꾸기(현재는 CSIv1 구조)
- [ ] 사용자별 샷 갯수, 리스트가 나오도록 처리하기.(멀티태스크이후)

### 데이터베이스(mongoDB) 설치, 실행

CentOS
```bash
$ sudo yum install mongodb mongodb-server
$ sudo service mongod start
```
- [CentOS7에서 mongoDB 설정](https://github.com/cgiseminar/curriculum/blob/master/docs/install_mongodb.md)

macOS
```bash
$ brew install mongodb
$ brew services start mongodb
```

Windows
- https://fastdl.mongodb.org/win32/mongodb-win32-x86_64-2008plus-ssl-4.0.10-signed.msi
- Download: https://www.mongodb.com/download-center/community?jmp=docs
- Setup: https://docs.mongodb.com/manual/tutorial/install-mongodb-on-windows-unattended/

### CSI 실행하기
CSI를 실행하기 전에 우선적으로 thumbnail 폴더가 필요합니다.
이 경로는 CSI를 운용하면서 생성되는 썸네일 이미지, 사용자 프로필 사진이 저장되는 경로로 사용됩니다.

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
- AWS 운용서버: https://csi.lazypic.org
- 회사 전용 빌드문의: hello@lazypic.org
- Maintainer: Jason / jason@lazypic.org
- Committer: Bailey / bailey@lazypic.org
- Contributors:

### Infomation
- [CSI의 역사](documents/history.md)
- License: BSD 3-Clause License