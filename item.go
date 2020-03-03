package main

import (
	"sort"
)

const ( // legacy
	// CLIENT 클라이언트 컨펌상태
	CLIENT = "9"
	// OMIT 작업취소 상태
	OMIT = "8"
	// CONFIRM 내부 컨펌상태
	CONFIRM = "7"
	// WIP 작업중 상태
	WIP = "6"
	// READY 작업준비중 상태
	READY = "5"
	// ASSIGN 작업자 배정을 기다리는 상태
	ASSIGN = "4"
	// OUT 외주상태(삭제예정이다.)
	OUT = "3"
	// DONE 작업완료 상태
	DONE = "2"
	// HOLD 작업중단 상태
	HOLD = "1"
	// NONE 상태없음. 예) 소스
	NONE = "0"
)

// TaskLevel 은 태스크 난이도이다.
type TaskLevel int

// TaskLevel list
const (
	TaskLevel0 = TaskLevel(iota) // 0 쉬운난이도
	TaskLevel1                   // 1
	TaskLevel2                   // 2
	TaskLevel3                   // 3
	TaskLevel4                   // 4
	TaskLevel5                   // 5 높은난이도
)

// Comment 자료구조는 글을 작성할 때 사용하는 자료구조이다.
type Comment struct {
	Date   string `json:"date"`   // 코맨트 등록시간 RFC3339
	Author string `json:"author"` // 작성자
	Text   string `json:"text"`   // 내용
	Media  string `json:"media"`  // media 경로
}

// Source 자료구조는 글을 작성할 때 사용하는 자료구조이다.
type Source struct {
	Date   string // 등록시간 RFC3339
	Author string // 작성자
	Title  string // 제목
	Path   string // 소스 경로
}

// Version 자료구조는 버전정보를 담을 때 사용하는 자료구조이다.
type Version struct {
	Main int // 메인버전 v01 또는 v01_w2 형태에서 앞부분 "1" 이다.
	Sub  int // 서브버전 v02_w02 형태에서 뒷부분 "2" 이다.
}

