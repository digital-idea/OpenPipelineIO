package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dchest/captcha"
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

// handleAPIAddUser 함수는 사용자를 추가한다.
// 버튼을 눌렀을 때 ajax로 Post 한다.
// 사용자 가입이 되도록 한다.
func handleAPIAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	/*
		_, _, err = TokenHandler(r, session)
		if err != nil {
			fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
			return
		}
	*/

	r.ParseForm() // 파싱이 되면 맵이 됩니다.
	if !captcha.VerifyString(r.FormValue("CaptchaID"), r.FormValue("CaptchaNum")) {
		err := errors.New("CaptchaID 값과 CaptchaNum 값이 다릅니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(r.Form)
	/*


		u := User{}
		for key, values := range r.Form {
			switch key {
			case "ID":
				u.ID = values[0]
			case "Password":
				pw, err := Encrypt(values[0])
				if err != nil {
					fmt.Fprintf(w, "{\"error\":\"%s\"}\n", "")
					return
				}
				u.Password = pw
				u.Token = base64.StdEncoding.EncodeToString([]byte(pw))
			case "FirstNameKor":
				u.FirstNameKor = values[0]
			case "LastNameKor":
				u.LastNameKor = values[0]
			case "FirstNameEng":
				u.FirstNameEng = strings.Title(strings.ToLower(r.FormValue("FirstNameEng")))
			case "LastNameEng":
				u.LastNameEng = strings.Title(strings.ToLower(r.FormValue("LastNameEng")))
			case "FirstNameChn":
				u.FirstNameChn = r.FormValue("FirstNameChn")
			case "LastNameChn":
				u.LastNameChn = r.FormValue("LastNameChn")
			case "Email":
				u.Email = r.FormValue("Email")
			case "EmailExternal":
				u.EmailExternal = r.FormValue("EmailExternal")
			case "Phone":
				u.Phone = r.FormValue("Phone")
			case "Hotline":
				u.Hotline = r.FormValue("Hotline")
			case "Location":
				u.Location = r.FormValue("Location")
			case "Tags":
				u.Tags = Str2Tags(r.FormValue("Tags"))
			case "Timezone":
				u.Timezone = r.FormValue("Timezone")
			default:
				fmt.Println("key:", key)
				fmt.Println("value:", values)
			}
		}

		host, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.LastIP = host
		u.LastPort = port

		// 이 데이터가 DB로 들어가야 한다.
		fmt.Println(u)
	*/
	/*
		err = addUser(session, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = addToken(session, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		SetSessionID(w, u.ID, u.AccessLevel, "")
		// 가입완료 페이지로 이동
		t, err := LoadTemplates()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "signup_success", u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/
}
