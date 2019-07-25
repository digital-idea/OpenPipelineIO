package main

// Division 은 본부 정보입니다.
type Division struct {
	ID      string `json:"id"`      // Division ID
	NameKor string `json:"namekor"` // 한글이름
	NameEng string `json:"nameeng"` // 영문이름
}

// Department 은 부 정보입니다.
type Department struct {
	ID      string `json:"id"`      // Division ID
	NameKor string `json:"namekor"` // 한글이름
	NameEng string `json:"nameeng"` // 영문이름
}

// Team 는 팀 정보입니다.
type Team struct {
	ID      string `json:"id"`      // Part ID
	NameKor string `json:"namekor"` // 한글이름
	NameEng string `json:"nameeng"` // 영문이름
}

// Role 은 직책 정보입니다.
type Role struct {
	ID      string `json:"id"`      // Role ID
	NameKor string `json:"namekor"` // 한글이름
	NameEng string `json:"nameeng"` // 영문이름
}

// Position 은 직급 정보입니다.
type Position struct {
	ID      string `json:"id"`      // Role ID
	NameKor string `json:"namekor"` // 한글이름
	NameEng string `json:"nameeng"` // 영문이름
}

// Organization 은 조직 정보입니다.
type Organization struct {
	Division   `json:"division"`
	Department `json:"department"`
	Team       `json:"team"`
	Role       `json:"role"`
	Position   `json:"position"`
}
