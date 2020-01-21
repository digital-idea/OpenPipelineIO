package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"gopkg.in/mgo.v2"
)

// handleAPITasksetting 함수는 Tasksetting 리소스의 restAPI 이다.
func handleAPITasksetting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project   string `json:"project"`
		Name      string `json:"name"`
		Task      string `json:"task"`
		Type      string `json:"type"`
		Assettype string `json:"assettype"`
		OS        string `json:"os"`
		Seq       string `json:"seq"`
		Cut       string `json:"cut"`
		Path      string `json:"path"`
		UserID    string `json:"userid"`
		Error     string `json:"error"`
	}
	rcp := Recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.UserID, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	for key, values := range r.PostForm {
		switch key {
		case "userid":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if rcp.UserID == "unknown" && v != "" {
				rcp.UserID = v
			}
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Project = v
		case "name":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Name = v
		case "task":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Task = v
		case "seq":
			if len(values) == 1 {
				rcp.Seq = values[0]
			}
		case "cut":
			if len(values) == 1 {
				rcp.Cut = values[0]
			}
		case "os":
			if len(values) == 1 {
				rcp.OS = values[0]
			}
		case "assettype":
			if len(values) == 1 {
				rcp.Assettype = values[0]
			}
		case "type":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Type = v
		}
	}
	typ := "shot"
	if rcp.Type == "asset" {
		typ = "asset"
	}
	t, err := getTaskSetting(session, rcp.Task+typ)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch rcp.OS {
	case "macos":
		rcp.Path = t.MacOSPath
	case "linux":
		rcp.Path = t.LinuxPath
	case "windows":
		rcp.Path = t.WindowPath
	default:
		rcp.Path = t.WFSPath // 기본적으로 WFS 경로를 사용한다.
	}
	// OS 경로를 불러와서 경로를 대입한다.
	tmpl, err := template.New("test").Parse(rcp.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Path = tpl.String()
	fmt.Println(tpl.String())
	// json 으로 결과 전송
	data, _ := json.Marshal(rcp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIShotTasksetting 함수는 Shot Task 항목을 반환하는 restAPI 이다.
func handleAPIShotTasksetting(w http.ResponseWriter, r *http.Request) {
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
	type Recipe struct {
		Tasksettings []Tasksetting
	}
	rcp := Recipe{}
	rcp.Tasksettings, err = getShotTaskSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, _ := json.Marshal(rcp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIAssetTasksetting 함수는 Shot Task 항목을 반환하는 restAPI 이다.
func handleAPIAssetTasksetting(w http.ResponseWriter, r *http.Request) {
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
	type Recipe struct {
		Tasksettings []Tasksetting
	}
	rcp := Recipe{}
	rcp.Tasksettings, err = getAssetTaskSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, _ := json.Marshal(rcp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPICategoryTasksettings 함수는 Category 문자를 입력받아 해당 Category Task항목을 반환하는 restAPI 이다.
func handleAPICategoryTasksettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Category     string        `json:"category"`
		Tasksettings []Tasksetting `json:"tasksettings"`
	}
	rcp := Recipe{}
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
	for key, values := range r.PostForm {
		switch key {
		case "category":
			if len(values) != 1 {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Category = values[0]
		}
	}
	rcp.Tasksettings, err = getCategoryTaskSettings(session, rcp.Category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
