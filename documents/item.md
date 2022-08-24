# Item
샷, 에셋 관련 터미널 명령어 사용법 입니다.

#### 샷등록
DB값만 생성되며, 샷 폴더가 생성되지는 않습니다.

```bash
openpipelineio -add item -project [projectname] -name [SS_0010] -type [org]
openpipelineio -add item -project [projectname] -name [SS_0010] -type [org] -eposide e01 -season 1 # 샷등록시 에피소드와 시즌을 추가
openpipelineio -add item -project [projectname] -name [SS_0010] -type [org] -eposide e01 -season 1 -netflixid 123435 # 샷 등록시 에피소드, 시즌, 넷플릭스ID를 같이 추가
```

#### Plate 정보와 함께 샷 등록하는 인수
SS_0010 샷을 생성시 아래 인수를 이용해서 정보를 샷,에셋 생성시 함께 추가할 수 있습니다.

- project: 프로젝트명 예)circle
- name: 샷이름 예)SS_0010
- type: 타입 예)org
- platesize: 플레이트 사이즈 예)2048x1152
- scanname: 데이터 네임 예)A007C006_160424_R28L
- scantimecodein: 스캔데이터 타임코드 In 예)10:00:00:00
- scantimecodeout: 스캔데이터 타임코드 Out 예)10:00:04:04
- scanframe: 스캔 프레임수 예)100
- scanin: 스캔 In프레임, 보통 6자리를 가지고 있습니다. 예)A007C006_160424_R28L.456812.exr 문자에서 456812 숫자
- scanout: 스캔 Out프레임, 보통 6자리를 가지고 있습니다. 예)A007C006_160424_R28L.456915.exr 문자에서 456915 숫자
- platein: 플레이트 In프레임 예)1001
- plateout: 플레이트 Out프레임 예)1101
- justin: Just구간 In프레임 예)1003
- justout: Just구간 Out프레임 예)1098
- justtimecodein: Just구간 타임코드 In 예)10:00:00:03
- justtimecodeout: Just구간 타임코드 Out 예)10:00:04:02
- episode: 에피소드명
- season: 시즌명
- netflixid: 넷플릭스ID

```bash
openpipelineio -add item -project circle -name SS_0010 -type org -platesize 2048x1152 -scanname A007C006_160424_R28L -scantimecodein 10:00:00:00 -scantimecodeout 10:00:04:04 -scanframe 100 -scanin 456812 -scanout 456912 -platein 1001 -plateout 1101 -justin 1003 -justout 1098 -justtimecodein 10:00:00:03 -justtimecodeout 10:00:04:02
```

만약 justin, justout, justtimecodein, justtimecodeout 값이 비어있다면,
justin값은 platein, justout값은 plateout, justtimecodein값은 scantimecodein, justtimecodeout값은 scantimecodeout 값으로 자동등록 됩니다.

#### 샷,에셋 삭제
circle 프로젝트 SS_0010 샷 삭제

```bash
sudo openpipelineio -rm item -project circle -name SS_0010 -type org
```

circle 프로젝트 stone01 에셋 삭제

```bash
sudo openpipelineio -rm item -project circle -name stone01 -type asset
```

#### 에셋등록

에셋이 prop 타입이고 component 형태일 때

```bash
openpipelineio -add item -type asset -project [projectname] -name [Assetname] -assettype prop -assettags prop,component
```

#### 소스등록

```bash
openpipelineio -add item -name OPN_0010 -type src2 -project TEMP -platepath /source/path
```

등록이 되면 자동으로 OPN_0010 샷에 OPN_0010_org 소스 제목으로 /source/path 경로가 자동 등록됩니다.
