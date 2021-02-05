package main

// Setting 자료구조는 관리자 설정 자료구조이다.
type Setting struct {
	ID                            string `json:"id"`                            // 셋팅ID
	RunScriptAfterSignup          string `json:"runscriptaftersignup"`          // 사용자 가입이후 실행될 쉘스크립트
	RunScriptAfterEditUserProfile string `json:"runscriptafteredituserprofile"` // 사용자 정보 수정후 실행될 쉘스크립트
	ExcludeProject                string `json:"excludeproject"`                // Search옵션에 제외할 프로젝트명, 마이그레이션 시 사용한다.
	OCIOConfig                    string `json:"ocioconfig"`                    // OpenColorIO Config Path 설정
	FFmpeg                        string `json:"ffmpeg"`                        // FFmpeg 경로 셋팅
	FFmpegThreads                 int    `json:"ffmpegthreads"`                 // FFmpeg 연산 Thread 셋팅
	RVPath                        string `json:"rvpath"`                        // RV 경로 셋팅
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
	AssetTypePath                 string `json:"assettypepath"`                 // Asset Type 경로 예) /show/{{.Project}}/assets/{{.Assettype}}
	AssetTypePathPermission       string `json:"assettypepathpermission"`       // Asset Type 경로의 권한
	AssetTypePathUID              string `json:"assettypepathuid"`              // Asset Type 경로의 User ID
	AssetTypePathGID              string `json:"assettypepathgid"`              // Asset Type 경로의 Group ID
	AssetPath                     string `json:"assetpath"`                     // 개별 Asset 예) /show/{{.Project}}/assets/{{.Assettype}}/{{.Name}}
	AssetPathPermission           string `json:"assetpathpermission"`           // 개별 Asset 경로의 권한
	AssetPathUID                  string `json:"assetpathuid"`                  // 개별 Asset 경로의 User ID
	AssetPathGID                  string `json:"assetpathgid"`                  // 개별 Asset 경로의 Group ID
	ThumbnailRootPath             string `json:"thumbnailrootpath"`             // 썸네일 이미지 Root 경로
	ThumbnailRootPathPermission   string `json:"thumbnailrootpathpermission"`   // 썸네일 이미지 Root 경로의 권한
	ThumbnailRootPathUID          string `json:"thumbnailrootpathuid"`          // 썸네일 이미지 Root 경로의 User ID
	ThumbnailRootPathGID          string `json:"thumbnailrootpathgid"`          // 썸네일 이미지 Root 경로의 Group ID
	ThumbnailImagePath            string `json:"thumbnailimagepath"`            // 썸네일 이미지 경로 /thumbnail/{{.Project}}/{{.Name}}_{{.Type}}
	ThumbnailImagePathPermission  string `json:"thumbnailimagepathpermission"`  // 썸네일 이미지 경로의 권한
	ThumbnailImagePathUID         string `json:"thumbnailimagepathuid"`         // 썸네일 이미지 경로의 User ID
	ThumbnailImagePathGID         string `json:"thumbnailimagepathgid"`         // 썸네일 이미지 경로의 Group ID
	ThumbnailMovPath              string `json:"thumbnailmovpath"`              // 썸네일 동영상 경로
	ThumbnailMovPathPermission    string `json:"thumbnailmovpathpermission"`    // 썸네일 동영상 경로의 권한
	ThumbnailMovPathUID           string `json:"thumbnailmovpathuid"`           // 썸네일 동영상 경로의 User ID
	ThumbnailMovPathGID           string `json:"thumbnailmovpathgid"`           // 썸네일 동영상 경로의 Group ID
	PlatePath                     string `json:"platepath"`                     // Plate 경로
	PlatePathPermission           string `json:"platepathpermission"`           // Plate 경로의 권한
	PlatePathUID                  string `json:"platepathuid"`                  // Plate 경로의 User ID
	PlatePathGID                  string `json:"platepathgid"`                  // Plate 경로의 Group ID
	ReviewDataPath                string `json:"reviewdatapath"`                // 리뷰 데이터가 저장되는 경로
	ReviewDataPathPermission      string `json:"reviewdatapathpermission"`      // 리뷰 데이터가 저장되는 경로의 권한
	ReviewDataPathUID             string `json:"reviewdatapathuid"`             // 리뷰 데이터가 저장되는 경로의 User ID
	ReviewDataPathGID             string `json:"reviewdatapathgid"`             // 리뷰 데이터가 저장되는 경로의 Group ID

	DefaultScaleRatioOfUndistortionPlate float64 `json:"defaultscaleratioofundistortionplate"` // 언디스토션 플레이트의 기본 스케일비율 1.1배
	ItemNumberOfPage                     int     `json:"itemnumberofpage"`                     // 한 페이지에 보이는 아이템 갯수
	MultipartFormBufferSize              int     `json:"multipartformbuffersize"`              // Multipart form buffer size
	ThumbnailImageWidth                  int     `json:"thumbnailimagewidth"`                  // Thumbnail Image 가로사이즈
	ThumbnailImageHeight                 int     `json:"thumbnailimageheight"`                 // Thumbnail Image 세로사이즈
}
