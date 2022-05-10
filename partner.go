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

}

type Date struct {
	Title string `json:"title"` // 이벤트제목
	Date  string `json:"date"`  // 날짜
}

type Contract struct {
	Title string `json:"title"` // 계약서 이름
	Date  string `json:"date"`  // 등록일
	Url   string `json:"url"`   // URL
}
