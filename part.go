package main

// Part 는 팀,부서 정보입니다.
type Part struct {
	ID       string `json:"id"`       // Part ID
	NameKor  string `json:"namekor"`  // Part 한글이름
	NameEng  string `json:"nameeng"`  // Part 영문이름
	Category string `json:"category"` // 상위부서 또는 관련 부서
}