// Item 자료구조는 하나의 항목에 대한 자료구조이다.
type Item struct {
	Project string `json:"project"` // 프로젝트명
	ID      string `json:"id"`      // ID

	// 현장정보
	// 현장에서 사용하는 카메라 데이터 이름. 슈퍼바이저 툴과 연동하기 위해서 Key로 사용된다.
	// 일반적으로 스캔이름과 같지만 항상 동일하지 않다.
	// 카메라 데이터 A037C012_160708_R717.[11694428-1172226].ari 형태에서 A037C012_160708_R717 부분이 데이터 이름이다.
	Dataname string `json:"dataname"` // 영화카메라(Red,Alexa등)이 자동 생성시키는 이미지 파일명이다.

	// 작업이 필요한 정보
	Scanname    string          `json:"scanname"`    // 스캔이름
	Platesize   string          `json:"platesize"`   // 플레이트 이미지사이즈
	Name        string          `json:"name"`        // 샷이름 SS_0010
	Seq         string          `json:"seq"`         // 시퀀스이름 SS_0010 에서 SS문자에 해당하는값. 에셋이면 "" 문자열이 들어간다.
	Cut         string          `json:"cut"`         // 시퀀스이름 SS_0010 에서 0010문자에 해당하는값. 에셋이면 "" 문자열이 들어간다.
	Type        string          `json:"type"`        // org, org1, src, asset..
	Assettype   string          `json:"assettype"`   // char, env, prop, comp, plant, vehicle, group
	CrowdAsset  bool            `json:"crowdasset"`  // 군중씬에서 사용하는 에셋인지 여부 체크
	UseType     string          `json:"usetype"`     // 재스캔상황시 실제로 사용해야하는 타입표기
	Scantime    string          `json:"scantime"`    // 스캔 등록시간 RFC3339
	Thumpath    string          `json:"thumpath"`    // 썸네일경로
	Thummov     string          `json:"thummov"`     // 썸네일 mov 경로
	Beforemov   string          `json:"beforemov"`   // 전에 들어갈 mov. 만약 2개 이상이라면 space로 구분한다.
	Aftermov    string          `json:"aftermov"`    // 후에 들어갈 mov. 만약 2개 이상이라면 space로 구분한다.
	Retimeplate string          `json:"retimeplate"` // 리타임 플레이트 경로
	Platepath   string          `json:"platepath"`   // 플레이트 경로
	Shottype    string          `json:"shottype"`    // "", "2d", "3d"
	Ddline3d    string          `json:"ddline3d"`    // 3D 데드라인 RFC3339
	Ddline2d    string          `json:"ddline2d"`    // 2D 데드라인 RFC3339
	Rnum        string          `json:"rnum"`        // 롤넘버, 영화를 권으로 나누었을 때 이 샷의 권 번호 예) A0001. A는 1권을 H는 8권을 의미한다.
	Tag         []string        `json:"tag"`         // 태그리스트
	Assettags   []string        `json:"assettags"`   // 에셋그룹 태그
	Finname     string          `json:"finname"`     // 파이널 파일이름
	Finver      string          `json:"finver"`      // 파이널된 버젼
	Findate     string          `json:"findate"`     // 파이널 데이터가 나간 날짜
	Clientver   string          `json:"clientver"`   // 클라이언트에게 보낸 버전
	Dsize       string          `json:"dsize"`       // 디스토션 사이즈
	Rendersize  string          `json:"rendersize"`  // 특수상황시 렌더사이즈. 예) 5k플레이트를 3D에서 2k영역만 잡아서 최종 아웃풋까지 이어질 때
	Status      string          `json:"status"`      // 샷 상태. legacy
	StatusV2    string          `json:"statusv2"`    // 샷 상태.
	Updatetime  string          `json:"updatetime"`  // 업데이트 시간 RFC3339
	Focal       string          `json:"focal"`       // 렌즈 미리수
	Stereotype  string          `json:"stereotype"`  // parallel(default), conversions
	Stereoeye   string          `json:"stereoeye"`   // left(default), right
	Outputname  string          `json:"outputname"`  // 프로젝트중 클라이언트가 제시하는 아웃풋 이름
	OCIOcc      string          `json:"ociocc"`      // Neutural Grading Pipeline에 사용하는 .cc 파일의 경로.
	Note        Comment         `json:"note"`        // 작업내용
	Sources     []Source        `json:"links"`       // 연결소스
	References  []Source        `json:"references"`  // 레퍼런스
	Comments    []Comment       `json:"comments"`    // 수정내용
	Tasks       map[string]Task `json:"tasks"`       // Task 리스트

	//시간에 관련된 데이터이다.
	ScanFrame       int                    `json:"scanframe"`       // 스캔 프레임수
	ScanTimecodeIn  string                 `json:"scantimecodein"`  // 스캔플레이트 타임코드 In
	ScanTimecodeOut string                 `json:"scantimecodeout"` // 스캔플레이트 타임코드 Out
	ScanIn          int                    `json:"scanin"`          // 스캔 Frame In
	ScanOut         int                    `json:"scanout"`         // 스캔 Frame Out
	HandleIn        int                    `json:"handlein"`        // 핸들 Frame In
	HandleOut       int                    `json:"handleout"`       // 핸들 Frame Out
	JustIn          int                    `json:"justin"`          // 저스트 Frame In
	JustOut         int                    `json:"justout"`         // 저스트 Frame Out
	JustTimecodeIn  string                 `json:"justtimecodein"`  // 저스트 타임코드 In
	JustTimecodeOut string                 `json:"justtimecodeout"` // 저스트 타임코드 Out
	PlateIn         int                    `json:"platein"`         // 플레이트 Frame In
	PlateOut        int                    `json:"plateout"`        // 플레이트 Frame Out
	Soundfile       string                 `json:"soundfile"`       // 사운드파일 필요시 사운드파일 경로
	Rollmedia       string                 `json:"rollmedia"`       // 현장데이터의 Rollmedia 문자. 수동으로 현장데이터와 연결할 때 사용한다.
	ObjectidIn      int                    `json:"objectidin"`      // ObjectID 시작번호. Deep이미지의 DeepID를 만들기 위해서 파이프라인상 필요하다.
	ObjectidOut     int                    `json:"objectidout"`     // ObjectID 끝번호. Deep이미지의 DeepID를 만들기 위해서 파인라인상 필요하다.
	OnsetCam        `json:"onsetcam"`      // 현장 카메라 정보
	ProductionCam   `json:"productioncam"` // 포스트 프로덕션 카메라 정보
}

