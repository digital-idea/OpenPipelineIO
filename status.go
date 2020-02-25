package main

import "errors"

// Status 자료구조는 상태 자료구조이다.
type Status struct {
	ID          string  `json:"id"`          // ID, 상태코드
	Color       string  `json:"color"`       // 상태 색상
	Description string  `json:"description"` // 설명
	Order       float64 `json:"order"`       // Status 우선순위
}

// CheckError 메소드는 Status 자료구조의 에러를 체크한다.
func (s *Status) CheckError() error {
	if s.ID == "" {
		return errors.New("ID가 빈 문자열 입니다")
	}
	if !regexWebColor.MatchString(s.Color) {
		return errors.New("웹컬러 문자가 아닙니다")
	}
	return nil
}
