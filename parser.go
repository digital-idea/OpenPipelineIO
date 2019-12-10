package main

import (
	"log"
	"net/url"
)

// Searchword 는 검색어를 종류별로 담는 자료구조이다.
type Searchword struct {
	words  []string // 아래 항목에 들어가지 않는 나머지 단어의 집합
	tasks  []string // task:ani,comp
	users  []string // user:김한웅 형태로 시작하는 문자
	tags   []string // tag:태그명 형태로 시작하는 문자
	status []string // status: 형태로 시작하는 문자
}

// URLUnescape 함수는 encode된 URL주소를 decode 한다.
func URLUnescape(u *url.URL) (url.Values, error) {
	unescape, err := url.QueryUnescape(u.RawQuery)
	if err != nil {
		log.Println(err)
	}
	q, err := url.ParseQuery(unescape)
	if err != nil {
		return q, err
	}
	return q, nil
}
