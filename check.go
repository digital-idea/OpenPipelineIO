package main

import (
	"os"
	"regexp"
	"strings"
)

// 해상도값 정규식 : 2048x1080 형태
var regexpImageSize = regexp.MustCompile(`\d{2,5}[xX]\d{2,5}$`)

// 샷네임값 정규식 : SS_0010 형태
var regexpShotname = regexp.MustCompile(`^[a-zA-Z0-9]+_[a-zA-Z0-9]+$`)

// 에셋 네임값 정규식 : stone01 형태
var regexpAssetname = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// 롤미디어 정규식 : 00_A03C001_180113_A001 형태
var regexpRollMedia = regexp.MustCompile(`^\d+_[A-Z0-9]+_\d+_[A-Z0-9]+$`)

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

// str2bool은 받아들인 문자열이 "true"나 "True"라면 참을, 그 외에는 거짓을 반환한다.
func str2bool(str string) bool {
	if str == "true" || str == "True" {
		return true
	}
	return false
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
func validAssettype(assettype string) bool {
	switch assettype {
	case "char", "comp", "env", "fx", "global", "matte", "plant", "prop", "vehicle", "concept":
		return true
	}
	return false
}
