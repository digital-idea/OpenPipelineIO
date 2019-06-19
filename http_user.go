package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"gopkg.in/mgo.v2"
)

// handleUser 함수는 유저리스트 페이지이다.
func handleUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "진행예정2순위 : '17.5.15~'17.6.30")
}

// handleSignup 함수는 회원가입 페이지이다.
func handleSignup(w http.ResponseWriter, r *http.Request) {
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		Company   string
		CaptchaID string
		MailDNS   string
	}
	rcp := recipe{}
	rcp.Company = strings.Title(*flagCompany)
	rcp.MailDNS = *flagMailDNS
	rcp.CaptchaID = captcha.New()
	err = t.ExecuteTemplate(w, "signup", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSignupSubmit 함수는 회원가입 페이지이다.
func handleSignupSubmit(w http.ResponseWriter, r *http.Request) {
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if !captcha.VerifyString(r.FormValue("CaptchaNum"), r.FormValue("CaptchaID")) {
		err := errors.New("CaptchaNum 값과 CaptchaID 값이 다릅니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("ID") == "" {
		err := errors.New("ID 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("Password") == "" {
		err := errors.New("Password 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("Password") != r.FormValue("RePassword") {
		err := errors.New("입력받은 2개의 패스워드가 서로 다릅니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.FormValue("ID")
	u := *NewUser(id)
	pw, err := Encrypt(r.FormValue("Password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Password = pw
	u.FirstNameKor = r.FormValue("FirstNameKor")
	u.LastNameKor = r.FormValue("LastNameKor")
	u.FirstNameEng = r.FormValue("FirstNameEng")
	u.LastNameEng = r.FormValue("LastNameEng")
	u.FirstNameChn = r.FormValue("FirstNameChn")
	u.LastNameChn = r.FormValue("LastNameChn")
	u.Email = r.FormValue("Email")
	u.EmailExternal = r.FormValue("EmailExternal")
	u.Phone = r.FormValue("Phone")
	u.Hotline = r.FormValue("Hotline")
	u.Location = r.FormValue("Location")
	u.Parts = Str2Tags(r.FormValue("Parts"))
	u.Timezone = r.FormValue("Timezone")
	u.LastIP = host
	u.LastPort = port
	// 이 데이터가 DB로 들어가야 한다.
	session, err := mgo.DialWithTimeout(*flagDBIP, 2*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = addUser(session, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 쿠키설정
	expire := time.Now().Add(1 * 4 * time.Hour) //MPAA기준이 4시간이다.
	cookie := http.Cookie{
		Name:    "ID",
		Value:   u.ID,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)

	// 가입완료 페이지로 이동
	err = t.ExecuteTemplate(w, "signup_success", u)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSignin 함수는 로그인 페이지이다.
func handleSignin(w http.ResponseWriter, r *http.Request) {
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		Company string
	}
	rcp := recipe{}
	rcp.Company = strings.Title(*flagCompany)
	err = t.ExecuteTemplate(w, "signin", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSigninSubmit 함수는 회원가입 페이지이다.
func handleSigninSubmit(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("ID") == "" {
		err := errors.New("ID 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("Password") == "" {
		err := errors.New("Password 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.FormValue("ID")
	pw := r.FormValue("Password")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	err = vaildUser(session, id, pw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 유효하면 ID 쿠키를 셋팅하고 / 로 이동한다.
	expire := time.Now().Add(1 * 4 * time.Hour) //MPAA기준이 4시간이다.
	cookie := http.Cookie{
		Name:    "ID",
		Value:   id,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// handleSignout 함수는 로그아웃 페이지이다.
func handleSignout(w http.ResponseWriter, r *http.Request) {
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:   "ID",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	err = t.ExecuteTemplate(w, "signout", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUsersInfo 함수는 유저 자료구조 페이지이다.
func handleUserinfo(w http.ResponseWriter, r *http.Request) {
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		Users []User
	}
	rcp := recipe{}
	rcp.Users, err = getUsers(session)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "userinfo", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
