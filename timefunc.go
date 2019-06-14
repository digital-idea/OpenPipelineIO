package main

import (
	"log"
	"time"
)

// Now 함수는 현재시간을 RFC3339로 반환한다.
func Now() string {
	return time.Now().Format(time.RFC3339)
}

func str2time(str string) time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		log.Println("시간을 바꿀 수 없습니다.")
		//리턴은 1970년 초기값을 반환한다.
	}
	return t
}
