package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Partner struct {
	ID                         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name                       string             `json:"name"`                       // 파트너사 이름
	StockSymbol                string             `json:"stocksymbol"`                // 주식종목코드
	Domain                     string             `json:"domain"`                     // 파트너사 사업 영역
	Size                       string             `json:"size"`                       // 파트너사 규모 명수
	Homepage                   string             `json:"homepage"`                   // 파트너사 홈페이지 주소
	Address                    string             `json:"address"`                    // 파트너사 위치 주소
	Phone                      string             `json:"phone"`                      // 파트너사 전화번호
	Email                      string             `json:"email"`                      // 파트너사 이메일 주소
	Timezone                   string             `json:"timezone"`                   // 파트너사 시간대
	Description                string             `json:"description"`                // 상세설명
	BusinessRegistrationNumber string             `json:"businessregistrationnumber"` // 사업자등록번호
	Manager                    string             `json:"manager"`                    // 담당자 이름
	ManagerPhone               string             `json:"managerphone"`               // 담당자 전화번호
	ManagerEmail               string             `json:"manageremail"`               // 담당자 이메일 주소
	FTP                        string             `json:"ftp"`                        // FTP 주소
	FTPID                      string             `json:"ftpid"`                      // FTP ID
	FTPPW                      string             `json:"ftppw"`                      // FTP PW
	Opentime                   string             `json:"opentime"`                   // Open Time
	Closedtime                 string             `json:"closedtime"`                 // Closed Time
	PaymentDate                string             `json:"paymentdate"`                // 지급일
	Bank                       string             `json:"bank"`                       // 은행명
	BankAccount                string             `json:"bankaccount"`                // 계좌번호
	MonetaryUnit               string             `json:"monetaryunit"`               // 달러인지 유로인지 위안인지 체크하기
	ProjectHistory             string             `json:"projecthistory"`             // 과거 했던 작품
	Reputation                 string             `json:"reputation"`                 // 평판
	Status                     string             `json:"status"`                     // 상태: R&R 프로잭트진행, 자금팀, 경영지원, 글씨로 쓰게 해두기.
	Codename                   string             `json:"codename"`                   // 벤더 코드명
	IsAbroad                   bool               `json:"isabroad"`                   // 국내, 해외 체크
	IsClient                   bool               `json:"isclient"`                   // 내 기준에서 갑(Clinet), 내 기준에서 을(Partner, Vender) 인지 체크
	Progress                   string             `json:"progress"`                   // 상태: 외부진행상황. %
	DeliveryDates              []Date             `json:"deliverydates"`              // 딜리버리 일정 관리(매주, 매월 나가야하는 것) 프로젝트별로 저장되어야 한다.
	Contracts                  []Contract         `json:"contracts"`                  // 계약서 등록 삭제기능, 과정을 볼 수 있도록 하기.
	Tags                       []string           `json:"tags"`                       // 태그
	Type                       string             `json:"type"`                       // 법인, 개인, 프리렌서 인가?
	// 컨택포인트(누구의 소개인지, 어디서 만났는지)

}

// 어떤 프로젝트에서 누가, 누구에게, 얼마를 주는가

// ProjectForPartner // 프로젝트와 파트너사이의 관계를 다루는 자료구조

// 파트너 <-> 프로젝트의 관계
// 프로젝트 비용
// 파트너에에게 가는 분량
// 나가는 비용이 몇 퍼센트인가?
// 과거 견적 히스토리 내용
// 과거 견적 액수
// 프로젝트 타입, 타입설정, "A","B","C"
// 프로젝트: 타입설정
// 계약상태: R&R
// 외주내용:
// 분량
// 최초견적
// 계약견적
// 추가견적
// 컷단가(프레임단가) 가격
// 프로젝트 시작일
// 프로젝트 완료일
// 밴더 내부관리자
// 예상맨데이(회계상필요) <- 감사시 필요함.
// 계약금 가격
// 계약금 지급일
// 중도금 가격...
// 중도금 지급일
// 잔금 가격
// 잔금 지급일
// 추가금 가격
// 추가금 지급일
// 지급회차 1/2, 4/6: 현재 지급단계, 총 지급횟수
// 프로젝트별 월별 지급일 <- ""
// 클라이언트에게 돈을 받는날짜, 프로젝트
// 벤더에게 주는 날짜, 프로젝트 진행시 벤더에게 돈을 주는 날짜
// 전매출의 10%(작년매출액기준) 넘는 프로젝트는 공시 -> IR 게시판, 프로젝트 매출
// 실무평가: 퀄리티, 스케줄, 커뮤니케이션, 종합, 총평 <- 실무자, 프렙의 경우 작업자, 프로젝트에 대한...
// Location 지역약자: 인도, 분당 <- partner
// 사용언어: 커뮤니케이션 언어 <- partner
// 사용 메신저 종류: <- partner
// 메신저 ID: <- partner

// R&R 과정에 대한 단계 시각화, 진행률 체크

// 컨텍 -> 계약완료까지의 과정
// 컨텍(소개)
// 상호 스케줄 체크
// NDA발송
// 날인된 NDA 받기
// 관련자료전달, 프로젝트 이름을 밝히지 않는다.
// 데이터 아웃, 프로젝트 이름을 밝히지 않는다.
// 견적받기
// 견적조정
// 견적시 우리가 던진금액이 얼마인지 체크하기
// 견적서 다시 받기
// 견적서 합의
// 수령후 계약서 전달
// 계약서 조율
// 날인된 최종 계약서 받기
// 2트렉 품의 올리기
// 이 금액으로 이 회사와 계약할 수 있게 확인
// 품의
// 결재중
// 결제완료
// 날인계약서 도착
// 결제완료된 계약서와 합의된 견적서를 회계팀에 전달하기
// 날인후 회계팀이 투자 회사에 전달
// 최종 계약서 스캔본 공유

type Date struct {
	Title string `json:"title"` // 이벤트제목
	Date  string `json:"date"`  // 날짜
}

type Contract struct {
	Title string `json:"title"` // 계약서 이름
	Date  string `json:"date"`  // 등록일
	Url   string `json:"url"`   // URL, 파일이 업로드된 url
}
