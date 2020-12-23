package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
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
	rcp.Review.AuthorNameKor = r.FormValue("authornamekor")
	if rcp.Review.AuthorNameKor == "" {
		// authornamekor 값이 비어있다면, 사용자의 아이디를 이용해서 DB에 등록된 이름을 가지고 온다.
		user, err := getUser(session, author)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Review.AuthorNameKor = user.LastNameKor + user.FirstNameKor
	}
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "path를 설정해주세요", http.StatusBadRequest)
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, path+"파일이 서버에 존재하지 않습니다", http.StatusBadRequest)
		return
	}
	rcp.Review.Path = path
	fpsString := r.FormValue("fps")
	fps, err := strconv.ParseFloat(fpsString, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s 는 fps로 사용할 수 없는 값 입니다", fpsString), http.StatusBadRequest)
		return
	}
	rcp.Review.Fps = fps
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
		if !(0 < n && n < 101) {
			http.Error(w, "progress의 값은 1~100 사이의 수가 되어야 합니다.", http.StatusBadRequest)
			return
		}
		rcp.Review.Progress = n
	}
	rcp.Review.Createtime = time.Now().Format("2006-01-02 15:04:05")
	rcp.Review.Updatetime = rcp.Review.Createtime
	mainVer, err := strconv.Atoi(r.FormValue("mainversion"))
	if err != nil {
		http.Error(w, "mainversion은 숫자로 입력되어 합니다", http.StatusBadRequest)
		return
	}
	rcp.MainVersion = mainVer
	subVer, err := strconv.Atoi(r.FormValue("subversion"))
	if err != nil {
		rcp.SubVersion = 0 // 서브버전은 없을 수 있다. 설정되지 않는다면 0값을 기본으로 한다.
	} else {
		rcp.SubVersion = subVer
	}
	rcp.Review.ID = bson.NewObjectId()
	rcp.Review.ProcessStatus = "wait" // ffmpeg 연산을 기다리는 상태로 등록한다.
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

