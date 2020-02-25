package main

// Status 자료구조는 상태 자료구조이다.
type Status struct {
	ID          string `json:"id"`          // ID, 상태코드
	Color       string `json:"color"`       // 상태 색상
	Description string `json:"description"` // 설명
	Order       int    `json:"order"`       // Status 우선순위
}
