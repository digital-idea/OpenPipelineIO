package main

type Partner struct {
	// 회사 정보
	ID          string `json:"id"`          // ID
	Name        string `json:"name"`        // 파트너사 이름
	Homepage    string `json:"homepage"`    // 파트너사 홈페이지 주소
	Address     string `json:"address"`     // 파트너사 위치 주소
	Phone       string `json:"phone"`       // 파트너사 전화번호
	Email       string `json:"email"`       // 파트너사 이메일 주소
	Timezone    string `json:"timezone"`    // 파트너사 시간대
	Description string `json:"description"` // 상세설명

	// 계정 정보
	Createtime string `json:"createtime"` // 파트너사 계정 생성시간
	Updatetime string `json:"updatetime"` // 파트너사 계정 수정시간
	Thumbnail  bool   `json:"thumbnail"`  // 썸네일 여부

	// 담당자
	Manager      string `json:"manager"`      // 담당자 이름
	ManagerPhone int    `json:"managerphone"` // 담당자 전화번호
	ManagerEmail string `json:"manageremail"` // 담당자 이메일 주소
}
