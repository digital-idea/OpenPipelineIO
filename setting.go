package main

// Setting 자료구조는 관리자 설정 자료구조이다.
type Setting struct {
	ID                            string `json:"id"`                            // 셋팅ID
	RunScriptAfterSignup          string `json:"runscriptaftersignup"`          // 사용자 가입이후 실행될 쉘스크립트
	RunScriptAfterEditUserProfile string `json:"runscriptafteredituserprofile"` // 사용자 정보 수정후 실행될 쉘스크립트
	ExcludeProject                string `json:"excludeproject"`                // Search옵션에 제외할 프로젝트명, 마이그레이션 시 사용한다.
	OCIOConfig                    string `json:"ocioconfig"`                    // OpenColorIO Config Path 설정
	LinuxRootPath                 string `json:"linuxrootpath"`                 // Root경로 예) /show
	LinuxProjectPath              string `json:"linuxprojectpath"`              // Project경로 예) /show/{{.Project}}
	LinuxShotPath                 string `json:"linuxshotpath"`                 // Shot경로 예) /show/{{.Project}}/seq/{{.Seq}}/{{.Name}}
	LinuxAssetPath                string `json:"linuxassetpath"`                // Asset경로 예) /show/{{.Project}}/assets/{{.Assettype}}/{{.Name}}
	UID                           string `json:"uid"`                           // User ID 예) 500, /etc/passwd 에 선언, 터미널에서 "$ id -u" 의 결과
	GID                           string `json:"gid"`                           // Group ID 예) 500, /etc/group 에 선언, 터미널에서 "$ id -g" 의 결과
	Permission                    string `json:"permission"`                    // Permission 숫자 예) 0775
	Umask                         string `json:"umask"`                         // Umask 값. 예) 0002
}
