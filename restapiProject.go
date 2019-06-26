package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

// handleAPIAddproject 함수는 프로젝트를 추가한다.
func handleAPIAddproject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	id := r.FormValue("id")
	p := *NewProject(id)
	err = addProject(session, p)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	project, err := getProject(session, id)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	err = json.NewEncoder(w).Encode(project)
	if err != nil {
		log.Println(err)
	}
}

// handleAPIProject 함수는 프로젝트 정보를 불러온다.
func handleAPIProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	q := r.URL.Query()
	id := q.Get("id")
	project, err := getProject(session, id)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	err = json.NewEncoder(w).Encode(project)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIProjectTags 함수는 프로젝트에 사용되는 태그리스트를 불러온다.
func handleAPIProjectTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	type recipe struct {
		Data  []string `json:"data"`
		Error string   `json:"error"`
	}
	rcp := recipe{}
	_, err = getProject(session, project)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	rcp.Data, err = Distinct(session, project, "tag")
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
	}
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIProjects 함수는 프로젝트 리스트를 반환한다.
func handleAPIProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	q := r.URL.Query()
	qStatus := q.Get("status")
	type recipe struct {
		Data  []string `json:"data"`
		Error string   `json:"error"`
	}
	rcp := recipe{}
	projectList, err := getProjects(session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	for _, p := range projectList {
		switch qStatus {
		case "unknown":
			if p.Status == UnknownProjectStatus {
				rcp.Data = append(rcp.Data, p.ID)
			}
		case "pre":
			if p.Status == PreProjectStatus {
				rcp.Data = append(rcp.Data, p.ID)
			}
		case "post":
			if p.Status == PostProjectStatus {
				rcp.Data = append(rcp.Data, p.ID)
			}
		case "layover":
			if p.Status == LayoverProjectStatus {
				rcp.Data = append(rcp.Data, p.ID)
			}
		case "backup":
			if p.Status == BackupProjectStatus {
				rcp.Data = append(rcp.Data, p.ID)
			}
		case "archive":
			if p.Status == ArchiveProjectStatus {
				rcp.Data = append(rcp.Data, p.ID)
			}
		case "lawsuit":
			if p.Status == LawsuitProjectStatus {
				rcp.Data = append(rcp.Data, p.ID)
			}
		default:
			// qStatus값이 빈 문자열이면 작업중인 프로젝트를 results 리스트에 추가한다.
			// 작업중인 상태는 pre(프리프로덕션), post(포스트프로덕션), backup(백업중)인 상태를 뜻한다.
			if p.Status == PreProjectStatus || p.Status == PostProjectStatus || p.Status == BackupProjectStatus {
				rcp.Data = append(rcp.Data, p.ID)
			}
		}
	}
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}
