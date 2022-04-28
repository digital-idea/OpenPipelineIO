package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Partner struct {
	ID                         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name                       string             `json:"name"`                       // 파트너사 이름
	Size                       int                `json:"size"`                       // 파트너사 규모
	Domain                     string             `json:"domain"`                     // 파트너사 사업 영역
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
}
