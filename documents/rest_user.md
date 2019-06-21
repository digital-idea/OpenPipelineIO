# RestAPI User

id가 `woong`인 유저의 정보를 가지고 오기

```bash
$ curl http://csi.lazypic.org/api/user?id=woong
```

합성팀, 1팀소속, 팀장을 가지고 오는 방법
```bash
curl http://192.168.219.104/api/users?searchword="합성팀,1팀,팀장"
```