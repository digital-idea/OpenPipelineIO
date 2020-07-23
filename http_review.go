package main

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
)

func handleDaily(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	today := time.Now().Format("2006-01-02")
	// 오늘날짜를 구하고 리다이렉트한다.
	http.Redirect(w, r, "/review?searchword=daily:"+today, http.StatusSeeOther)
}

// handleReview 함수는 리뷰 페이지이다.
func handleReview(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	q := r.URL.Query()
	type recipe struct {
		User        User
		Projectlist []string
		Devmode     bool
		SearchOption
		Searchword       string
		Status           []Status // css 생성을 위해서 필요함
		CurrentReview    Review   // 현재 리뷰 자료구조
		Reviews          []Review // 옆 Review 항목
		ReviewGroup      []Review // 하단 Review 항목
		TasksettingNames []string
	}
	rcp := recipe{}
	rcp.Searchword = q.Get("searchword")
	id := q.Get("id")
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Devmode = *flagDevmode
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.TasksettingNames, err = TasksettingNames(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Status, err = AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Reviews, err = searchReview(session, rcp.Searchword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if id != "" {
		review, err := getReview(session, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.CurrentReview = review
		rcp.ReviewGroup, err = searchReview(session, fmt.Sprintf("project:%s name:%s", review.Project, review.Name))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = TEMPLATES.ExecuteTemplate(w, "review", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleReviewData 함수는 리뷰 영상데이터를 전송한다.
func handleReviewData(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	admin, err := GetAdminSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q := r.URL.Query()
	id := q.Get("id")
	http.ServeFile(w, r, fmt.Sprintf("%s/%s.mp4", admin.ReviewDataPath, id))
}

// handleReviewSubmit 함수는 리뷰 검색창의 검색어를 입력받아 새로운 URI로 리다이렉션 한다.
func handleReviewSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	searchword := r.FormValue("SearchReview")
	redirectURL := fmt.Sprintf(`/review?searchword=%s`,
		searchword,
	)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