// handleAPISearchReview 함수는 review를 검색하는 핸들러이다.
func handleAPISearchReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		UserID  string   `json:"userid"`
		Reviews []Review `json:"review"`
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
	searchword := r.FormValue("searchword")
	if searchword == "" {
		http.Error(w, "searchword를 설정해주세요", http.StatusBadRequest)
		return
	}
	reviews, err := searchReview(session, searchword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Reviews = reviews
	data, err := json.Marshal(rcp.Reviews)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIReview 함수는 id를 받아서 review 데이터를 반환하는 핸들러이다.
func handleAPIReview(w http.ResponseWriter, r *http.Request) {
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
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
		return
	}
	review, err := getReview(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetReviewStatus 함수는 review의 상태를 설정하는 RestAPI 이다.
func handleAPISetReviewStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		UserID string `json:"userid"`
		ID     string `json:"id"`
		Status string `json:"status"`
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
	reviewID := r.FormValue("id")
	if reviewID == "" {
		http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = reviewID

	status := r.FormValue("status")
	if status == "" {
		http.Error(w, "status를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Status = status
	review, err := getReview(session, rcp.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = setReviewStatus(session, rcp.ID, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log
	err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Set Review Status: %s, %s", rcp.ID, rcp.Status), review.Project, review.Name, "csi3", rcp.UserID, 180)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// slack log
	err = slacklog(session, review.Project, fmt.Sprintf("Set Review Status: %s, \nProject: %s, Name: %s, Author: %s", status, review.Project, review.Name, rcp.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// handleAPIAddReviewComment 함수는 review에 comment를 설정하는 RestAPI 이다.
func handleAPIAddReviewComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		UserID     string `json:"userid"`
		ID         string `json:"id"`
		Text       string `json:"text"`
		Media      string `json:"media"`
		MediaTitle string `json:"mediatitle"`
		Author     string `json:"author"`
		Date       string `json:"date"`
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "comment를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Text = text
	rcp.Media = r.FormValue("media")
	rcp.MediaTitle = r.FormValue("mediatitle")

	review, err := getReview(session, rcp.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cmt := Comment{}
	cmt.Date = time.Now().Format(time.RFC3339)
	rcp.Date = cmt.Date
	cmt.Author = rcp.UserID
	rcp.Author = rcp.UserID
	cmt.Text = rcp.Text
	cmt.Media = rcp.Media
	cmt.MediaTitle = rcp.MediaTitle
	err = addReviewComment(session, rcp.ID, cmt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log
	err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Add Review Comment: %s, %s", rcp.ID, rcp.Text), review.Project, review.Name, "csi3", rcp.UserID, 180)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// slack log
	err = slacklog(session, review.Project, fmt.Sprintf("Add Review Comment: %s, \nProject: %s, Name: %s, Author: %s", rcp.Text, review.Project, review.Name, rcp.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// handleAPIEditReviewComment 함수는 리뷰에서 코멘트를 수정합니다.
func handleAPIEditReviewComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		ID     string `json:"id"`
		Time   string `json:"time"`
		Text   string `json:"text"`
		UserID string `json:"userid"`
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

	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	reviewID := r.FormValue("id")
	if reviewID == "" {
		http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = reviewID
	reviewTime := r.FormValue("time")
	if reviewTime == "" {
		http.Error(w, "time을 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Time = reviewTime
	reviewText := r.FormValue("text")
	if reviewText == "" {
		http.Error(w, "text를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Text = reviewText
	err = EditReviewComment(session, rcp.ID, rcp.Time, rcp.Text)
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

// handleAPIRmReviewComment 함수는 리뷰에서 수정사항을 삭제합니다.
func handleAPIRmReviewComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		ID      string `json:"id"`
		Time    string `json:"time"`
		Project string `json:"project"`
		Name    string `json:"name"`
		Text    string `json:"text"`
		UserID  string `json:"userid"`
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
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	reviewID := r.FormValue("id")
	if reviewID == "" {
		http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = reviewID
	reviewTime := r.FormValue("time")
	if reviewTime == "" {
		http.Error(w, "time을 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Time = reviewTime

	// ID를 이용해서 삭제할 리뷰아이템을 가져와 Project, Name, Text를 반환될 json에 설정합니다.
	review, err := getReview(session, rcp.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Project = review.Project
	rcp.Name = review.Name
	for _, t := range review.Comments {
		if t.Date == reviewTime {
			rcp.Text = t.Text
			break
		}
	}

	// 리뷰데이터를 삭제합니다.
	err = RmReviewComment(session, rcp.ID, rcp.Time)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// log
	err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Rm Review Comment: %s, %s", rcp.ID, rcp.Text), rcp.Project, rcp.Name, "csi3", rcp.UserID, 180)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Rm Review Comment: %s\nProject: %s, Name: %s, Author: %s", rcp.Text, rcp.Project, rcp.Name, rcp.UserID))
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

// handleAPIRmReview 함수는 리뷰를 삭제합니다.
func handleAPIRmReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		ID     string `json:"id"`
		UserID string `json:"userid"`
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
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	reviewID := r.FormValue("id")
	if reviewID == "" {
		http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = reviewID
	adminSetting, err := GetAdminSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 데이터 삭제
	reviewfile := fmt.Sprintf("%s/%s.mp4", adminSetting.ReviewDataPath, rcp.ID)
	if _, err := os.Stat(reviewfile); err == nil {
		err = os.Remove(reviewfile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// DB삭제
	err = RmReview(session, rcp.ID)
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

// handleAPISetReviewProject 함수는 리뷰에서 Project를 설정합니다.
func handleAPISetReviewProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		ID      string `json:"id"`
		Project string `json:"project"`
		UserID  string `json:"userid"`
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

	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	reviewID := r.FormValue("id")
	if reviewID == "" {
		http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = reviewID
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project 를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Project = project
	err = SetReviewProject(session, rcp.ID, rcp.Project)
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

// handleAPISetReviewTask 함수는 리뷰에서 Task를 설정합니다.
func handleAPISetReviewTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		ID     string `json:"id"`
		Task   string `json:"task"`
		UserID string `json:"userid"`
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

	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	reviewID := r.FormValue("id")
	if reviewID == "" {
		http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = reviewID
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "task 를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	err = SetReviewTask(session, rcp.ID, rcp.Task)
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
