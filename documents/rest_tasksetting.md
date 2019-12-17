# Tasksetting RestAPI
Tasksetting Restapi 입니다.

## Get
| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/shottasksetting | shot tasksetting 정보를 가지고 온다. | `$ curl http://csi.lazypic.org/api/shottasksetting` |
| /api/assettasksetting | asset tasksetting 정보를 가지고 온다. | `$ curl http://csi.lazypic.org/api/assettasksetting` |

## POST
| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/tasksetting | 인수를 입력받고 task에서 사용하는 경로를 반환한다. | project, name, task, type, assettype, os, seq, cut, userid | `$ curl -X POST -d "project=TEMP&seq=SS&cut=0010&task=comp" http://csi.lazypic.org/api/tasksetting` |