package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Endpoint struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	DNS           string             `json:"dns"`           // 도메인 네임 서버. 예) platform.lazypic.com
	Endpoint      string             `json:"endpoint"`      // 엔드포인트 주소 예) /api/v1/users
	Method        []string           `json:"method"`        // GET,PUT,DELETE,POST,OPTIONS
	Parameter     []string           `json:"parameter"`     // Endpoint 에서 지원하는 옵션
	CORS          string             `json:"cors"`          // Cross-Origin Resource Sharing: origin || *
	Type          string             `json:"type"`          // Resource 페이지, RestAPI, Webpage 인지 설정하는 옵션
	StorageType   string             `json:"storagetype"`   // AWS S3, 자체서버를 구축해서 파일을 저장할지 결정
	DB            string             `json:"db"`            // Endpoint의 특성과 잘 맞는 DB 예) hBase, mongoDB, mysql, dynamodb, redis
	Description   string             `json:"description"`   // 설명
	IsServerless  bool               `json:"isserverless"`  // 이 기능을 수행함에 따라 서버를 사용하는지 서버리스를 사용하는지 여부, 예) AWS Lambda
	User          bool               `json:"user"`          // 사용자 공개
	Developer     bool               `json:"developer"`     // 개발자 공개
	Admin         bool               `json:"admin"`         // 관리자 공개
	Security      bool               `json:"security"`      // 보안이 필요한가
	IsAsset       bool               `json:"isasset"`       // 회사의 자산이 될 수 있는 Endpoint 인가?
	IsPatent      bool               `json:"ispatent"`      // 특허와 관련된 기술이 처리되는 Endpoint 인가?
	IsUpload      bool               `json:"isupload"`      // 업로드가 되는 Endpoint인가?
	Token         string             `json:"token"`         // 보안토큰
	Tags          []string           `json:"tags"`          // 태그
	SupportFormat []string           `json:"supportformat"` // 지원포멧 예) json, graphQL, form, xml, binary
	Curl          string             `json:"curl"`          // curl 예시
	Category      string             `json:"category"`      // 카테고리: 유저,결제,관리,CS(고객대응),메일,알림,메시징,내저징데이터,정산,검색,통계,마이페이지,공통사항,인증,웹소캣,공지사항,시스템 관리,메인페이지,모바일지원,광고,혜택,이벤트,친구관리,찜,장바구니,공유,채팅,기기관리,통화,음성,카메라,구매,보안,선물,뉴스,모니터링,구독,데이터분석,머신러닝,번역
	Partner       string             `json:"partner"`       // Endpoint 제작사 또는 유지보수 관리사
	Progress      string             `json:"progress"`      // 진행률. bootstrap progress바의 형이 string이다.
}
