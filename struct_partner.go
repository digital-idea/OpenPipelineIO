package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Partner struct {
	ID                         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name                       string             `json:"name"`                       // 파트너사 이름
	Symbol                     string             `json:"symbol"`                     // 주식종목코드
	Domain                     string             `json:"domain"`                     // 파트너사 사업 영역
	Size                       string             `json:"size"`                       // 파트너사 규모 명수
	Homepage                   string             `json:"homepage"`                   // 파트너사 홈페이지 주소
	Address                    string             `json:"address"`                    // 파트너사 위치 주소
	Phone                      string             `json:"phone"`                      // 파트너사 전화번호
	Email                      string             `json:"email"`                      // 파트너사 이메일 주소
	Location                   string             `json:"location"`                   // Location 지역약자: 인도, 상암, 분당
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
	Tags                       []string           `json:"tags"`                       // 태그
	Companytype                string             `json:"companytype"`                // 법인, 개인, 프리렌서 인가?
	ContactPoint               string             `json:"contactpoint"`               // 컨택포인트(누구의 소개인지, 어디서 만났는지)
	PmsURL                     string             `json:"pmsurl"`                     // 샷그리드, Ftrack, 윕홀, OpenPiplineIO 등 파이프라인툴 URL추가
}

// Money 는 돈과 관련된 자료구조이다. 어떤 프로젝트에서 누가, 누구에게, 언제 얼마를 주는가에 대한 정보
type Money struct {
	ID                           primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	MoneytypeID                  string             `json:"moneytypeid"`                  // firstestimate 최초견적, contractestimate 계약견적, additionalestimate 추가견적, downpayment 계약금, interimpayment 중도금, balance 잔금, surchage 추가금
	Project                      string             `json:"project"`                      // 프로젝트
	Sender                       string             `json:"sender"`                       // 보내는이
	Recipient                    string             `json:"recipient"`                    // 받는이
	Amount                       float64            `json:"amount"`                       // 액수
	Date                         string             `json:"date"`                         // 전달 날짜
	MonetaryUnit                 string             `json:"monetaryunit"`                 // 단위 : KRW,USD,EUR,CNY,JPY,VND ... Policy : ISO4217
	Description                  string             `json:"description"`                  // 내용
	Status                       string             `json:"status"`                       // 절차검토, 품의...
	IssuanceElectronicTaxInvoice bool               `json:"issuanceelectronictaxinvoice"` // 세금계산서 발행
	IssuanceInvoice              bool               `json:"issuanceinvoice"`              // 인보이스 발행여부
}

type Moneytype struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Order       float64            `json:"order"`       // 순서
	Name        string             `json:"name"`        // 영문이름
	Description string             `json:"description"` // 설명
	// firstestimate 최초견적, contractestimate 계약견적, additionalestimate 추가견적, downpayment 계약금, interimpayment 중도금, balance 잔금, surchage 추가금
}

// ProjectForPartner 자료구조는 프로젝트와 파트너사이의 관계를 다루는 자료구조
type ProjectForPartner struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	ProjectName   string             `json:"projectname"`   // 프로젝트 이름
	PartnerName   string             `json:"partnername"`   // 파트너이름
	RNR           string             `json:"rnr"`           // 계약,진행상태
	ProjectBudget float64            `json:"projectbudget"` // 프로젝트 비용

	StartDate               string  `json:"startdate"`               // 프로젝트 시작일
	EndDate                 string  `json:"enddate"`                 // 프로젝트 완료일
	PercentageOfTotalBudget float64 `json:"percentageoftotalbudget"` // 나가는 비용이 프로젝트 총 비용의 몇 퍼센트인가?
	Description             string  `json:"description"`             // 외주내용
	PricePerShot            float64 `json:"pricepershot"`            // 컷당 가격
	PricePerFrame           float64 `json:"priceperframe"`           // 프레임당 가격
	PartnerInternalManager  string  `json:"partnerinternalmanager"`  // 파트너 내부관리자
	ManagerEmail            string  `json:"manageremail"`            // 파트너 내부관리자 메일
	PaymentCycle            string  `json:"paymentcycle"`            // 지급회차 1/2, 4/6: 현재 지급단계, 총 지급횟수
	PaymentDateForClient    int     `json:"paymentdateforclient"`    // 클라이언트에게 돈을 받는날짜, 프로젝트(프로젝트 월별 지급일)
	PaymentDateForVender    int     `json:"paymentdateforvender"`    // 벤더에게 주는 날짜, 프로젝트 진행시 벤더에게 돈을 주는 날짜
	Language                string  `json:"language"`                // 사용언어: 커뮤니케이션 언어
	Messenger               string  `json:"messenger"`               // 사용 메신저 종류
	MessengerID             string  `json:"messengerid"`             // 메신저 ID
	Manday                  int     `json:"manday"`                  // 예상 맨데이, 회계 작성시 필요, 감사시 필요
	LeftoverShot            int     `json:"leftovershot"`            // 남은컷
	TotalShot               int     `json:"totalshot"`               // 총컷
}

type Event struct {
	ID      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Date    string             `json:"date"`    // 이벤트 날짜(필수)
	Title   string             `json:"title"`   // 이벤트 제목(필수)
	Partner string             `json:"partner"` // 관련 파트너사가 있다면 입력한다.
	Project string             `json:"project"` // 관련 프로젝트가 있다면 입력한다.
}

type Contract struct {
	ID      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Partner string             `json:"partner"` // 계약된 파트너 이름
	Title   string             `json:"title"`   // 계약서 이름
	Date    string             `json:"date"`    // 등록일
	Url     string             `json:"url"`     // URL, pdf 일이 업로드된 url
	Price   float64            `json:"price"`   // 가격
	Amount  string             `json:"amount"`  // 분량, 컷수
}

// Evaluation 은 평가 자료구조이다.
type Evaluation struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Partner       string             `json:"partner"`       // 파트너 이름
	Quality       int                `json:"quality"`       // 퀄리티
	Schedule      int                `json:"schedule"`      // 스케줄
	Communication int                `json:"communication"` // 커뮤니케이션
	Total         int                `json:"total"`         // 종합점수
	GeneralReview string             `json:"generalreview"` // 총평, 총평을 위주로 본다.
	CreateDate    string             `json:"createdate"`    // 평가일
}

// Step 자료구조는 워크플로우 또는 R&R 과정에 대한 단계 시각화, 진행률 체크하기 위한 자료구조이다.
// 아래는 예시이다.
// - 컨텍중(소개), 미팅접점 <- 경영지원, 마케팅팀
// - 상호 스케줄 체크 <- PD
// - NDA발송 <- PM
// - 날인된 NDA 받기 <- 경재실
// - 관련자료전달, 프로젝트 이름을 밝히지 않는다.
// - 데이터 아웃, 프로젝트 이름을 밝히지 않는다.
// - 견적받기
// - 견적조정
// - 견적시 우리가 던진금액이 얼마인지 체크하기
// - 견적서 다시 받기
// - 견적서 합의
// - 수령후 계약서 전달
// - 계약서 조율
// - 날인된 최종 계약서 받기
// - 2트렉 품의 올리기
// - 이 금액으로 이 회사와 계약할 수 있게 확인
// - 품의
// - 결재중
// - 결제완료
// - 날인계약서 도착
// - 결제완료된 계약서와 합의된 견적서를 회계팀에 전달하기
// - 날인후 회계팀이 투자 회사에 전달
// - 최종 계약서 스캔본 공유
type Step struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Order       float64            `json:"order"`       // 순서
	Name        string             `json:"name"`        // R&R 이름
	Description string             `json:"description"` // 내용
}
