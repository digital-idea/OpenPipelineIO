package main

import "errors"

// PublishKey 자료구조는 상태 자료구조이다.
type PublishKey struct {
	ID          string `json:"id"`          // ID
	Description string `json:"description"` // 설명
}

// CheckError 메소드는 PublishKey 자료구조의 에러를 체크한다.
func (key *PublishKey) CheckError() error {
	if key.ID == "" {
		return errors.New("ID가 빈 문자열 입니다")
	}
	return nil
}
