# 개발환경 셋팅

CSI는 go 언어로 작성되는 프로젝트입니다.
공동개발을 위한 개발환경 셋팅방법을 다루는 문서입니다.
CSI 개발자가 추가되었을 때 교육문서로도 활용됩니다.
서버에 위치한 go 보다 로컬에 go를 설치하고 사용하시면 더 빠르게 개발할 수 있습니다.

#### 개발자가 새로오면 준비할 사항
- 오리엔테이션 및 회사 개발환경 소개
- golang 소개
- mongoDB 소개
- 회사카드로 자신이 볼 책을 결제할 수 있음을 알리기
- 개발 & 코드리뷰 파이프라인 설명

#### .bashrc 설정
- 아래내용을 ~/.bashrc에 추가하세요.

```
export GOROOT="/home/{사번}/goroot" # 이곳에 GO 를 설치합니다. 사본위치 : /lustre/INHouse/Tool/go1.9
export GOPATH="/home/{사번}/gopath"
export GOBIN="$GOPATH/bin"
export PATH="$PATH:$GOROOT/bin:$GOBIN:/dept/rnd/bin"
```

#### 서비스배포
- 배포 이후 장애를 잘 처리하기 위해서 급한 퍼블리쉬가 아니면 월,화,수,목 오전에만 사이트를 갱신합니다.

#### 라이브러리
사용되는 라이브러리는 아래와 같습니다.

- [dipath](https://github.com/digital-idea/dipath)
- [ditime](https://github.com/digital-idea/ditime)
- [mgo](https://github.com/go-mgo/mgo)
