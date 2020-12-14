package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2"
)

// handleAPIAddproject 함수는 프로젝트를 추가한다.
func handleAPIAddproject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	id := r.FormValue("id")
	p := *NewProject(id)
	err = addProject(session, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	project, err := getProject(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIProject 함수는 프로젝트 정보를 불러온다.
func handleAPIProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	q := r.URL.Query()
	id := q.Get("id")
	project, err := getProject(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIProjectTags 함수는 프로젝트에 사용되는 태그리스트를 불러온다.
func handleAPIProjectTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	_, err = getProject(session, project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tags, err := Distinct(session, project, "tag")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 태그중에 빈 문자열을 제거한다.
	var filteredTags []string
	for _, t := range tags {
		if t == "" {
			continue
		}
		filteredTags = append(filteredTags, t)
	}
	data, err := json.Marshal(filteredTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIProjectAssetTags 함수는 프로젝트에 사용되는 에셋 태그리스트를 불러온다.
func handleAPIProjectAssetTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	_, err = getProject(session, project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	assettags, err := Distinct(session, project, "assettags")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 에셋태그 중 빈 문자열을 제거한다.
	var filteredTags []string
	for _, t := range assettags {
		if t == "" {
			continue
		}
		filteredTags = append(filteredTags, t)
	}
	data, err := json.Marshal(filteredTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIProjects 함수는 프로젝트 리스트를 반환한다.
func handleAPIProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, p := range projectList {
		switch qStatus {
		case "test":
			if p.Status == TestProjectStatus {
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
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
