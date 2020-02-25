package main

// StatusV2 자료구조는 상태 자료구조이다.
// 기존에 Status 이름을 사용중이다. 새로운 Status는 StatusV2로 사용한다.
type StatusV2 struct {
	ID          string `json:"id"`          // ID, 상태코드
	Color       string `json:"color"`       // 상태 색상
	Description string `json:"description"` // 설명
	Order       int    `json:"order"`       // Status 우선순위
}
