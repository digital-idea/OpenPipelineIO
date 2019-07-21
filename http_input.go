package main

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

// handleInputTags 함수는 태그를 편하게 입력하는 페이지 이다.
func handleInputTags(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		SearchOption
	}
	rcp := recipe{}
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
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q := r.URL.Query()
	op := SearchOption{
		Project:      q.Get("project"),
		Searchword:   q.Get("searchword"),
		Sortkey:      q.Get("sortkey"),
		Assign:       str2bool(q.Get("assign")),
		Ready:        str2bool(q.Get("ready")),
		Wip:          str2bool(q.Get("wip")),
		Confirm:      str2bool(q.Get("confirm")),
		Done:         str2bool(q.Get("done")),
		Omit:         str2bool(q.Get("omit")),
		Hold:         str2bool(q.Get("hold")),
		Out:          str2bool(q.Get("out")),
		None:         str2bool(q.Get("none")),
		Template:     q.Get("template"),
		PostEndpoint: "/searchv2",
	}
	rcp.SearchOption = op
	rcp.Items, err = Searchv1(session, op)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, op.Template, rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
