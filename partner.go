package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Partner struct {
	// 회사 정보
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name        string             `json:"name"`        // 파트너사 이름
	Size        int                `json:"size"`        // 파트너사 규모
	Domain      string             `json:"domain"`      // 파트너사 사업 영역
	Homepage    string             `json:"homepage"`    // 파트너사 홈페이지 주소
	Address     string             `json:"address"`     // 파트너사 위치 주소
	Phone       string             `json:"phone"`       // 파트너사 전화번호
	Email       string             `json:"email"`       // 파트너사 이메일 주소
	Timezone    string             `json:"timezone"`    // 파트너사 시간대
	Description string             `json:"description"` // 상세설명

	// 계정 정보
	CreateTime string `json:"createtime"` // 파트너사 계정 생성시간
	UpdateTime string `json:"updatetime"` // 파트너사 계정 수정시간
	Thumbnail  bool   `json:"thumbnail"`  // 썸네일 여부

	// 담당자
	Manager      string `json:"manager"`      // 담당자 이름
	ManagerPhone int    `json:"managerphone"` // 담당자 전화번호
	ManagerEmail string `json:"manageremail"` // 담당자 이메일 주소
}
