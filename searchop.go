package main

import (
	"encoding/base64"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

// SearchOption 은 웹 검색창의 옵션 자료구조이다.
type SearchOption struct {
	Project    string // 선택한 프로젝트
	Searchword string // 검색어
	Sortkey    string // 정렬방식
	Template   string // 템플릿 이름
	Task       string // Task명
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
		Task:       q.Get("task"),
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

// LoadCookie 메소드는 request에 이미 설정된 쿠키값을을 SearchOption 자료구조에 추가한다.
func (op *SearchOption) LoadCookie(session *mgo.Session, r *http.Request) error {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "Project" {
			op.Project = cookie.Value
		}
		if cookie.Name == "Task" {
			op.Task = cookie.Value
		}
		if cookie.Name == "Searchword" {
			cookieByte, err := base64.StdEncoding.DecodeString(cookie.Value)
			if err != nil {
				log.Println(err)
			}
			op.Searchword = string(cookieByte)
		}
		if cookie.Name == "Assign" {
			op.Assign = str2bool(cookie.Value)
		}
		if cookie.Name == "Ready" {
			op.Ready = str2bool(cookie.Value)
		}
		if cookie.Name == "Wip" {
			op.Wip = str2bool(cookie.Value)
		}
		if cookie.Name == "Confirm" {
			op.Confirm = str2bool(cookie.Value)
		}
		if cookie.Name == "Done" {
			op.Done = str2bool(cookie.Value)
		}
		if cookie.Name == "Omit" {
			op.Omit = str2bool(cookie.Value)
		}
		if cookie.Name == "Hold" {
			op.Hold = str2bool(cookie.Value)
		}
		if cookie.Name == "Out" {
			op.Out = str2bool(cookie.Value)
		}
		if cookie.Name == "None" {
			op.None = str2bool(cookie.Value)
		}
		if cookie.Name == "Template" {
			op.Template = cookie.Value
		}
	}
	if op.Project == "" {
		plist, err := Projectlist(session)
		if err != nil {
			return err
		}
		op.Project = plist[0]
	}
	if op.Template == "" {
		op.Template = "index2"
	}
	return nil
}
