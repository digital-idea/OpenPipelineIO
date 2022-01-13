package main

// Comment 자료구조는 코맨트 글을 작성할 때 사용하는 자료구조이다.
type Comment struct {
	Date       string `json:"date"`       // 코맨트 등록시간 RFC3339
	Author     string `json:"author"`     // 작성자 ID
	AuthorName string `json:"authorname"` // 작성자 표기명
	Text       string `json:"text"`       // 내용
	Media      string `json:"media"`      // media 경로
	MediaTitle string `json:"mediatitle"` // media 제목
	Stage      string `json:"stage"`      // 코멘트가 달린 리뷰 Stage
	ItemStatus string `json:"itemstatus`  // 코멘트가 달렸던 ItemStatus
	Frame      int    `json:"frame"`      // 코맨트 프레임수
}
