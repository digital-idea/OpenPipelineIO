package main

import (
	"encoding/base64"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

// handleInputMode 함수는 수정을 편하게 입력하는 페이지 이다.
func handleInputMode(w http.ResponseWriter, r *http.Request) {
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
	type recipe struct {
		User
		SessionID   string
		Devmode     bool
		Projectlist []string
		Items       []Item
		Project     string
		Ddline3d    []string
		Ddline2d    []string
		Tags        []string
		Assettags   []string
		SearchOption
		Searchnum Infobarnum
		Totalnum  Infobarnum
	}
	rcp := recipe{}
	rcp.SearchOption.LoadCookie(r)

	rcp.Devmode = *flagDevmode
	rcp.SessionID = ssid.ID
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Projectlist, err = OnProjectlist(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(rcp.Projectlist) == 0 {
		http.Redirect(w, r, "/noonproject", http.StatusSeeOther)
		return
	}
	rcp.Ddline3d, err = DistinctDdline(session, rcp.SearchOption.Project, "ddline3d")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Ddline2d, err = DistinctDdline(session, rcp.SearchOption.Project, "ddline2d")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Tags, err = Distinct(session, rcp.SearchOption.Project, "tag")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Assettags, err = Distinct(session, rcp.SearchOption.Project, "assettags")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Totalnum, err = Totalnum(session, rcp.SearchOption.Project)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.SearchOption = handleRequestToSearchOption(r)
	q := r.URL.Query()
	template := q.Get("template")
	rcp.SearchOption.Template = template
	rcp.Items, err = Searchv2(session, rcp.SearchOption)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Searchnum, err = Searchnum(rcp.SearchOption.Project, rcp.Items)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 최종적으로 사용된 프로젝트명을 쿠키에 저장한다.
	cookie := http.Cookie{
		Name:   "Project",
		Value:  rcp.SearchOption.Project,
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Task",
		Value:  rcp.SearchOption.Task,
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Searchword",
		Value:  base64.StdEncoding.EncodeToString([]byte(rcp.SearchOption.Searchword)), //  쿠키는 UTF-8을 저장할 때 에러가 발생한다.
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)

	err = TEMPLATES.ExecuteTemplate(w, rcp.SearchOption.Template, rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
