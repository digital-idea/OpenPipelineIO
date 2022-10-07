package main

import "go.mongodb.org/mongo-driver/bson/primitive"

//https://fullcalendar.io/docs/event-object
// https://fullcalendar.io/docs/event-parsing
type FullCalendarEvent struct {
	ID               primitive.ObjectID        `bson:"_id" json:"id,omitempty"`
	GroupID          string                    `json:"groupId"`          // 그룹ID
	AllDay           bool                      `json:"allDay"`           // 하루종일 진행되는 이벤트인가?
	Start            string                    `json:"start"`            // 시작시간
	End              string                    `json:"end"`              // 끝나는 시간
	StartStr         string                    `json:"startStr"`         // ISO8601 Time
	EndStr           string                    `json:"endStr"`           // ISO8601 Time
	Title            string                    `json:"title"`            // 제목
	Url              string                    `json:"url"`              // 클릭하면 이동되는 URL
	ClassNames       []string                  `json:"classNames"`       // html에 렌더링할 때 attached 할 클레스 이름
	Editable         bool                      `json:"editorble"`        // 수정 여부
	StartEditable    bool                      `json:"startEditable"`    // 시작시간 수정 여부
	EndEditable      bool                      `json:"endEditable"`      // 끝시간 수정 여부
	DurationEditable bool                      `json:"durationEditable"` // 기간수정 여부
	ResourceEditable bool                      `json:"resourceEditable"` // 리소스 수정가능 여부
	Display          string                    `json:"display"`          // 'auto','block','list-item','background','inverse-background','none' -> https://fullcalendar.io/docs/eventDisplay
	Color            string                    `json:"color"`            // Background + Border Color
	BackgroundColor  string                    `json:"backgroundColor"`  // 배경색
	BorderColor      string                    `json:"borderColor"`      // 테두리색
	TextColor        string                    `json:"textColor"`        // 글자색
	ExtendedProps    FullCalendarExtendedProps `json:"extendedProps"`    // 나머지 필요한 기능
	ResourceId       string                    `json:"resourceId"`       // Resource ID
	ResourceIds      string                    `json:"resorceIds"`       // Resource IDs
	// Constraint // 일정을 설정할 때 제한을 둘때 사용한다.
	// Overlap // false 또는 함수를 다룬다. Go 자료구조에서는 다루지 않는다.
}

// https://fullcalendar.io/docs/resource-object
// https://fullcalendar.io/docs/resource-parsing
type FullCalendarResource struct {
	ID                   primitive.ObjectID        `bson:"_id" json:"id,omitempty"`
	Title                string                    `json:"title"`                // 제목
	ExtendedProps        FullCalendarExtendedProps `json:"extendedProps"`        // 나머지 필요한 기능
	EventColor           string                    `json:"eventColor"`           // 배경색 + 테두리색
	EventBackgroundColor string                    `json:"eventBackgroundColor"` // 배경색
	EventBorderColor     string                    `json:"eventBorderColor"`     // 테두리색
	EventTextColor       string                    `json:"eventTextColor"`       // 글자색
	EventClassNames      []string                  `json:"eventClassNames"`      // html에 렌더링할 때 attached 할 클레스 이름
	Children             []FullCalendarResource    `json:"children"`             // 자식의 Resource IDs
	ParentId             string                    `json:"parentId"`             // 부모의 Resource ID
	// EventAllow
	// EventConstraint // 일정을 설정할 때 제한을 둘때 사용한다.
	// EventOverlap // false 또는 함수를 다룬다. Go 자료구조에서는 다루지 않는다.
}

type FullCalendarExtendedProps struct {
	ItemID          string   `json:"itemid"`          // 아이템 아이디(샷,에셋)
	UserID          string   `json:"userid"`          // User 아이디
	PartnerCodename string   `json:"partnercodename"` // 파트너사 코드네임
	Project         string   `json:"project"`         // 프로젝트명, 감독컨펌, 데이터 아웃 날짜를 구하기 위해 프로젝트가 필요하다.
	Deadline3d      string   `json:"deadline3d"`      // 3D 데드라인(검토하기)
	Deadline2d      string   `json:"deadline2d"`      // 2D 데드라인(Comp 데드라인과 같다.)
	Task            string   `json:"task"`            // Task 이름
	TaskStartDate   string   `json:"taskstartdate"`   // 작업 시작일 RFC3339
	TaskPreDeadline string   `json:"taskpredeadline"` // 1차 마감일 RFC3339
	TaskDeadline    string   `json:"taskdeadline"`    // 2차 마감일 RFC3339
	Pipelinestep    string   `json:"pipelinestep"`    // pipelinestep
	Tags            []string `json:"tags"`            // Tags
}