// Task 자료구조는 태크스 정보를 담는 자료구조이다.
type Task struct {
	Title        string             `json:"title"`        // 테스크 네임
	User         string             `json:"user"`         // 아티스트명
	Status       string             `json:"status"`       // 상태 legacy
	StatusV2     string             `json:"statusv2"`     // 샷 상태.
	BeforeStatus string             `json:"beforestatus"` // 이전상태
	Startdate    string             `json:"startdate"`    // 작업시작일 RFC3339
	Predate      string             `json:"predate"`      // 1차 마감일 RFC3339
	Date         string             `json:"date"`         // 2차 마감일 RFC3339
	Mov          string             `json:"mov"`          // mov 경로
	Mdate        string             `json:"mdate"`        // mov 업데이트 날짜 RFC3339
	Pubfile      string             `json:"pubfile"`      // Pubfile
	Due          int                `json:"due"`          // 예측 멘데이
	Promday      int                `json:"promday"`      // 실제 멘데이
	UserNote     string             `json:"usernote"`     // 아티스트와 관련된 엘리먼트등의 정보를 입력하기 위해 사용.
	TaskLevel    `json:"tasklevel"` // 샷 레벨
	Version      `json:"version"`   // Pubfile 버전정보
}

// updateStatus는 각 팀의 상태를 조합해서 샷 상태를 업데이트하는 함수이다.
func (item *Item) updateStatus() {
	maxstatus := "0"
	for _, value := range item.Tasks {
		if value.Status > maxstatus {
			maxstatus = value.Status
		}
	}
	item.Status = maxstatus
}

// setRumTag 는 특정 항목이 입력이 되었을때 알맞은 태그를 자동으로 넣거나 삭제할 때 사용한다.
// 예를 들어 "A0001" 이라는 롤넘버가 셋팅되면 태그리스트에 "1권" 이라는 단어를 넣어준다.
func (item *Item) setRnumTag() {
	if item.Rnum == "" {
		return
	}
	var rnumTag string
	switch item.Rnum[0] {
	case 'a', 'A':
		rnumTag = "1권"
	case 'b', 'B':
		rnumTag = "2권"
	case 'c', 'C':
		rnumTag = "3권"
	case 'd', 'D':
		rnumTag = "4권"
	case 'e', 'E':
		rnumTag = "5권"
	case 'f', 'F':
		rnumTag = "6권"
	case 'g', 'G':
		rnumTag = "7권"
	case 'h', 'H':
		rnumTag = "8권"
	}
	var newTags []string
	for _, t := range item.Tag {
		if !(t == "1권" || t == "2권" || t == "3권" || t == "4권" || t == "5권" || t == "6권" || t == "7권" || t == "8권") {
			newTags = append(newTags, t)
		}
	}
	newTags = append(newTags, rnumTag)
	sort.Strings(newTags)
	item.Tag = newTags
}

// setAssettags는 Assettype 이 변경시 Assettags 리스트를 변경한다.
func (item *Item) setAssettags() {
	var tags []string
	if item.Assettype != "" {
		tags = append(tags, item.Assettype)
	}
	for _, t := range item.Assettags {
		b, _ := validAssettype(t)
		if b {
			continue
		}
		tags = append(tags, t)
	}
	sort.Strings(tags)
	item.Assettags = tags
}
