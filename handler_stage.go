package main

import (
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2"
)

// handleAddStage 함수는 Stage를 추가하는 페이지이다.
func handleAddStage(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 4 {
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
	type recipe struct {
		User    User
		Devmode bool
		SearchOption
		Stages []Stage
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
	rcp.Stages, err = AllStages(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "addstage", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddStageSubmit 함수는 Stage를 추가합니다.
func handleAddStageSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 4 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	order, err := strconv.ParseFloat(r.FormValue("order"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := Stage{
		ID:          r.FormValue("id"),
		Description: r.FormValue("description"),
		TextColor:   "#" + r.FormValue("textcolor"),
		BGColor:     "#" + r.FormValue("bgcolor"),
		BorderColor: "#" + r.FormValue("bordercolor"),
		Order:       order,
		InitStage:   str2bool(r.FormValue("initstage")),
		NextStage:   r.FormValue("nextstage"),
		NextStatus:  r.FormValue("nextstatus"),
	}
	err = s.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = AddStage(session, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/stage", http.StatusSeeOther)
}

// handleStage 함수는 Stage 를 보는 페이지이다.
func handleStage(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 4 {
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
	type recipe struct {
		User    User
		Devmode bool
		Stages  []Stage
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.Devmode = *flagDevmode
	rcp.Stages, err = AllStages(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	err = TEMPLATES.ExecuteTemplate(w, "stage", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleRmStage 함수는 Stage를 삭제하는 페이지이다.
func handleRmStage(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 5 {
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
	type recipe struct {
		User    User
		Devmode bool
		SearchOption
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
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
	err = TEMPLATES.ExecuteTemplate(w, "rmstage", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleRmStageSubmit 함수는 Stage를 삭제합니다.
func handleRmStageSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 5 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	id := r.FormValue("id")
	err = RmStage(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/stage", http.StatusSeeOther)
}

// handleEditStage 함수는 Stage를 편집하는 페이지이다.
func handleEditStage(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 4 {
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
	type recipe struct {
		User    User
		Devmode bool
		SearchOption
		Stage
		Stages []Stage
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
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
	q := r.URL.Query()
	id := q.Get("id")
	rcp.Stage, err = GetStage(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Stages, err = AllStages(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "editstage", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleEditStageSubmit 함수는 Stage를 수정합니다.
func handleEditStageSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 4 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	order, err := strconv.ParseFloat(r.FormValue("order"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := Stage{
		ID:          r.FormValue("id"),
		Description: r.FormValue("description"),
		TextColor:   "#" + r.FormValue("textcolor"),
		BGColor:     "#" + r.FormValue("bgcolor"),
		BorderColor: "#" + r.FormValue("bordercolor"),
		Order:       order,
		InitStage:   str2bool(r.FormValue("initstage")),
		NextStage:   r.FormValue("nextstage"),
		NextStatus:  r.FormValue("nextstatus"),
	}
	err = s.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = SetStage(session, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/stage", http.StatusSeeOther)
}
