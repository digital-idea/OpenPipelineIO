# Daily
OpenPipelineIO 에 등록된 미디어를 RV로 모아서 보는 기능과 관련된 터미널 명령어 사용법

#### 2016-12-05 날짜 아티스트가 업로드한 mov를 rvplayer로 한번에 모아보기.

```bash
$ openpipelineio -date 2016-12-05 -play &
```

#### 특정 프로젝트의 데일리 mov rvplayer로 모아보기.

```bash
$ openpipelineio -date 2016-12-05 -play -project [projectname] &
$ openpipelineio -date 2016-12-05 -play -project [projectname] -task model & // 해당 프로젝트의 model 테스크만 보기
```