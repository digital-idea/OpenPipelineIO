# Statistics

통계 관련된 RestAPI를 다룹니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/statistics/projectnum | 프로젝트 갯수를 가지고 옵니다. | . | curl -X GET -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api/statistics/projectnum"
| /api1/statistics/shot | status V1 shot의 샹태를 가지고옵니다. 프로젝트 옵션이 없다면 전체 프로젝트 갯수를 처리합니다. | (project) | curl -X GET -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/shot?project=TEMP"


## Curl을 이용해서 DB dump 파일을 다운로드

빅데이터 처리처럼 한번에 DB의 모든 데이터를 가지고 와서 다른 툴에서 분석을 위해 모든 데이터를 가지고 와야 할 때 사용합니다.

```bash
$ curl -H "Authorization: Basic <TOKEN>" -o /download/path/dump.zip https://csi.lazypic.com/export-dump-project
```
