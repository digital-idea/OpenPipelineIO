package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
)

// handleAPIUser 함수는 사용자의 id를 받아서 사용자 정보를 반환한다.
func handleAPIUser(w http.ResponseWriter, r *http.Request) {
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
	id := q.Get("id")
	user, err := getUser(session, id)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	type recipe struct {
		Data User `json:"data"`
	}
	rcp := recipe{}
	// 불필요한 정보는 초기화 시킨다.
	user.Password = ""
	user.Token = ""
	rcp.Data = user
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISearchUser 함수는 단어를 받아서 조건에 맞는 사용자 정보를 반환한다.
func handleAPISearchUser(w http.ResponseWriter, r *http.Request) {
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
	searchword := q.Get("searchword")
	users, err := searchUsers(session, strings.Split(searchword, ","))
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	type recipe struct {
		Data []User `json:"data"`
	}
	rcp := recipe{}
	// 불필요한 정보는 초기화 시킨다.
	for _, user := range users {
		user.Password = ""
		user.Token = ""
		rcp.Data = append(rcp.Data, user)
	}
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIValidUser 함수는 사용자가 유효한지 체크하는 핸들러 입니다.
func handleAPIValidUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
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
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var id string
	var pw string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "id":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			id = v
		case "pw":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			pw = v
		}
	}
	err = vaildUser(session, id, pw)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	fmt.Fprintf(w, "{\"error\":\"\"}\n")
}
