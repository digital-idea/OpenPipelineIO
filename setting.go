package main

// Setting 자료구조는 관리자 설정 자료구조이다.
type Setting struct {
	ID                            string `json:"id"`                            // 셋팅ID
	RunScriptAfterSignup          string `json:"runscriptaftersignup"`          // 사용자 가입이후 실행될 쉘스크립트
	RunScriptAfterEditUserProfile string `json:"runscriptafteredituserprofile"` // 사용자 정보 수정후 실행될 쉘스크립트
	ExcludeProject                string `json:"excludeproject"`                // Search옵션에 제외할 프로젝트명, 마이그레이션 시 사용한다.
	OCIOConfig                    string `json:"ocioconfig"`                    // OpenColorIO Config Path 설정
}
