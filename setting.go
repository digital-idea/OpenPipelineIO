package main

// Setting 자료구조는 관리자 설정 자료구조이다.
type Setting struct {
	ID                            string `json:"id"`                            // 셋팅ID
	RunScriptAfterSignup          string `json:"runscriptaftersignup"`          // 사용자 가입이후 실행될 쉘스크립트
	RunScriptAfterEditUserProfile string `json:"runscriptafteredituserprofile"` // 사용자 정보 수정후 실행될 쉘스크립트
	ExcludeProject                string `json:"excludeproject"`                // Search옵션에 제외할 프로젝트명, 마이그레이션 시 사용한다.
	OCIOConfig                    string `json:"ocioconfig"`                    // OpenColorIO Config Path 설정
	Umask                         string `json:"umask"`                         // Umask 값. 예) 0002
	RootPath                      string `json:"rootpath"`                      // Root경로 예) /show
	ProjectPath                   string `json:"projectpath"`                   // Project경로 예) /show/{{.Project}}
	ProjectPathPermission         string `json:"projectpathpermission"`         // Project경로의 권한
	ProjectPathUID                string `json:"projectpathuid"`                // Project경로의 User ID
	ProjectPathGID                string `json:"projectpathgid"`                // Project경로의 Group ID
	ShotRootPath                  string `json:"shotrootpath"`                  // Shot Root 경로 예) /show/{{.Project}}/seq/
	ShotRootPathPermission        string `json:"shotrootpathpermission"`        // Shot Root 경로의 권한
	ShotRootPathUID               string `json:"shotrootpathuid"`               // Shot Root 경로의 User ID
	ShotRootPathGID               string `json:"shotrootpathgid"`               // Shot Root 경로의 Group ID
	SeqPath                       string `json:"seqpath"`                       // Seq 경로 예) /show/{{.Project}}/seq/{{.Seq}}
	SeqPathPermission             string `json:"seqpathpermission"`             // Seq 경로의 권한
	SeqPathUID                    string `json:"seqpathuid"`                    // Seq 경로의 User ID
	SeqPathGID                    string `json:"seqpathgid"`                    // Seq 경로의 Group ID
	ShotPath                      string `json:"shotpath"`                      // 개별 Shot 경로 예) /show/{{.Project}}/seq/{{.Seq}}/{{.Name}}
	ShotPathPermission            string `json:"shotpathpermission"`            // 개별 Shot 경로의 권한
	ShotPathUID                   string `json:"shotpathuid"`                   // 개별 Shot 경로의 User ID
	ShotPathGID                   string `json:"shotpathgid"`                   // 개별 Shot 경로의 Group ID
	AssetRootPath                 string `json:"assetrootpath"`                 // Asset Root 경로 예) /show/{{.Project}}/assets/
	AssetRootPathPermission       string `json:"assetrootpathpermission"`       // Asset Root 경로의 권한
	AssetRootPathUID              string `json:"assetrootpathuid"`              // Asset Root 경로의 User ID
	AssetRootPathGID              string `json:"assetrootpathgid"`              // Asset Root 경로의 Group ID
	AssetPath                     string `json:"assetpath"`                     // 개별 Asset 예) /show/{{.Project}}/assets/{{.Assettype}}/{{.Name}}
	AssetPathPermission           string `json:"assetpathpermission"`           // 개별 Asset 경로의 권한
	AssetPathUID                  string `json:"assetpathuid"`                  // 개별 Asset 경로의 User ID
	AssetPathGID                  string `json:"assetpathgid"`                  // 개별 Asset 경로의 Group ID
}
