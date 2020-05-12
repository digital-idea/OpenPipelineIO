package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/digital-idea/dilog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// handleAPIAddReview 함수는 review를 추가하는 핸들러이다.
func handleAPIAddReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		UserID string `json:"userid"`
		Review
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

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Review.Project = project

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name을 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Review.Name = name

	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "task를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Review.Task = task

	author := r.FormValue("author")
	if author == "" {
		rcp.Review.Author = rcp.UserID
	}
	rcp.Review.Author = author

	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "path를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Review.Path = path
	rcp.Review.Status = "wait"

	rcp.Review.Description = r.FormValue("description")
	rcp.Review.CameraInfo = r.FormValue("camerainfo")
	progress := r.FormValue("progress")
	if progress != "" {
		n, err := strconv.Atoi(progress)
		if err != nil {
			http.Error(w, "progress의 값이 숫자가 아닙니다.", http.StatusBadRequest)
			return
		}
		if !(0 <= n && n <= 101) {
			http.Error(w, "progress의 값은 0~100 사이의 수가 되어야 합니다.", http.StatusBadRequest)
			return
		}
		rcp.Review.Progress = n
	}
	rcp.Review.Createtime = time.Now().Format("2006-01-02 15:04:05")
	rcp.Review.Updatetime = rcp.Review.Createtime
	rcp.ID = bson.NewObjectId()
	err = addReview(session, rcp.Review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log
	err = dilog.Add(*flagDBIP, host, fmt.Sprintf("AddReview: %s, %s, %s", rcp.Review.Name, rcp.Review.Task, rcp.Review.Path), rcp.Review.Project, rcp.Review.Name, "csi3", rcp.UserID, 180)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("AddReview: %s, %s\nProject: %s, Name: %s, Author: %s", rcp.Review.Task, rcp.Review.Path, rcp.Review.Project, rcp.Review.Name, rcp.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rcp.Review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
