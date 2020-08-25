
# RestAPI Review

## Post

| EndPoint | Description | Attributes | Use case |
| --- | --- | --- | --- |
| /api/addreview | 리뷰데이터 추가 | project, name, task, path, author, mainversion, subversion, description, fps, (camerainfo), (progress) | `$ curl -X POST -d "project=TEMP&name=SS_0010&task=comp&path=test.mov&description=3팀&fps=24&mainversion=1&sebversion=1" https://csi.lazypic.org/api/addreview` |