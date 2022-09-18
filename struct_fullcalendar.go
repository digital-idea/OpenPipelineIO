package main

// https://fullcalendar.io/docs/event-object
type FullCalendarEvent struct {
	ID               string            `json:"id"`               // ID
	GroupID          string            `json:"groupId"`          // 그룹ID
	AllDay           bool              `json:"allDay"`           // 하루종일 진행되는 이벤트인가?
	Start            string            `json:"start"`            // 시작시간
	End              string            `json:"end"`              // 끝나는 시간
	StartStr         string            `json:"startStr"`         //  ISO8601 Time
	EndStr           string            `json:"endStr"`           //  ISO8601 Time
	Title            string            `json:"title"`            // 제목
	Url              string            `json:"url"`              // 클릭하면 이동되는 URL
	ClassNames       []string          `json:"classNames"`       // html에 렌더링할 때 attached 할 클레스 이름
	Editable         bool              `json:"editorble"`        // 수정 여부
	StartEditable    bool              `json:"startEditable"`    // 시작시간 수정 여부
	EndEditable      bool              `json:"endEditable"`      // 끝시간 수정 여부
	DurationEditable bool              `json:"durationEditable"` // 기간수정 여부
	ResourceEditable bool              `json:"resourceEditable"` // 리소스 수정가능 여부
	Display          string            `json:"display"`          // 'auto','block','list-item','background','inverse-background','none' -> https://fullcalendar.io/docs/eventDisplay
	BackgroundColor  string            `json:"backgroundColor"`  // 배경색
	BorderColor      string            `json:"borderColor"`      // 테두리색
	TextColor        string            `json:"textColor"`        // 글자색
	ExtendedProps    map[string]string `json:"extendedProps"`    // 나머지 필요한 기능
	// Constraint // 일정을 설정할 때 제한을 둘때 사용한다.
	// Overlap // false 또는 함수를 다룬다. Go 자료구조에서는 다루지 않는다.
}

// https://fullcalendar.io/docs/resource-object
type FullCalendarResource struct {
	ID                   string            `json:"id"`                   // ID
	Title                string            `json:"title"`                // 제목
	ExtendedProps        map[string]string `json:"extendedProps"`        // 나머지 필요한 기능
	EventBackgroundColor string            `json:"eventBackgroundColor"` // 배경색
	EventBorderColor     string            `json:"eventBorderColor"`     // 테두리색
	EventTextColor       string            `json:"eventTextColor"`       // 글자색
	EventClassNames      []string          `json:"eventClassNames"`      // html에 렌더링할 때 attached 할 클레스 이름
	// EventAllow
	// EventConstraint // 일정을 설정할 때 제한을 둘때 사용한다.
	// EventOverlap // false 또는 함수를 다룬다. Go 자료구조에서는 다루지 않는다.
}
