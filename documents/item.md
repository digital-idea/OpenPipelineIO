# Item
샷, 에셋 관련 터미널 명령어 사용법 입니다.

#### 샷등록
DB값만 생성되며, 샷 폴더가 생성되지는 않습니다.

```bash
$ csi3 -add item -project [projectname] -name [SS_0010] -type [org]
```

#### 샷,에셋 삭제

```bash
$ csi3 -rm item -project [projectname] -name [SS_0010] -type [org]
```

#### 에셋등록

에셋이 prop 타입이고 component 형태일 때

```bash
$ csi3 -add item -type asset -project [projectname] -name [Assetname] -assettype prop -assettags prop,component
```
