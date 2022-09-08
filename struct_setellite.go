package main

// 아이패드용 현장툴인 setellite 어플리케이션 CSV 자료구조이다.
// mongoDB구조 : setellite.project
// 홈페이지 https://www.setellite.nl

import (
	"reflect"
	"strings"
)

// Setellite CSV 포멧양식이다.
// 입력의 예는 영화 "상해보루" 기록에 직접 사용된 예시를 작성했다.
// Setellite 포멧에서 , 문자는 다음데이터를 뜻한다.
// Setellite 포멧에서 ; 문자는 여러 정보를 묶을 때 사용한다.
// 자료구조의 key는 Setellite의 " ",'/','[',']','-'를 제외한 영문명을 그대로 적용했다.
// CSV내부에 Notes값은 여러곳에 분포되어있고 Table Key가 중복된다. Notes는 자료구조에서는 다른 자료와 같은방식의 ";"로 구분된 string으로 표현한다.
// 데이터의 중복을 막기위해 ID를 사용한다. 조합방법: ShootDay + SlateNumber + TakeNumber + RollMedia
// 데이터가 크기때문에 OpenPipelineIO 샷정보에 링크하고, OpenPipelineIO 아이템 자료구조에 바로넣지 않는다.
type Setellite struct {
	ID                    string `json:"id"`                    // DB에 존재하는지 체크값으로 사용한다. ShootDay+SlateNumber + TakeNumber +  RollMedia
	Project               string `json:"project"`               // 프로젝트명
	SlateNumber           string `json:"slatenumber"`           // 슬레이트 넘버. 예) 1
	Episode               string `json:"episode"`               // 에피소드명.
	SceneNumber           string `json:"scenenumber"`           // 씬넘버. 예) X182
	StoryboardNumber      string `json:"storyboardnumber"`      // 스토리보드 넘버. 예) Shot#1
	Unit                  string `json:"unit"`                  // 촬영유닛. 예) B
	Timestamp             string `json:"timestamp"`             // 기록시작시간. 예) 02-10-2017 07:33:53
	ShootDay              string `json:"shootday"`              // 촬영일. 예) 1
	VFXShot               string `json:"vfxshot"`               // 작업내용. 예) 포식자 합성
	VFXWork               string `json:"vfxwork"`               // CG요소. 예) 3D character;wire removal;스파크 심기
	InteriorExterior      string `json:"interiorexterior"`      // 실내,실외체크 예) [INT | EXT]
	ScriptTime            string `json:"scripttime"`            // 낮,밤체크 예) [Day|Dusk|Night|Dawn]
	ScriptLocation        string `json:"scriptlocation"`        // 촬영장소정보 예) Holy vine camp
	ShotDescription       string `json:"shotdescription"`       // 샷 설명
	SetLocation           string `json:"setlocation"`           // 장소. 예) 공장
	ElementType           string `json:"elementtype"`           // 촬영요소 타입. 예) CG/3D Elements; Reference Material,
	SetMedia              string `json:"setmedia"`              // 연관된 미디어. 예) Balls; Clean Plate; Color Chart
	Weather               string `json:"weather"`               // 날씨. 예) Rain
	Wrangler              string `json:"wrangler"`              // 정보수집자. 예)Digital Idea 또는 작성자명.
	HDRiFilename          string `json:"hdrifilename"`          // HDRi파일명
	CameraName            string `json:"cameraname"`            // 카메라이름. 예) Camera E
	CameraGrip            string `json:"cameragrip"`            // 그립. 예) Handheld
	CameraMovement        string `json:"cameramovement"`        // 카메라 이동형태. 예) Free Move
	CameraModel           string `json:"cameramodel"`           // 카메라 모델명. 예) Arri Alexa M
	FormatStock           string `json:"formatstock"`           // 포멧. 예) ARRIRAW 12-bit
	ISORating             string `json:"isorating"`             // 감도. 예) ISO 800
	ColorTemperature      string `json:"colortemperature"`      // 색온도. 예) 5600K
	Resolution            string `json:"resolution"`            // 해상도. 예) 3.4K
	AspectRatio           string `json:"aspectratio"`           // 렌즈비. 예) 2:1 (아나모픽렌즈경우)
	PixelAspectRatio      string `json:"pixelaspectratio"`      // 픽셀비. 예) 1:1 (square)
	CompressionRatio      string `json:"compressionratio"`      // 압축비.
	CameraHead            string `json:"camerahead"`            // 카메라 헤드형태기록
	LensModel             string `json:"lensmodel"`             // 렌즈명. 예) Arri Master Primes 21mm
	LensFocalLength       string `json:"lensfocallength"`       // 렌즈mm. 예) 21 mm
	LensType              string `json:"lenstype"`              // 렌즈타입. 예) Prime Lens
	StereoCameraModel     string `json:"stereocameramodel"`     // 입체카메라 모델
	StereoRigOrientation  string `json:"stereorigorientation"`  // 입체카메라 리그의 성향
	StereoThruCamEye      string `json:"stereothrucameye"`      // 입체카메라 기존화면 [left|right]
	StereoCameraHead      string `json:"stereocamerahead"`      // 입체카메라 헤드
	StereoLensModel       string `json:"stereolensmodel"`       // 입체 렌즈모델
	StereoLensFocalLength string `json:"setreolensfocallength"` // 입체카메라 렌즈mm
	StereoLensType        string `json:"stereolenstype"`        // 스테레오 렌즈타입
	Notes                 string `json:"notes"`                 // 카메라 관련 메모. CSV에서 Notes는 여러값을 가진다.
	TakeNumber            string `json:"takenumber"`            // 테이크넘버. 예) 1
	TakeName              string `json:"takename"`              // 테이크이름. 예) [Rehearsal|OK|Keep|Clean plate|Ref]
	TimestampStart        string `json:"timestampstart"`        // 테이크시작시간. 예) 02-10-2017 15:49:45
	TimestampStop         string `json:"timestampstop"`         // 테이크종료시간. 예) 15:50:00
	RollMedia             string `json:"rollmedia"`             // 롤미디어이름.중요! 예) E002C001
	LensFocalLengthStart  string `json:"lensfocallengthstart"`  // 렌즈미리수(시작점). 예) 21
	LensFocalLengthStop   string `json:"lensfocallengthstop"`   // 렌즈미리수(끝점).
	LensHeight            string `json:"lensheight"`            // 렌즈높이
	LensHeightStart       string `json:"lensheightstart"`       // 렌즈높이(시작점)
	LensHeightStop        string `json:"lensheightstop"`        // 렌즈높이(끝점)
	LensHeightUnit        string `json:"lensheightunit"`        // 렌즈높이단위. 예) cm
	LensFocus             string `json:"lensfocus"`             // 렌즈포커스
	LensFocusStart        string `json:"lensfocusstart"`        // 렌즈포커스(시작점)
	LensFocusStop         string `json:"lensfocusstop"`         // 렌즈포커스(끝점)
	LensFocusUnit         string `json:"lensfocusunit"`         // 카메라포커스 단위. 예) cm
	SubjectDistance       string `json:"subjectdistance"`       // 피사체거리
	SubjectDistanceStart  string `json:"subjectdistancestart"`  // 피사체거리(시작점)
	SubjectDistanceStop   string `json:"subjectdistancestop"`   // 피사체거리(끝점)
	SubjectDistanceUnit   string `json:"subjectdistanceunit"`   // 피사체거리 단위. 예) cm
	CameraFStop           string `json:"camerafstop"`           // 카메라 조리개
	CameraFStopStart      string `json:"camerafstopstart"`      // 카메라 조리개(시작점) F-Stop 예) 6
	CameraFStopStop       string `json:"camerafstopstop"`       // 카메라 조리개(끝점)
	CameraTilt            string `json:"cameratilt"`            // 카메라 틸트(위,아래움직임)
	CameraTiltStart       string `json:"cameratiltstart"`       // 카메라 틸트(시작점)
	CameraTiltStop        string `json:"cameratiltstop"`        // 카메라 틸트(끝점)
	CameraDutch           string `json:"cameradutch"`           // 카메라 더치틸트(좌,우로 회전시키는 움직임)
	Framerate             string `json:"framerate"`             // FPS. 예) 24
	ShutterAngle          string `json:"shutterangle"`          // 셔터앵글. 예) 172.8
	Filter                string `json:"filter"`                // 렌즈 필터기록
	StereoConvergence     string `json:"stereoconvergence"`     // 입체 좌우 카메라각도.
	StereoIA              string `json:"stereoia"`              // 이 값은 Setellite에서 사용하는듯함.
	Updatetime            string `json:"updatetime"`            // 아이템 업데이트 시간
}

// setSetellite함수는 Setellite자료구조의 필드에 값을 설정하는 함수이다.
func setSetellite(s *Setellite, field string, value string) {
	v := reflect.ValueOf(s).Elem().FieldByName(field)
	if v.IsValid() {
		v.SetString(value)
	}
}

// setId 메소드는 Setellite자료구조에 ID를 설정한다.
// ID = ShootDay + SlateNumber + TakeNumber + RollMedia
func (s *Setellite) setID() {
	id := s.ShootDay + s.SlateNumber + s.TakeNumber + s.RollMedia
	// mongodb의 필드명은 "$"문자로 시작되면 안된다.
	if strings.HasPrefix(s.ShootDay, "$") {
		id = s.SlateNumber + s.TakeNumber + s.RollMedia
	}
	// mongodb의 필드명에 "."문자가 포함되면 안된다.
	s.ID = strings.Replace(id, ".", "", -1)
}
