
# RestAPI Review

## Post

| EndPoint | Description | Attributes | Use case |
| --- | --- | --- | --- |
| /api/review | 리뷰데이터 가지고 오기 | id | `$ curl -X POST -d "id=5f87f82641a789486f3970d1" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/review` |
| /api/rmreview | 리뷰데이터 삭제 | id | `$ curl -X POST -d "id=5f87f82641a789486f3970d1" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/rmreview` |
| /api/searchreview | 리뷰검색 | searchword | `$ curl -X POST -d "searchword=합성3팀" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/searchreview` |
| /api/addreview | 리뷰데이터 추가 | project, name, task, stage, path, author, authornamekor mainversion, subversion, description, fps, (camerainfo), (progress) | `$ curl -X POST -d "project=TEMP&name=SS_0010&task=comp&stage=team&path=test.mov&description=3팀&fps=24&mainversion=1&sebversion=1&authornamekor=김한웅" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/addreview` |
| /api/setreviewstatus | 리뷰상태 변경 | id, status(wait, approve, comment) | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&status=approve" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewstatus` |
| /api/setreviewstage | 리뷰Stage 변경 | id, stage | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&stage=team" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewstage` |
| /api/setreviewproject | 리뷰의 Project 변경 | id, project | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&project=projectname" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewproject` |
| /api/setreviewtask | 리뷰의 Task 변경 | id, task | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&task=task" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewtask` |
| /api/setreviewname | 리뷰의 Name 변경 | id, name | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&name=SS_0010" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewname` 
| /api/setreviewpath | 리뷰의 Path 변경 | id, path | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&path=/show/review/path/reviewdata.mov" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewpath` |
| /api/setreviewmainversion | 리뷰의 MainVersion 변경 | id, mainversion | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&mainversion=1" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewmainversion` |
| /api/setreviewsubversion | 리뷰의 SubVersion 변경 | id, subversion | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&mainversion=1" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewsubversion` |
| /api/setreviewfps | 리뷰의 Fps 변경 | id, fps | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&fps=23.98" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewfps` |
| /api/setreviewdescription | 리뷰의 Description 변경 | id, description | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&description=설명" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewdescription` |
| /api/setreviewcamerainfo | 리뷰의 CameraInfo 변경 | id, camerainfo | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&camerainfo=24mm" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/setreviewcamerainfo` |
| /api/addreviewcomment | 리뷰 Comment 추가 | id, text, stage, media, mediatitle | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&text=수정사항&stage=team&media=/show/drawing.jpg&mediatitle=참고이미지" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/addreviewcomment` |
| /api/editreviewcomment | 리뷰 Comment 수정 | id, time, text | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&status=" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/editreviewcomment` |
| /api/rmreviewcomment | 리뷰 Comment 삭제 | id, time | `$ curl -X POST -d "id=5f87f82641a789486f3970d1&time=2020-05-21T09:00:00%2B09:00" -H "Authorization: Basic <Token>" https://csi.lazypic.org/api/rmreviewcomment` |
| /api/uploadreviewdrawing | 리뷰 드로잉 이미지 업로드 | id, frame | `$ curl -X POST -H "Authorization: Basic <Token>" -F  id=5f4edbe16e59c4695abb12d1 -F frame=101 -F "image=@/path/reviewdrawing.png" https://csi.lazypic.org/api/uploadreviewdrawing`|