# DB관리

## DB백업,복원 및 DB 마이그레이션

DB 백업하기

```bash
mongodump -o /dbdump/path/20220405
```

DB 복원하기

```bash
mongorestore --drop /dbdump/path/20220405
```

## mongomirror

- mongomirror: https://www.mongodb.com/docs/atlas/import/mongomirror/

```bash
mongomirror --host "MySourceRS/host1.example.net:27017,host2.example.net:27017,host3.example.net:27017" \
   --ssl \
   --username "mySourceUser" \
   --password "mySourceP@$$word" \
   --authenticationDatabase "admin" \
   --destination "myAtlasRS/00.foo.mongodb.net:27017,01.foo.mongodb.net:27017,02.foo.mongodb.net:27017" \
   --destinationUsername "myAtlasAdminUser" \
   --destinationPassword "atlasPassword"
```

## DB체크 : DB 부하체크하기

- http://127.0.0.1:28017

## DB체크 : DB에 들어간 값을 RestAPI를 이용해서 관찰하기

- db시작시 --rest 옵션을 붙히면 웹에서 json 데이터를 관찰할 수 있습니다.

```text
http://127.0.0.1:28017/project/projectname/
```

## 프로젝트 하나를 bson으로 백업하기

추후 분석을 위해서 프로젝트를 `.bson`으로 백업합니다.

```bash
mongodump --db project --collection TEMP --out /backupdb/TEMP
```

## AWS 백업

DB는 클라우드에 백업하는 것을 추천합니다. 먼저 awscli를 설치합니다.

```bash
pip3 install awscli --upgrade --user
```

crontab에 백업스크립트를 등록합니다.
백업이 되지 않을 때를 대비해 확인할 log처리도 해둡니다.

```bash
$ crontab -e
* 3 * * * sh ~/OpenPipelineIO/script/backup_aws.sh >> /tmp/cron.out 2>&1
```
