package main

// OnsetCam 자료구조는 현장카메라에 대한 자료구조이다.
type OnsetCam struct {
	CameraName     string  `json:"cameraname"`     // 카메라명 예)Arri Alexa XT+
	LensName       string  `json:"lensname"`       // 렌즈명 Arri Master Prime
	DataName       string  `json:"dataname"`       // 샷의 데이터네임 예) S#53 C#2 T#2 / A110-C009
	RigName        string  `json:"rigname"`        // 리그 이름
	Lensmm         float64 `json:"lensmm"`         // 렌즈mm
	SerialNum      string  `json:"serialnum"`      // 렌즈 시리얼넘버. 입체시 Left
	SerialNumRight string  `json:"serialnumright"` // 렌즈 시리얼넘버. Right
	Iso            int     `json:"iso"`            // 800
	Temp           int     `json:"temp"`           // 색온도 예)4800K
	Unit           string  `json:"unit"`           // A,B,C 캠
	StereoIo       float64 `json:"stereoio"`       // 입체리그-카메라 거리, 만약 다이나믹 IO라면 평균값이 들어간다.
	StereoConv     float64 `json:"stereoconv"`     // 입체리그-Convergence 거리 : 각도구하는 법 = degree(atan(2*conv/io))
	StereoUnit     string  `json:"stereounit"`     // 입체리그-단위 : mm, feet
	StereoMaincam  string  `json:"stereomaincam"`  // 입체리그-기준카메라 위치 예) left, right
}

// ProductionCam 자료구조는 CG작업에 사용하는 카메라 구조이다.
type ProductionCam struct {
	PubPath    string `json:"pubpath"`    // 퍼블리쉬 카메라 경로
	PubTask    string `json:"pubtask"`    // 카메라 퍼블리쉬 Task : mm,layout,ani
	Projection bool   `json:"projection"` // 프로젝션을 사용하는 카메라인가?
	Lensmm     string `json:"lensmm"`     // 렌즈mm
}
