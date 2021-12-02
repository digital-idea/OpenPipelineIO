package main

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

func handlePartners(w http.ResponseWriter, r *http.Request) {
	// 사용자 로그인 여부 확인
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	// mongoDB 연결
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	type recipe struct {
		User                 // 로그인한 사용자 정보
		Partners   []Partner // 검색된 사용자를 담는 리스트
		Tags       []string  // 부서,파트 태그
		Searchword string    // searchform에 들어가는 문자
		Partnernum int       // 검색된 파트너사수
		Devmode    bool      // 개발모드
		SearchOption
		Setting
	}
	rcp := recipe{}
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Setting = CachedAdminSetting
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "partners", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func handleAddPartner(w http.ResponseWriter, r *http.Request) {
	// 사용자 로그인 여부 확인
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	// mongoDB 연결
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	type recipe struct {
		User // 로그인한 사용자 정보
		Setting
	}
	rcp := recipe{}
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Setting = CachedAdminSetting
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addpartner", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
