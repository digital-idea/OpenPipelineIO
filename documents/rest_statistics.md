# Statistics

통계 관련된 RestAPI를 다룹니다.

## GET(Status V1)

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/statistics/projectnum | 프로젝트 갯수를 가지고 옵니다. | . | curl -X GET -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api/statistics/projectnum"
| /api1/statistics/shot | status V1 shot의 샹태를 가지고옵니다. 전체 프로젝트 갯수를 처리합니다. | (project) | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/shot"
| /api1/statistics/shot | status V1 shot의 샹태를 가지고옵니다. 프로젝트 옵션이 있다면 해당 프로젝트 갯수를 처리합니다. | (project) | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/shot?project=TEMP"
| /api1/statistics/asset | status V1 asset의 샹태를 가지고옵니다. 전체 프로젝트 갯수를 처리합니다. | (project) | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/asset"
| /api1/statistics/asset | status V1 asset의 샹태를 가지고옵니다. 프로젝트 옵션이 있다면 해당 프로젝트 갯수를 처리합니다. | (project) | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/asset?project=TEMP"
| /api1/statistics/task | status V1 task의 샹태를 가지고옵니다. 전체 프로젝트 갯수를 처리합니다. | task,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/task?task=comp&type=asset"
| /api1/statistics/task | status V1 task의 샹태를 가지고옵니다. 프로젝트 옵션이 있다면 해당 프로젝트 갯수를 처리합니다. | task,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/task?project=TEMP&task=comp&type=shot"
| /api1/statistics/tag | status V1 tag의 샹태를 가지고옵니다. 전체 프로젝트 갯수를 처리합니다. | name,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/tag?name=tagname&type=asset"
| /api1/statistics/tag | status V1 tag의 샹태를 가지고옵니다. 프로젝트 옵션이 있다면 해당 프로젝트 갯수를 처리합니다. | name,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/tag?project=TEMP&name=tagname&type=shot"
| /api1/statistics/user | status V1 tag의 샹태를 가지고옵니다. 전체 프로젝트 갯수를 처리합니다. | name,task,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/user?name=김한웅&type=asset&task=comp"
| /api1/statistics/user | status V1 tag의 샹태를 가지고옵니다. 프로젝트 옵션이 있다면 해당 프로젝트 갯수를 처리합니다. | name,task,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api1/statistics/user?project=TEMP&name=김한웅&task=comp&type=shot"

## GET(Status V2)

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/statistics/projectnum | 프로젝트 갯수를 가지고 옵니다. | . | curl -X GET -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api/statistics/projectnum"
| /api2/statistics/shot | status V2 shot의 샹태를 가지고옵니다. 전체 프로젝트 갯수를 처리합니다. | (project) | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api2/statistics/shot"
| /api2/statistics/shot | status V2 shot의 샹태를 가지고옵니다. 프로젝트 옵션이 있다면 해당 프로젝트 갯수를 처리합니다. | (project) | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api2/statistics/shot?project=TEMP"
| /api2/statistics/asset | status V2 asset의 샹태를 가지고옵니다. 전체 프로젝트 갯수를 처리합니다. | (project) | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api2/statistics/asset"
| /api2/statistics/asset | status V2 asset의 샹태를 가지고옵니다. 프로젝트 옵션이 있다면 해당 프로젝트 갯수를 처리합니다. | (project) | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api2/statistics/asset?project=TEMP"
| /api2/statistics/task | status V2 task의 샹태를 가지고옵니다. 전체 프로젝트 갯수를 처리합니다. | task,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api2/statistics/task?task=comp&type=asset"
| /api2/statistics/task | status V2 task의 샹태를 가지고옵니다. 프로젝트 옵션이 있다면 해당 프로젝트 갯수를 처리합니다. | task,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api2/statistics/task?project=TEMP&task=comp&type=shot"
| /api2/statistics/tag | status V2 tag의 샹태를 가지고옵니다. 전체 프로젝트 갯수를 처리합니다. | name,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api2/statistics/tag?name=tagname&type=asset"
| /api2/statistics/tag | status V2 tag의 샹태를 가지고옵니다. 프로젝트 옵션이 있다면 해당 프로젝트 갯수를 처리합니다. | name,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api2/statistics/tag?project=TEMP&name=tagname&type=shot"
| /api2/statistics/user | status V2 user의 샹태를 가지고옵니다. 전체 프로젝트 갯수를 처리합니다. | name,task,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api2/statistics/user?name=김한웅&type=asset&task=comp"
| /api2/statistics/user | status V2 user의 샹태를 가지고옵니다. 프로젝트 옵션이 있다면 해당 프로젝트 갯수를 처리합니다. | name,task,(project),type | curl -H "Authorization: Basic <TOKEN>" "https://csi.lazypic.com/api2/statistics/user?project=TEMP&name=김한웅&task=comp&type=shot"

## Curl을 이용해서 DB dump 파일을 다운로드

빅데이터 처리처럼 한번에 DB의 모든 데이터를 가지고 와서 다른 툴에서 분석을 위해 모든 데이터를 가지고 와야 할 때 사용합니다.

```bash
$ curl -H "Authorization: Basic <TOKEN>" -o /download/path/dump.zip https://csi.lazypic.com/export-dump-project
```
