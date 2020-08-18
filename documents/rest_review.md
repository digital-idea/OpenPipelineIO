
# RestAPI Review

## Post

| uri | description | attribute name | example |
| --- | --- | --- | --- |
| /api/addreview | 리뷰데이터 추가 | project, name, task, path, author, mainversion, subversion, (progress) | `$ curl -X POST -d "project=TEMP&name=SS_0010&task=comp&path=test.mov" https://csi.lazypic.org/api/addreview` |