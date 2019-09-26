# DB관리

### DB백업
- DB백업은 mongodump 명령을 사용한다. mongoDB를 설치하면 해당 명령어가 있습니다.
- 원칙상 백업은 서비스를 끄고 백업을 하는 것이 관례입니다. 현재 입력되는 값의 백업을 포기한다면, DB가 돌고 있을때도 mongodump는 사용가능합니다.
- 만약 /backuppath/20190613 경로에 CSI에서 사용하는 project DB를 백업한다면 명령어는 아래와 같습니다.

```bash
$ mongodump -d project -o /backuppath/20190906
$ mongodump -d user -o /backuppath/20190906
$ mongodump -d setellite -o /backuppath/20190906
$ mongodump -d projectinfo -o /backuppath/20190906
$ mongodump -d organization -o /backuppath/20190906
```

### DB복원
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

### DB복원 : bson 파일을 이용하기
- 아래는 projectnameadventure 프로젝트를 18일 날짜의 데이터로 DB 복구하는 예제이다.

```bash
$ mongorestore --drop -d project -c adventure /dbbackup/project/18/project/adventure.bson
$ mongorestore --drop -d projectinfo -c adventure /dbbackup/projectinfo/18/projectinfo/adventure.bson
```

### DB체크 : DB 부하체크하기
- http://127.0.0.1:28017

### DB체크 : DB에 들어간 값을 RestAPI를 이용해서 관찰하기
- db시작시 --rest 옵션을 붙히면 웹에서 json데이터를 관찰할 수 있다.

```
http://127.0.0.1:28017/project/projectname/
```

### 프로젝트 하나를 bson으로 백업하기
추후 분석을 위해서 프로젝트를 `.bson`으로 백업합니다.

```bash
$ mongodump --db project --collection TEMP --out /backupdb/TEMP
```