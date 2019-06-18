# 개발환경 셋팅

CSI는 Go, mongoDB를 사용하여 진행되는 프로젝트입니다.
공동개발을 위한 개발환경 셋팅방법을 다루는 문서입니다.

### 개발환경셋팅
macOS

```bash
$ brew install go
$ brew install mongodb
```

### 테스트서버
http://csi.lazypic.org

### TravisCI
테스트를 위해서 TravisCI 툴을 사용합니다.

- mongoDB setting: https://docs.travis-ci.com/user/database-setup/

### 라이브러리
사용되는 라이브러리는 아래와 같습니다.

- [dipath](https://github.com/digital-idea/dipath)
- [ditime](https://github.com/digital-idea/ditime)
- [mgo](https://github.com/go-mgo/mgo)
