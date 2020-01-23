package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// 해상도값 정규식: 2048x1080 형태
var regexpImageSize = regexp.MustCompile(`\d{2,5}[xX]\d{2,5}$`)

// 샷네임값 정규식: SS_0010 형태
var regexpShotname = regexp.MustCompile(`^[a-zA-Z0-9]+_[a-zA-Z0-9]+$`)

// 에셋 네임값 정규식: stone01 형태
var regexpAssetname = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// ID값 정규식: organization
var regexpID = regexp.MustCompile(`^[a-z0-9]+$`)

// userInfo 정규칙: id(name,team) 형태
var regexpUserInfo = regexp.MustCompile(`^[a-zA-Z0-9]+\(.+\)$`)

// Task값 정규식: task
var regexpTask = regexp.MustCompile(`^[a-z0-9]+$`)

// 롤미디어 정규식: 00_A03C001_180113_A001 형태
var regexpRollMedia = regexp.MustCompile(`^\d+_[A-Z0-9]+_\d+_[A-Z0-9]+$`)

// Timecode 정규식: 00:00:00:00 또는 00:00:00;00 형태 ref: https://en.wikipedia.org/wiki/SMPTE_timecode
var regexpTimecode = regexp.MustCompile(`^\d{2}[:;.]\d{2}[:;.]\d{2}[:;.]\d{2}$`)

// Rnum 정규식: A0001~H9999
var regexpRnum = regexp.MustCompile(`^[A-H]\d{4}$`)

// Handle 정규식: 5, 10
var regexpHandle = regexp.MustCompile(`^\d{1,2}$`)

// Version 정규식: 5, 10, 103, v5 v10, v100, v1001, V001, V1
var regexpVersion = regexp.MustCompile(`[vV]?\d{1,4}`)

// Alexa 카메라의 형태 : N_AAAACCCC_YYMMDD_RRRR
// - N : order
// - AAAACCCC : reel name
// - YYMMDD : 년,월,일
// - RRRR : Unique hex code
//   -  참고 : embedded 네임은 AAAARRRR 형태이다.
//
// Red카메라의 형태 : C RRR TTT dd mm HH NNN
// - C : Camera letter, from A to Z
// - RRR : three digit reel number
// - TTT : three digit take number
// - dd mm : date and  month
// - HH : unique hex code that ensures another file won't have the exact same name
// - NNN : spanned clip number

// str2bool은 받아들인 문자열이 "true","True" ... 라면 참을, 그 외에는 거짓을 반환한다.
func str2bool(str string) bool {
	if strings.ToLower(str) == "true" {
		return true
	}
	return false
}

func bool2str(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// mov경로를 체크하는 함수이다.
func isMov(path string) bool {
	// 빈 문자열인지 체크
	if path == "" {
		return false
	}
	// .mov로 끝나는지 체크
	if !strings.HasSuffix(strings.ToLower(path), ".mov") {
		return false
	}
	// 경로가 존재하는지 체크
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// 에셋타입이 유효한지 체크하는 함수이다.
func validAssettype(typ string) (bool, error) {
	switch typ {
	case "char", "comp", "env", "fx", "matte", "plant", "prop", "vehicle", "concept":
		return true, nil
	default:
		return false, fmt.Errorf("%s 이름을 에셋타입으로 사용할 수 없습니다", typ)
	}
}

// 롤넘버가 유효한지 체크하는 함수이다.
func validRnumTag(rnum string) bool {
	switch rnum {
	case "1권", "2권", "3권", "4권", "5권", "6권", "7권", "8권":
		return true
	default:
		return false
	}
}

// validShottype 함수는 샷타입이 유효한지 체크하는 함수이다.
func validShottype(typ string) error {
	switch typ {
	case "2d", "3d", "":
		return nil
	default:
		return fmt.Errorf("%s 이름을 샷타입으로 사용할 수 없습니다", typ)
	}
}

// UniqueSlice 함수는 중복 문자를 제거한다.
func UniqueSlice(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// isASCII 함수는 문자가 아스키로 구성되어있는지 체크한다.
func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}
