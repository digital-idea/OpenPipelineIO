# CSI3

![travisCI](https://secure.travis-ci.org/digital-idea/csi3.png)

프로젝트 매니징 웹 파이프라인툴 입니다.

### 다운로드
- [Linux 64bit](https://github.com/digital-idea/csi3/releases/download/v1.0/csi3_linux_x86-64.tgz)
- [Linux 64bit Digitalidea](https://github.com/digital-idea/csi3/releases/download/v1.0/csi3_linux_di_x86-64.tgz): 회사가 필요한 인수가 자동으로 설정되어 있습니다.
- [macOS 64bit](https://github.com/digital-idea/csi3/releases/download/v1.0/csi3_darwin_x86-64.tgz)
- [macOS 64bit Digitalidea](https://github.com/digital-idea/csi3/releases/download/v1.0/csi3_darwin_di_x86-64.tgz): 회사가 필요한 인수가 자동으로 설정되어 있습니다.
- Windows 64bit: 필요시 진행합니다.
- Windows 64bit Digitalidea: 필요시 진행합니다.


### mongoDB 설치, 실행

```bash
$ sudo yum install mongodb mongodb-server // CentOS
$ brew install mongodb // macOS
```

mongodb를 실행합니다.
```
# service mongod start // CentOS
$ brew services start mongodb // macOS
```

- 참고: [CentOS7에서 mongoDB 설치](https://github.com/cgiseminar/curriculum/blob/master/docs/install_mongodb.md)

### CSI 실행
기본적으로 해당 리포지터리의 media, template를 사용합니다. 리포지터리로 이동하여 실행해주세요.
신뢰할 수 있는 네트워크망에서 사용하세요.

```
# csi3 -http :80 // CentOS
$ csi3 -http :80 // macOS
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

### 썸네일 경로
- 썸네일주소 : `/path/in/data/thumbnail/{project}/{slug}.jpg`
- 사용자이미지 : `/path/in/data/user/{ID}.jpg`

### DB관리

#### DB백업
- DB백업은 mongodump 명령을 사용한다. mongoDB를 설치하면 해당 명령어가 있습니다.
- 원칙상 백업은 서비스를 끄고 백업을 하는 것이 관례입니다. 현재 입력되는 값의 백업을 포기한다면, DB가 돌고 있을때도 mongodump는 사용가능합니다.
- 만약 /backuppath/20190613 경로에 CSI에서 사용하는 project DB를 백업한다면 명령어는 아래와 같습니다.
```
mongodump -d project -o /backuppath/20190613
```

#### DB복원
- CSI DB가 백업되어있는 경로로 이동한다.
- 백업경로 : /backuppath/project/20190613
- 터미널에서 아래 명령을 타이핑한다.
- 테스트서버의 DB를 삭제하고, 20190613 데이터로 복구하는 명령어는 아래와 같다.

```bash
$ mongorestore --drop -d project /backuppath/20190613/project
$ mongorestore --drop -d projectinfo /backuppath/20190613/projectinfo
$ mongorestore --drop -d user /backuppath/20190613/user
```

- 개발서버에서 모든 DB를 깨끗하게 비우고 테스트하고 싶다면 아래처럼 mongo쉘에서 명령어를 타이핑합니다.

```
> use project
> db.dropDatabase()

> use projectinfo
> db.dropDatabase()
```

- 아래처럼 터미널에서도 타이핑이 가능합니다.

```bash
$ mongo project --eval "db.dropDatabase()"
$ mongo projectinfo --eval "db.dropDatabase()"
```

#### DB복원 : bson 파일을 이용하기
- 아래는 projectnameadventure 프로젝트를 18일 날짜의 데이터로 DB 복구하는 예제이다.

```bash
$ mongorestore --drop -d project -c adventure /dbbackup/project/18/project/adventure.bson
$ mongorestore --drop -d projectinfo -c adventure /dbbackup/projectinfo/18/projectinfo/adventure.bson
```

#### DB체크 : DB 부하체크하기
- http://127.0.0.1:28017

#### DB체크 : DB에 들어간 값을 RestAPI를 이용해서 관찰하기
- db시작시 --rest 옵션을 붙히면 웹에서 json데이터를 관찰할 수 있다.

```
http://127.0.0.1:28017/project/projectname/
```

### Onset툴 Setellite 와 CSI의 연동
- Setellite는 데이터가 쌓일수록 사양을 많이 타게됩니다. 테스트시 Ipad Pro 2세대 이상을 사용할것.
- 홍콩에서 아이패드 구입하면 촬영시 소리가 나지 않습니다.
- 우리나라에서 판매되는 아이패드는 법규로 인해 촬영시 소리가 나도록 설정되어있습니다. 촬영장의 오디오팀에서 촬영시 소리에 민감할 수 있다.

### 디자인규칙
#### 아이콘
- SVG 포멧을 사용합니다.
- 벡터 파일이며 레이어들이 살아있어서 유지보수하기 좋습니다.
- 아스키 구조라서 git을 사용하기 좋습니다.

#### 디자인테스트방법
- CSI는 가변형 디자인을 추구합니다.
- 사파리는 개발자용 탭에서 모든 브라우저의 테스트가 가능합니다.
- 개발자용 > 응답형 디자인 모드 시작 옵션을 활성화하면 여러 장비테스트가 가능합니다.
- 크롬브라우저는 속도테스트, 디자인 세션별 디버깅에 좋습니다.
- 사파리에서 캐쉬를 비우고 새로고침하는 방법 : Opt+Cmd+E, Cmd+R

### Travis
- mongoDB setting: https://docs.travis-ci.com/user/database-setup/

### Process
- 디자인 프로세스
- 개발 프로세스

### Developer
- Maintainer: HanwoongKim(hello@lazypic.org)
- Committer: 
- Contributors: 

### Infomation
- [History](documents/history.md): csi의 역사
- License: BSD 3-Clause License
- 회사 전용 빌드문의: hello@lazypic.org
