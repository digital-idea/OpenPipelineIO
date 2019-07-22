# CSI3

![travisCI](https://secure.travis-ci.org/digital-idea/csi3.png)

![screenshot](figures/screenshot.png)

콘텐츠 제작을 위한 프로젝트 매니징 웹 어플리케이션 입니다.

- 속도, 간결함(검색어 Based), 교육의 최소화, 단일파일 배포를 중점으로 개발되고 있습니다.
- 내부, 외부 서버 모두 설치가 가능합니다.
- 사용자별 세션키, 암호화키를 사용합니다. 직책별 페이지 접근레벨을 설정할 수 있습니다.

### 다운로드

- [Linux 64bit](https://github.com/digital-idea/csi3/releases/download/v3.0.3/csi3_linux_x86-64.tgz)
- [Linux 64bit for Digitalidea](https://github.com/digital-idea/csi3/releases/download/v3.0.3/csi3_linux_di_x86-64.tgz): 회사가 필요한 인수가 자동으로 설정되어 있습니다.
- [macOS 64bit](https://github.com/digital-idea/csi3/releases/download/v3.0.3/csi3_darwin_x86-64.tgz)
- [macOS 64bit for Digitalidea](https://github.com/digital-idea/csi3/releases/download/v3.0.3/csi3_darwin_di_x86-64.tgz): 회사가 필요한 인수가 자동으로 설정되어 있습니다.
- [Windows 64bit](https://github.com/digital-idea/csi3/releases/download/v3.0.3/csi3_windows_x86-64.tgz)
- [Windows 64bit for Digitalidea](https://github.com/digital-idea/csi3/releases/download/v3.0.3/csi3_windows_di_x86-64.tgz): 회사가 필요한 인수가 자동으로 설정되어 있습니다.


### Roadmap
- [x] 빌드환경 구축 ![build](http://progressed.io/bar/100)
- [x] 유저자료구조 추가 ![user](http://progressed.io/bar/100)
- [x] 과거에 사용된 Python CSI2 API를 웹 restAPI로 전환(총46개함수) ![restAPI](http://progressed.io/bar/100)
- [ ] 리서치: 웹용 PM 입력기(MultiTask로 변경되면 기존 X,Y축 입력방식이 아닌 웹입력기 제작이 필요 ![input](http://progressed.io/bar/5)
    - [x] 태그입력 기술테스트 완료.
    - [ ] 작업배정 기술테스트
    - [ ] 시간정보 수정
    - [ ] 작업내용 개별입력
    - [ ] 작업내용 일괄입력
    - [ ] 수정내용 개별입력
    - [ ] 수정내용 일괄입력
- [ ] Multi Task 기능추가. (참고: CSI의 내부구조가 완전히 달라진다.)
- [ ] 3D 파이프라인에 추가적으로 필요한 자료구조 및 API 생성(샷트레킹, 버전)
- [ ] 통계툴을 내부에서 처리하도록 변경.(기존 Statistics 페이지 통합)
- [ ] 현 작업내용,현장정보,수정사항이 "날짜;툴;작성자;내용" 형태인것을 맵으로 바꾸기(현재는 CSIv1 구조)
- [ ] 사용자별 샷 갯수, 리스트가 나오도록 처리하기.(멀티태스크이후)

### mongoDB 설치, 실행

CentOS
```bash
$ sudo yum install mongodb mongodb-server
$ sudo service mongod start
```

macOS
```bash
$ brew install mongodb
$ brew services start mongodb
```

Windows
- https://fastdl.mongodb.org/win32/mongodb-win32-x86_64-2008plus-ssl-4.0.10-signed.msi
- Download: https://www.mongodb.com/download-center/community?jmp=docs
- Setup: https://docs.mongodb.com/manual/tutorial/install-mongodb-on-windows-unattended/

### CSI 실행

```bash
$ sudo csi3 -http :80
```

> macOS는 아파치 서버가 Built-in 되어있습니다. 80포트를 사용하기 위해서는 `$ sudo apachectl stop` 명령어를 입력해주세요.

csi는 [wfs](https://github.com/digital-idea/wfs), [dilog](https://github.com/digital-idea/dilog), [dilink](https://github.com/digital-idea/dilink)과 연동됩니다. 아래 서비스도 같이 실행해주세요.

```bash
$ dilog -http :8080
$ wfs -http :8081
```

### 명령어 사용법
- 프로젝트 생성
```bash
$ csi3 -add project -name [projectname]
```

- 프로젝트 삭제: Root 계정에서만 작동됩니다.
```
# csi3 -rm project -name [projectname]
```

- 샷등록
DB값만 생성되며, 샷 폴더가 생성되지는 않습니다.

```bash
$ csi3 -add item -project [projectname] -name [SS_0010] -type [org]
```

- 샷,에셋 삭제
```bash
$ csi3 -rm item -project [projectname] -name [SS_0010] -type [org]
```

- 에셋등록 예: 에셋이 prop 타입이고 component 형태일 때
```
$ csi3 -add item -type asset -project [projectname] -name [Assetname] -assettype prop -assettags prop,component
```

- 2016-12-05 에 아티스트가 업로드한 mov를 rvplayer로 한번에 모아보기.
```bash
$ csi3 -date 2016-12-05 -play &
```

- 특정 프로젝트의 데일리 mov rvplayer로 모아보기.
```bash
$ csi3 -date 2016-12-05 -play -project [projectname] &
$ csi3 -date 2016-12-05 -play -project [projectname] -task model & // 해당 프로젝트의 model 테스크만 보기
```

- 사용자 패스워드 초기화
사용자의 패스워드를 `Welcome2csi!`로 초기화 하고 싶다면 아래 명령어를 사용합니다.
관리자만 처리할 수 있습니다.

```bash
$ sudo csi3 -initpass Welcome2csi! -id [userid]
```

- 사용자 제거

```bash
$ sudo csi3 -rm user -id [userid]
```

### 썸네일 경로
- 썸네일주소 : `/thumbnail/{projectname}/{slug}.jpg`
- 사용자이미지 : `/thumbnail/user/{id}.jpg`

### RestAPI
- [Project](documents/rest_project.md)
- [Item](documents/rest_item.md)
- [User](documents/rest_user.md)

### Process
- [디자인 프로세스](documents/process_designer.md)
- [개발 프로세스](documents/process_developer.md)
- [Onset Setellite](documents/setellite.md)
- [DB관리](documents/dbbackup.md)

### Developer
- Maintainer: HanwoongKim(hello@lazypic.org)
- Committer: 
- Contributors: 

### Infomation
- [History](documents/history.md): csi의 역사
- License: BSD 3-Clause License
- 회사 전용 빌드문의: hello@lazypic.org
- 참고: [CentOS7에서 mongoDB 설치](https://github.com/cgiseminar/curriculum/blob/master/docs/install_mongodb.md)
