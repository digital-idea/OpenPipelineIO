package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
)

// handleAPISetelliteItems 함수는 project, rollmedia을 받아서 setellite 정보를 반환한다.
func handleAPISetelliteItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	rollmedia := q.Get("rollmedia")
	type recipe struct {
		Data []Setellite `json:"data"`
	}
	rcp := recipe{}
	rcp.Data, err = SetelliteItems(session, project, rollmedia)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"err\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetelliteSearch 함수는 project, searchword의 검색어를 통해 setellite 정보를 반환한다.
func handleAPISetelliteSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	searchword := q.Get("searchword")
	type recipe struct {
		Data []Setellite `json:"data"`
	}
	rcp := recipe{}
	rcp.Data, err = SetelliteSearch(session, project, searchword)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"err\":\"%v\"}\n", err)
		return
	}
}
