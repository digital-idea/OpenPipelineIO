# Item
샷, 에셋 관련 터미널 명령어 사용법 입니다.

#### 샷등록
DB값만 생성되며, 샷 폴더가 생성되지는 않습니다.

```bash
$ csi3 -add item -project [projectname] -name [SS_0010] -type [org]
```

#### Plate 정보와 함께 샷 등록
SS_0010 샷을 생성시 아래 정보를 함께 추가할 수 있습니다.

- Project(프로젝트): circle
- Name(샷이름): SS_0010
- Type(타입): org
- Platesize(플레이트 사이즈): 2048x1152
- Scanname(데이터 네임): A007C006_160424_R28L
- ScanTimecodeIn(스캔데이터 타임코드 In): 10:00:00:00
- ScanTimecodeOut(스캔데이터 타임코드 Out): 10:00:04:04
- ScanFrame(스캔 프레임수): 100
- ScanIn(스캔 In 프레임, 보통 6자리를 가지고 있습니다 A007C006_160424_R28L.######.exr): 456812
- ScanOut(스캔 Out 프레임, 보통 6자리를 가지고 있습니다 A007C006_160424_R28L.######.exr): 456912
- PlateIn(플레이트 In 프레임): 1001
- PlateOut(플레이트 Out 프레임): 1101
- JustIn(Just구간 In 프레임): 1003
- JustOut(Just구간 Out 프레임): 1098
- JustTimecodeIn(Just구간 타임코드 In): 10:00:00:03
- JustTimecodeOut(Just구간 타임코드 Out): 10:00:04:02

```bash
$ csi3 -add item -project circle -name SS_0010 -type org -platesize 2048x1152 -scanname A007C006_160424_R28L -scantimecodein 10:00:00:00 -scantimecodeout 10:00:04:04 -scanframe 100 -scanin 456812 -scanout 456912 -platein 1001 -plateout 1101 -justin 1003 -justout 1098 -justtimecodein 10:00:00:03 -justtimecodeout 10:00:04:02
```

만약 justin, justout, justtimecodein, justtimecodeout 값이 비어있다면,
justin값은 platein, justout값은 plateout, justtimecodein값은 scantimecodein, justtimecodeout값은 scantimecodeout 값으로 자동등록 됩니다.

#### 샷,에셋 삭제
circle 프로젝트 SS_0010 샷 삭제

```bash
$ csi3 -rm item -project circle -name SS_0010
```

#### 에셋등록

에셋이 prop 타입이고 component 형태일 때

```bash
$ csi3 -add item -type asset -project [projectname] -name [Assetname] -assettype prop -assettags prop,component
```

#### 소스등록

```bash
csi3 -add item -name OPN_0010 -type src2 -project TEMP -platepath /source/path
```

등록이 되면 자동으로 OPN_0010 샷에도 OPN_0010_org제목과 /source/path 경로로 소스가 등록됩니다.