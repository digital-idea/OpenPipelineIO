
# RestAPI Review

## Post

| EndPoint | Description | Attributes | Use case |
| --- | --- | --- | --- |
| /api/addreview | 리뷰데이터 추가 | project, name, task, path, author, authornamekor mainversion, subversion, description, fps, (camerainfo), (progress) | `$ curl -X POST -d "project=TEMP&name=SS_0010&task=comp&path=test.mov&description=3팀&fps=24&mainversion=1&sebversion=1&authornamekor=김한웅" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/addreview` |
| /api/editreviewcomment | 리뷰 Comment 수정 | id, time, text | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&time=2020-10-15T16:20:35%2B09:00&text=commet_text" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/editreviewcomment` |