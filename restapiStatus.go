package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2"
)

// handleAPIStatus 함수는 Status 모든 항목을 반환하는 restAPI 이다.
func handleAPIStatus(w http.ResponseWriter, r *http.Request) {
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
	status, err := AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIAddStatus 함수는 Status를 추가하는 RestAPI 이다.
func handleAPIAddStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	rcp := Status{}
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "itemtype을 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Description = r.FormValue("description")
	orderString := r.FormValue("order")
	if orderString == "" {
		http.Error(w, "order를 설정해주세요", http.StatusBadRequest)
		return
	}
	order, err := strconv.ParseFloat(orderString, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.Order = order
	textcolor := r.FormValue("textcolor")
	if !regexWebColor.MatchString(textcolor) {
		http.Error(w, "textcolor 값이 #FFFFFF 형태가 아닙니다", http.StatusBadRequest)
		return
	}
	rcp.TextColor = textcolor
	bgcolor := r.FormValue("bgcolor")
	if !regexWebColor.MatchString(bgcolor) {
		http.Error(w, "bgcolor 값이 #FFFFFF 형태가 아닙니다", http.StatusBadRequest)
		return
	}
	rcp.BGColor = bgcolor
	bordercolor := r.FormValue("bordercolor")
	if !regexWebColor.MatchString(bordercolor) {
		http.Error(w, "bordercolor 값이 #FFFFFF 형태가 아닙니다", http.StatusBadRequest)
		return
	}
	rcp.BorderColor = bordercolor
	rcp.DefaultOn = str2bool(r.FormValue("defaulton"))

	err = AddStatus(session, rcp)
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
