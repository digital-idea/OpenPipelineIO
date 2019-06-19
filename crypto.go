package main

import (
	"crypto/md5"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Passcheck 함수는 사용 가능한 패스워드인지 체크하는 함수이다.
// MPAA의 기본조건은 7자리 이상 길이의 패스워드이다.
// 또한 MPAA 규칙중 10자리이상, 대문자, 소문자, 특수문자, 숫자, 연속된 문자금지 조항중 조건이 4개이상 일치해야한다.
func Passcheck(pw string) bool {
	if len(pw) <= 7 {
		return false
	}
	hasLong, hasUpper, hasLower, hasSpecial, hasNum, hasNotSame := false, false, false, false, false, true
	if len(pw) >= 10 { // 패스워드가 10자리 이상을 만족함.
		hasLong = true
	}
	var before rune
	for _, r := range pw {
		if before == r { // 전 문자열과 비교한다.
			hasNotSame = false
		}
		if strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ", string(r)) {
			hasUpper = true
		}
		if strings.Contains("abcdefghijklmnopqrstuvwxyz", string(r)) {
			hasLower = true
		}
		if strings.Contains("0123456789", string(r)) {
			hasNum = true
		}
		if strings.Contains("~`!@#$%^&*()_-+={[}]|:;\"'<,>.?/", string(r)) {
			hasSpecial = true
		}
		before = r
	}
	condition := 0
	conditions := []bool{hasLong, hasUpper, hasLower, hasSpecial, hasNum, hasNotSame}
	for _, c := range conditions {
		if c {
			condition++
		}
	}
	if condition < 4 {
		return false
	}
	return true
}

// Encrypt 함수는 문자를 받아서 해쉬문자로 변환한다.
func Encrypt(s string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Str2md5 함수는 문자를 받아서 md5로 암호화된 문자를 반환한다.
func Str2md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
