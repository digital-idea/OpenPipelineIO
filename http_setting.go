package main

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

// handleAdminSetting 함수는 admin setting을 수정하는 페이지이다.
func handleAdminSetting(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel != 10 { // Admin
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		User        User
		Projectlist []string
		Devmode     bool
		SearchOption
		Setting
	}
	rcp := recipe{}
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
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Setting, err = GetAdminSetting(session)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "adminsetting", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAdminSettingSubmint 함수는 adminsetting을 업데이트 한다.
func handleAdminSettingSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel != 10 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	runScriptAfterSingup := r.FormValue("RunScriptAfterSignup")
	s := Setting{}
	s.ID = "admin"
	s.RunScriptAfterSignup = runScriptAfterSingup

	err = SetAdminSetting(session, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type recipe struct {
		User    User
		Devmode bool
		SearchOption
	}
	rcp := recipe{}
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
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "setadminsetting", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSetAdminSetting 함수는 admin setting을 수정하는 페이지이다.
func handleSetAdminSetting(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel != 10 { // Admin
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		User        User
		Projectlist []string
		Devmode     bool
		SearchOption
	}
	rcp := recipe{}
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
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "update_adminsetting", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
