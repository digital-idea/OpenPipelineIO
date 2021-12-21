package main

import "strings"

// replaceSearchFilder 함수는 해당 값에서 key가 존재 한다면, value를 교체한다. 없으면 추가한다.
func setSearchFilter(org, key, value string) string {
	var newWords []string
	// key가 없다면 포함시킨다.
	if !strings.Contains(org, key) && value == "" {
		newWords = append(newWords, key+":"+value)
	}
	words := strings.Split(org, " ")
	for _, word := range words {
		if strings.HasPrefix(word, key+":") {
			// 키가 존재하면 값을 바꾼다.
			newWords = append(newWords, key+":"+value)
		} else {
			newWords = append(newWords, word)
		}

	}
	return strings.Join(newWords, " ")
}
