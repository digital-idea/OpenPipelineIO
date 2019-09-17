package main

// AdminSetting 자료구조는 관리자 설정 자료구조이다.
type AdminSetting struct {
	RunScriptAfterSignup string `json:"runscriptaftersignup"` // 사용자 가입이후 실행될 쉘스크립트
}
