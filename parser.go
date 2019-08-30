package main

import (
	"log"
	"net/url"
	"strings"
)

// Searchword 는 검색어를 종류별로 담는 자료구조이다.
type Searchword struct {
	words  []string // 아래 항목에 들어가지 않는 나머지 단어의 집합
	tasks  []string // task:ani,comp
	user   []string // user:김한웅 형태로 시작하는 문자
	tags   []string // tag:태그명 형태로 시작하는 문자
	status []string // status: 형태로 시작하는 문자
}

// SearchwordParser 함수는 단어를 받아서 Searchword 자료구조로 반환한다.
func SearchwordParser(word string) Searchword {
	searchword := Searchword{}
	for _, w := range strings.Split(word, " ") {
		// 공백(스페이스) 체크
		if strings.TrimSpace(w) == "" {
			continue
		}

		switch {
		case strings.HasPrefix(w, "task:"): //'task:'로 시작하는지 체크
			// taskValue := strings.TrimLeft(w, "task:")
			// string.TrimLeft를 쓰는 경우 task:ani 등의 :a로 시작하면 a문자열이 잘린다. 유니코드문제
			taskValue := strings.TrimSpace(strings.Split(w, ":")[1])
			taskList := strings.Split(taskValue, ",")
			for _, list := range taskList {
				for _, task := range TASKS {
					if strings.ToLower(list) != strings.ToLower(task) {
						continue
					}
					searchword.tasks = append(searchword.tasks, strings.ToLower(list))
					break
				}
			}
		case strings.HasPrefix(w, "user:"): // 'user:'로 시작하는지 체크
			searchword.user = append(searchword.user, strings.TrimLeft(w, "user:"))
		case strings.HasPrefix(w, "tag:"): // 'tag:'로 시작하는지 체크
			searchword.tags = append(searchword.tags, strings.TrimLeft(w, "tag:"))
		case strings.HasPrefix(w, "status:"): // 'status:'로 시작하는지 체크
			searchword.status = append(searchword.status, strings.TrimLeft(w, "status:"))
		default:
			//중복 검색어 체크
			hasWord := false
			if len(searchword.words) >= 1 {
				for _, word := range searchword.words {
					if word == strings.TrimSpace(w) {
						hasWord = true
						break
					}
				}
				if hasWord {
					continue
				}
			}
			searchword.words = append(searchword.words, strings.TrimSpace(w))
		}
	}
	// 서치 데이터에 값이 없으면 오류가 발생하므로 words에 값이 없다면 ''값을 넣어준다.
	if len(searchword.words) == 0 {
		searchword.words = append(searchword.words, "''")
	}
	return searchword
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
