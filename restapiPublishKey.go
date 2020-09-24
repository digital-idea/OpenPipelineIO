package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2"
)

// handleAPIPulishKey 함수는 모든 PublishKey 항목을 반환하는 restAPI 이다.
func handleAPIPublishKeys(w http.ResponseWriter, r *http.Request) {
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
	keys, err := AllPublishKeys(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, err := json.Marshal(keys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIGetPulish 함수는 Publish값을 반환하는 restAPI 이다.
func handleAPIGetPublish(w http.ResponseWriter, r *http.Request) {
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
	r.ParseForm()
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project 값이 빈 문자열입니다", http.StatusBadRequest)
		return
	}
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "id 값이 빈 문자열입니다", http.StatusBadRequest)
		return
	}
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "task 값이 빈 문자열입니다", http.StatusBadRequest)
		return
	}
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "key 값이 빈 문자열입니다", http.StatusBadRequest)
		return
	}
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "path 값이 빈 문자열입니다", http.StatusBadRequest)
		return
	}
	createtime := r.FormValue("createtime")
	if createtime == "" {
		http.Error(w, "createtime 값이 빈 문자열입니다", http.StatusBadRequest)
		return
	}
	item, err := getItem(session, project, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hasTask := false
	for _, t := range item.Tasks {
		if t.Title == task {
			hasTask = true
		}
	}
	if !hasTask {
		http.Error(w, task+" Task가 존재하지 않습니다", http.StatusBadRequest)
		return
	}
	hasKey := false
	hasTime := false
	hasPath := false
	pubInfo := Publish{}
	for k, pubList := range item.Tasks[task].Publishes {
		if key != k {
			continue
		}
		hasKey = true
		for _, p := range pubList {
			if p.Createtime == createtime && p.Path == path {
				hasTime = true
				hasPath = true
				pubInfo = p
			}
		}
	}
	if !hasKey {
		http.Error(w, key+" key로 Publish한 데이터가 존재하지 않습니다", http.StatusInternalServerError)
		return
	}
	if !hasPath {
		http.Error(w, path+" 값으로 Publish한 데이터가 존재하지 않습니다", http.StatusInternalServerError)
		return
	}
	if !hasTime {
		http.Error(w, createtime+" 시간으로 Publish한 데이터가 존재하지 않습니다", http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, err := json.Marshal(pubInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
