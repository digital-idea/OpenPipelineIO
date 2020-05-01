package main

// Review 는 리뷰데이터 자료구조 이다.
type Review struct {
	Project    string    `json:"project"`    // 프로젝트
	Task       string    `json:"task"`       // 태스크
	Createtime string    `json:"createtime"` // 생성시간
	Author     string    `json:"author"`     // 작성자
	Path       string    `json:"path"`       // 리뷰경로
	Comments   []Comment `json:"comments"`   // 댓글
}
