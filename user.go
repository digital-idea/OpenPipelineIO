package main

import (
	"fmt"
	"time"
)

// AccessLevel 사용자의 엑세스 레벨이다.
type AccessLevel int

const (
	// UnknownAccessLevel 접근이 가장 낮은 레벨이다. 0레벨
	UnknownAccessLevel = AccessLevel(iota)
	// GuestAccessLevel 손님 1레벨
	GuestAccessLevel
	// ClientsAccessLevel 클라이언트 2레벨
	ClientsAccessLevel
	// ArtistAccessLevel 아티스트 3레벨
	ArtistAccessLevel
	// LeadAccessLevel 팀장 4레벨
	LeadAccessLevel
	// PmAccessLevel 프로젝트메니징 5레벨
	PmAccessLevel
	// SupervisorAccessLevel 슈퍼바이저 6레벨
	SupervisorAccessLevel
	// IoAccessLevel IO 매니저 7레벨(서버권한 이슈)
	IoAccessLevel
	// PdAccessLevel PD레벨 8레벨(자금 이슈)
	PdAccessLevel
	// HqAccessLevel 경영지원(인사,재무)레벨 9레벨(전사적 자금 이슈)
	HqAccessLevel
	// DeveloperAccessLevel 개발자 10레벨
	DeveloperAccessLevel
	// AdminAccessLevel 관리자 11레벨(Root권한)
	AdminAccessLevel
)

// User 는 사용자 정보입니다.
type User struct {
	ID                string         `json:"id"`                // 사용자 ID(사번). 손님 및 클라이언트는 사번이 없다.(예외)
	Password          string         `json:"password"`          // 사용자 비밀번호
	PasswordAttempt   int            `json:"passwordattempt"`   // 잘못된 패스워드 시도횟수
	FirstNameKor      string         `json:"firstnamekor"`      // 한글이름: 이름
	LastNameKor       string         `json:"lastnamekor"`       // 한글이름: 성
	FirstNameEng      string         `json:"firstnameeng"`      // 영문이름, Firstname, 외국인은 이름이 길기때문에 이 이름을 닉네임으로 사용한다.
	LastNameEng       string         `json:"lastnameeng"`       // 영문이름, Lastname
	FirstNameChn      string         `json:"firstnamechn"`      // 한자이름: 이름, 중국에서 한자명으로 엔딩크레딧 요청이 있다
	LastNameChn       string         `json:"lastnamechn"`       // 한자이름: 성, 중국에서 한자명으로 엔딩크레딧 요청이 있다
	Email             string         `json:"email"`             // 사내 메일
	EmailExternal     string         `json:"emailexternal"`     // 외부 이메일
	Phone             string         `json:"phone"`             // 핸드폰
	Hotline           string         `json:"hotline"`           // 사내전화
	Location          string         `json:"location"`          // 사내위치: 여러층에 나누어져 있을 때 층수, 대략의 위치정보
	Tags              []string       `json:"tags"`              // 사용자 태그
	Timezone          string         `json:"timezone"`          // 타임존(이슈지역 : 한국, 중국, 캐나다, 미국)
	AccessLevel       AccessLevel    `json:"accesslevel"`       // 소프트웨어의 액세스 레벨
	Updatetime        string         `json:"updatetime"`        // 업데이트 시간.
	Createtime        string         `json:"createtime"`        // 계정생성 시간.
	IsLeave           bool           `json:"isleave"`           // 퇴사여부. 약자로 BSR(빤스런) 이라고 불린다.
	LastIP            string         `json:"lastip"`            // 최근 접속 IP
	LastPort          string         `json:"lastport"`          // 최근 접속 Port
	Thumbnail         bool           `json:"thumbnail"`         // 썸네일 유무
	Token             string         `json:"token"`             // restAPI Token 키
	Organizations     []Organization `json:"organizations"`     // 조직정보
	OrganizationsForm string         `json:"organizationsform"` // 가입시 사용된 조직정보 문자
	AccessProjects    []string       `json:"accessprojects"`    // 사용자에게 허가된 프로젝트 리스트
	EmployeeNumber    string         `json:"employeenumber"`    // 사원번호
}

// Token 자료구조. 사용자가 가입될 때 user.token DB에 저장된다. 모든 유저의 Token를 매번 비교하지않고, Token 키의 유효성을 바로 체크하기 위해서 사용한다.
type Token struct {
	Token       string      `json:"token"`       // restAPI Token 키
	AccessLevel AccessLevel `json:"accesslevel"` // 소프트웨어의 액세스 레벨
	ID          string      `json:"id"`          // restAPI에 대한 로그 기록을 남기기 위해서는 토큰키로 아이디를 가지고올 수 있어야 한다.
}

// NewUser 는 새로운 유저를 생성할 때 사용한다.
func NewUser(id string) *User {
	return &User{
		ID:          id,
		AccessLevel: AccessLevel(*flagSignUpAccessLevel),
		Updatetime:  time.Now().Format(time.RFC3339),
		Createtime:  time.Now().Format(time.RFC3339),
		IsLeave:     false,
	}
}

// SetTags 메소드는 조직정보를 분석하여 사용자에게 적용해야할 태그를 설정한다.
func (u *User) SetTags() {
	var tags []string
	for _, org := range u.Organizations {
		if org.Division.Name != "" {
			tags = append(tags, org.Division.Name)
		}
		if org.Department.Name != "" {
			tags = append(tags, org.Department.Name)
		}
		if org.Team.Name != "" {
			tags = append(tags, org.Team.Name)
		}
		if org.Role.Name != "" {
			tags = append(tags, org.Role.Name)
		}
		if org.Position.Name != "" {
			tags = append(tags, org.Position.Name)
		}
	}
	tags = append(tags, u.Tags...)
	u.Tags = UniqueSlice(tags)
}

func (u *User) emailString(lang string) string {
	switch lang {
	case "ko", "kor", "korean":
		if u.LastNameKor == "" && u.FirstNameKor == "" {
			return fmt.Sprintf("%s<%s>", u.Email, u.Email)
		} else {
			return fmt.Sprintf("%s%s<%s>", u.LastNameKor, u.FirstNameKor, u.Email)
		}
	case "eg", "eng", "english":
		if u.LastNameEng == "" && u.FirstNameEng == "" {
			return fmt.Sprintf("%s<%s>", u.Email, u.Email)
		} else {
			return fmt.Sprintf("%s%s<%s>", u.FirstNameEng, u.LastNameEng, u.Email)
		}
	case "mailonly", "onlymail":
		return fmt.Sprintf("%s<%s>", u.Email, u.Email)
	default:
		// 한국에서 자주 사용하는 툴이라서 한국어로 기본 설정한다.
		if u.LastNameKor == "" && u.FirstNameKor == "" {
			return fmt.Sprintf("%s<%s>", u.Email, u.Email)
		} else {
			return fmt.Sprintf("%s%s<%s>", u.LastNameKor, u.FirstNameKor, u.Email)
		}
	}

}
