package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
)

// handleAPI2User 함수는 사용자관련 REST API이다. GET, DELETE를 지원한다.
func handleAPI2User(w http.ResponseWriter, r *http.Request) {
	// GET 메소드는 사용자의 id를 받아서 사용자 정보를 반환한다.
	if r.Method == http.MethodGet {
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer session.Close()
		_, _, err = TokenHandler(r, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
			return
		}
		user, err := getUser(session, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 불필요한 정보는 초기화 시킨다.
		user.Password = ""
		user.Token = ""
		// json 으로 결과 전송
		data, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
		// DELETE 메소드는 사용자의 ID를 받아 해당 사용자를 DB에서 삭제한다.
	} else if r.Method == http.MethodDelete {
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer session.Close()
		// accesslevel 체크. user 삭제는 admin만 가능하다.
		_, accesslevel, err := TokenHandler(r, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if int(accesslevel) < 11 {
			http.Error(w, "permission is low", http.StatusUnauthorized)
			return
		}
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
			return
		}
		// 토큰 삭제
		err = rmToken(session, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 유저 삭제
		err = rmUser(session, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//responce
		data, err := json.Marshal("deleted")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	} else {
		http.Error(w, "Not Supported Method", http.StatusMethodNotAllowed)
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

// handleAPISetLeaveUser 함수는 사용자의 퇴사여부를 셋팅하는 핸들러 입니다.
func handleAPISetLeaveUser(w http.ResponseWriter, r *http.Request) {
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
	var leave string
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
		case "leave":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			leave = v
		}
	}
	err = setLeaveUser(session, id, str2bool(leave))
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	fmt.Fprintf(w, "{\"error\":\"\"}\n")
}

// handleAPIAutoCompliteUsers 함수는 form에서 autocomplite 에 사용되는 사용자 데이터를 반환한다.
func handleAPIAutoCompliteUsers(w http.ResponseWriter, r *http.Request) {
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
	if *flagAuthmode { // 보안모드로 실행하면, 철저하게 검사해야 한다.
		_, accesslevel, err := TokenHandler(r, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if int(accesslevel) < 3 {
			http.Error(w, "permission is low", http.StatusUnauthorized)
			return
		}
	}
	users, err := allUsers(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type AutocompliteUser struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Team       string `json:"team"`
		Searchword string `json:"searchword"`
	}
	type recipe struct {
		Users []AutocompliteUser `json:"users"`
	}
	rcp := recipe{}
	for _, user := range users {
		id := user.ID
		name := user.LastNameKor + user.FirstNameKor
		var team string
		for _, o := range user.Organizations {
			if o.Primary {
				team = o.Team.Name
				break
			}
			team = o.Team.Name
		}
		u := AutocompliteUser{
			ID:         id,
			Name:       name,
			Team:       team,
			Searchword: id + name + team,
		}
		rcp.Users = append(rcp.Users, u)
	}
	// json 으로 결과 전송
	data, _ := json.Marshal(rcp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIInitPassword 함수는 User의 패스워드를 Adminsetting에 설정된 패스워드를 이용해서 패스워드를 리셋한다.
func handleAPIInitPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		ID          string      `json:"id"`
		AccessLevel AccessLevel `json:"accesslevel"`
		UserID      string      `json:"userid"`
	}
	rcp := Recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.UserID, rcp.AccessLevel, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// 관리자만 이 API를 사용할 수 있도록 제한합니다.
	if rcp.AccessLevel != AdminAccessLevel {
		http.Error(w, "사용자의 패스워드를 초기화 하기 위해서 관리자 권한이 필요합니다", http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	// 사용자의 패스워드를 초기화 합니다.
	err = initPassUser(session, rcp.ID, CachedAdminSetting.InitPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 바뀐 패스워드의 유저를 불러옵니다.
	u, err := getUser(session, rcp.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 기존 토근을 제거합니다.
	err = rmToken(session, u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 유저에 새로운 토큰을 추가합니다.
	err = addToken(session, u)
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
