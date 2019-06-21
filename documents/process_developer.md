# 개발환경 셋팅

CSI는 Go, mongoDB를 사용하여 진행되는 프로젝트입니다.
공동개발을 위한 개발환경 셋팅방법을 다루는 문서입니다.

### 개발환경셋팅
macOS

```bash
$ brew install go
$ brew install mongodb
```

#### 에디터
[MS Visual Code](https://code.visualstudio.com)를 사용하며 툴내부 마켓플레이스에서 Go와 관련된 모든 편리한 기능을 인스톨 하고 사용하고 있습니다.
디버그, 실수방지, 개발 속도를 많이 올릴 수 있으니 협업시에는 위 셋팅을 추천합니다.

### 테스트서버
- http://csi.lazypic.org

### TravisCI
테스트를 위해서 [TravisCI](https://docs.travis-ci.com) 를 사용합니다.

### 배포
항상 바이너리 파일 하나를 지향합니다.