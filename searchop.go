package main

import (
	"net/http"
)

// SearchOption 은 웹 검색창의 옵션 자료구조이다.
type SearchOption struct {
	Project    string // 선택한 프로젝트
	Searchword string // 검색어
	Sortkey    string // 정렬방식
	Template   string // 템플릿 이름
	// 상태
	Assign  bool
	Ready   bool
	Wip     bool
	Confirm bool
	Done    bool
	Omit    bool
	Hold    bool
	Out     bool
	None    bool
	// 요소
	Shot   bool
	Assets bool
	Type3d bool
	Type2d bool
}

// SearchOption과 관련된 메소드

func (op *SearchOption) setStatusAll() {
	op.Assign = true
	op.Ready = true
	op.Wip = true
	op.Confirm = true
	op.Done = true
	op.Omit = true
	op.Hold = true
	op.Out = true
	op.None = true
}

func (op *SearchOption) setStatusDefault() {
	op.Assign = true
	op.Ready = true
	op.Wip = true
	op.Confirm = true
	op.Done = false
	op.Omit = false
	op.Hold = false
	op.Out = false
	op.None = false
}

func (op *SearchOption) setStatusNone() {
	op.Assign = false
	op.Ready = false
	op.Wip = false
	op.Confirm = false
	op.Done = false
	op.Omit = false
	op.Hold = false
	op.Out = false
	op.None = false
}

func handleRequestToSearchOption(r *http.Request) SearchOption {
	q := r.URL.Query()
	op := SearchOption{
		Project:    q.Get("project"),
		Searchword: q.Get("searchword"),
		Sortkey:    q.Get("sortkey"),
		Template:   q.Get("template"),
		Assign:     str2bool(q.Get("assign")),
		Ready:      str2bool(q.Get("ready")),
		Wip:        str2bool(q.Get("wip")),
		Confirm:    str2bool(q.Get("confirm")),
		Done:       str2bool(q.Get("done")),
		Omit:       str2bool(q.Get("omit")),
		Hold:       str2bool(q.Get("hold")),
		Out:        str2bool(q.Get("out")),
		None:       str2bool(q.Get("none")),
	}
	return op
}
