package main

import (
	"fmt"
	"log"
	"time"
)

func str2time(str string) time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		log.Println("시간을 바꿀 수 없습니다.")
		//리턴은 1970년 초기값을 반환한다.
	}
	return t
}

// ToDday 함수는 RFC3339 날짜를 받아서 D-100 형태의 문자를 출력한다.
func ToDday(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	deadline, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return str, err
	}
	//days := deadline.Day() - time.Now().Day()
	days := time.Now().Sub(deadline).Hours() / 24
	sign := "+"
	if days < 0 {
		sign = ""
	}
	return fmt.Sprintf("D%s%d", sign, int(days)), nil
}

func shortTimecode(timecode string) (string, error) {
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return "", fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	return timecode[0:8], nil
}
