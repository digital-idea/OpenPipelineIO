package main

import "errors"

// Status 자료구조는 상태 자료구조이다.
type Status struct {
	ID                string  `json:"id"`                // ID, 상태코드
	TextColor         string  `json:"textcolor"`         // TEXT 색상
	BGColor           string  `json:"bgcolor"`           // BG 상태 색상
	BorderColor       string  `json:"bordercolor"`       // Border 색상
	Description       string  `json:"description"`       // 설명
	Order             float64 `json:"order"`             // Status 우선순위
	DefaultOn         bool    `json:"defaulton"`         // 검색바 기본선택 여부
	InitStatus        bool    `json:"initstatus"`        // 아이템 생성시 최초 설정되는 Status 설정값
	ReviewStatusEvent string  `json:"reviewstatusevent"` // 리뷰상태 변경시 발생하는 이벤트
	ProgressCategory  string  `json:"progressgategory"`  // 상태에 대한 큰 카테고리, In Progress, Final Approved, On Hold 같은 개념정의
}

// CheckError 메소드는 Status 자료구조의 에러를 체크한다.
func (s *Status) CheckError() error {
	if s.ID == "" {
		return errors.New("ID가 빈 문자열 입니다")
	}
	if !regexpStatus.MatchString(s.ID) {
		return errors.New("status 이름은 영문 대,소문자 또는 숫자로만 이루어져야 합니다")
	}
	if !regexWebColor.MatchString(s.TextColor) {
		return errors.New("text 컬러가 웹컬러(#FFFFFF 형태) 문자열이 아닙니다")
	}
	if !regexWebColor.MatchString(s.BGColor) {
		return errors.New("BG 컬러가 웹컬러(#FFFFFF 형태) 문자열이 아닙니다")
	}
	if !regexWebColor.MatchString(s.BorderColor) {
		return errors.New("테두리 컬러가 웹컬러(#FFFFFF 형태) 문자열이 아닙니다")
	}
	return nil
}
