package main

// Review 는 리뷰데이터 자료구조 이다.
type Review struct {
	Project     string    `json:"project"`     // 프로젝트
	Task        string    `json:"task"`        // 태스크
	Createtime  string    `json:"createtime"`  // 생성시간
	Author      string    `json:"author"`      // 작성자
	Path        string    `json:"path"`        // 리뷰경로
	PathMp4     string    `json:"pathmp4"`     // Mp4 패스경로
	PathOgg     string    `json:"pathogg"`     // Ogg 패스경로
	Status      string    `json:"status"`      // 상태 approve, comment, waiting
	Sketches    []Sketch  `json:"sketches"`    // 스케치 프레임
	Playlist    []string  `json:"playlist"`    // 플레이리스트 목록
	Comments    []Comment `json:"comments"`    // 댓글
	Description string    `json:"description"` // 설명
}

// Sketch 는 스케치 자료구조이다.
type Sketch struct {
	Frame      int    `json:"frame"`      // 프레임수
	Author     string `json:"author"`     // 스케치를 그린사람
	SketchPath string `json:"sketchpath"` // 스케치 경로
}
