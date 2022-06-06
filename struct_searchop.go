package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2"
)

// SearchOption 은 웹 검색창의 옵션 자료구조이다.
type SearchOption struct {
	Project           string `json:"project"`           // 선택한 프로젝트
	Searchword        string `json:"searchword"`        // 검색어
	Sortkey           string `json:"sortkey"`           // 정렬방식
	SearchbarTemplate string `json:"searchbartemplate"` // 검색바 탬플릿 이름
	Task              string `json:"task"`              // Task명
	Assign            bool   `json:"assign"`            // legacy
	Ready             bool   `json:"ready"`             // legacy
	Wip               bool   `json:"wip"`               // legacy
	Confirm           bool   `json:"confirm"`           // legacy
	Done              bool   `json:"done"`              // legacy
	Omit              bool   `json:"omit"`              // legacy
	Hold              bool   `json:"hold"`              // legacy
	Out               bool   `json:"out"`               // legacy
	None              bool   `json:"none"`              // legacy
	// 상태
	TrueStatus []string `json:"truestatus"` // true 상태리스트
	// 요소
	Shot   bool `json:"shot"`
	Assets bool `json:"assets"`
	Type3d bool `json:"type3d"`
	Type2d bool `json:"type2d"`
	// Page
	Page int `json:"page"`
}

// SearchOption과 관련된 메소드

func (op *SearchOption) setStatusAll() error {
	op.Assign = true  // legacy
	op.Ready = true   // legacy
	op.Wip = true     // legacy
	op.Confirm = true // legacy
	op.Done = true    // legacy
	op.Omit = true    // legacy
	op.Hold = true    // legacy
	op.Out = true     // legacy
	op.None = true    // legacy
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		return err
	}
	defer session.Close()
	status, err := AllStatus(session)
	if err != nil {
		return err
	}
	for _, s := range status {
		op.TrueStatus = append(op.TrueStatus, s.ID)
	}
	return nil
}

func (op *SearchOption) setStatusDefaultV1() error {
	op.Assign = true  // legacy
	op.Ready = true   // legacy
	op.Wip = true     // legacy
	op.Confirm = true // legacy
	op.Done = false   // legacy
	op.Omit = false   // legacy
	op.Hold = false   // legacy
	op.Out = false    // legacy
	op.None = false   // legacy
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		return err
	}
	defer session.Close()
	status, err := AllStatus(session)
	if err != nil {
		return err
	}
	for _, s := range status {
		if !s.DefaultOn {
			continue
		}
		op.TrueStatus = append(op.TrueStatus, s.ID)
	}
	return nil
}

func (op *SearchOption) setStatusDefaultV2() error {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		return err
	}
	defer session.Close()
	status, err := AllStatus(session)
	if err != nil {
		return err
	}
	for _, s := range status {
		if !s.DefaultOn {
			continue
		}
		op.TrueStatus = append(op.TrueStatus, s.ID)
	}
	return nil
}

func (op *SearchOption) setStatusNone() {
	op.Assign = false  // legacy
	op.Ready = false   // legacy
	op.Wip = false     // legacy
	op.Confirm = false // legacy
	op.Done = false    // legacy
	op.Omit = false    // legacy
	op.Hold = false    // legacy
	op.Out = false     // legacy
	op.None = false    // legacy
	op.TrueStatus = []string{}
}

// isStatusOff 메소드는 모든 상태가 꺼저 있는지 체크한다.
func (op *SearchOption) isStatusOff() bool {
	if op.Assign || op.Ready || op.Wip || op.Confirm || op.Done || op.Omit || op.Hold || op.Out || op.None {
		return false
	}
	return true
}

func handleRequestToSearchOption(r *http.Request) SearchOption {
	q := r.URL.Query()
	op := SearchOption{}
	op.Project = q.Get("project")
	op.Searchword = q.Get("searchword")
	op.Sortkey = q.Get("sortkey")
	op.SearchbarTemplate = q.Get("searchbartemplate")
	op.Task = q.Get("task")
	// 페이지를 구한다.
	page, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		op.Page = 1 // 에러가 발생하면 1페이지로 이동한다.
	} else {
		op.Page = page
	}
	op.Assign = str2bool(q.Get("assign"))   // legacy
	op.Ready = str2bool(q.Get("ready"))     // legacy
	op.Wip = str2bool(q.Get("wip"))         // legacy
	op.Confirm = str2bool(q.Get("confirm")) // legacy
	op.Done = str2bool(q.Get("done"))       // legacy
	op.Omit = str2bool(q.Get("omit"))       // legacy
	op.Hold = str2bool(q.Get("hold"))       // legacy
	op.Out = str2bool(q.Get("out"))         // legacy
	op.None = str2bool(q.Get("none"))       // legacy
	for _, s := range strings.Split(q.Get("truestatus"), ",") {
		if s == "" {
			continue
		}
		op.TrueStatus = append(op.TrueStatus, s)
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
		if cookie.Name == "SearchbarTemplate" {
			op.SearchbarTemplate = cookie.Value
		}
	}
	if op.Project == "" {
		plist, err := Projectlist(session)
		if err != nil {
			return err
		}
		op.Project = plist[0] // 프로젝트가 빈 문자열이면 첫번째 프로젝트를 설정합니다.
	}
	return nil
}
