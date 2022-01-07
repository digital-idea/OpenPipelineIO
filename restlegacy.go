package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/digital-idea/dilog"
	"github.com/digital-idea/dipath"
	"gopkg.in/mgo.v2"
)

// handleAPISetTags 함수는 아이템에 태그를 교체합니다.
func handleAPISetTags(w http.ResponseWriter, r *http.Request) {
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
	tokenID, _, err := TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var tags string
	args := r.PostForm
	for key, values := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "tag", "tags":
			if len(values) == 0 {
				tags = ""
			} else {
				tags = values[0]
			}
		}
	}
	err = SetTags(session, project, name, Str2List(tags))
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	// 로그처리
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dilog.Add(*flagDBIP, host, fmt.Sprintf("Edit tag: %s", tags), project, name, "csi3", tokenID, 180)
	fmt.Fprintf(w, "{\"error\":\"\"}\n")
}

// handleAPISetTaskMov 함수는 Task에 mov를 설정한다.
func handleAPISetTaskMov(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Name    string `json:"name"`
		Task    string `json:"task"`
		Mov     string `json:"mov"`
		UserID  string `json:"userid"`
		Error   string `json:"error"`
	}
	rcp := Recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.UserID, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	for key, values := range r.PostForm {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Project = v
		case "name", "shot", "asset":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Name = v
		case "task":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Task = v
		case "mov": // 앞뒤샷 포함 여러개의 mov를 등록할 수 있다.
			rcp.Mov = strings.Join(values, ";")
		case "userid":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if rcp.UserID == "unknown" && v != "" {
				rcp.UserID = v
			}
		default:
			http.Error(w, key+"키는 사용할 수 없습니다.(project, shot, asset, task, mov 키값만 사용가능합니다.)", http.StatusBadRequest)
			return
		}
	}
	rcp.Mov = dipath.Win2lin(rcp.Mov) // 내부적으로 모든 경로는 unix 경로를 사용한다.
	_, err = setTaskMov(session, rcp.Project, rcp.Name, rcp.Task, rcp.Mov)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// log
	err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Setmov: %s %s", rcp.Task, rcp.Mov), rcp.Project, rcp.Name, "csi3", rcp.UserID, 180)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Setmov: %s %s\nProject: %s, Name: %s, Author: %s", rcp.Task, rcp.Mov, rcp.Project, rcp.Name, rcp.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, _ := json.Marshal(rcp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIUser 함수는 사용자관련 REST API이다. GET, DELETE를 지원한다. // legacy
func handleAPIUser(w http.ResponseWriter, r *http.Request) {
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
		type recipe struct {
			Data User `json:"data"`
		}
		rcp := recipe{}
		// 불필요한 정보는 초기화 시킨다.
		user.Password = ""
		user.Token = ""
		rcp.Data = user
		// json 으로 결과 전송
		data, err := json.Marshal(rcp)
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

// handleAPIItem 함수는 아이템 자료구조를 불러온다. // legacy
func handleAPIItem(w http.ResponseWriter, r *http.Request) {
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
	if project == "" {
		fmt.Fprintf(w, "{\"error\":\"%s\"}\n", "project가 빈 문자열입니다.")
		return
	}
	// 과거 slug 방식의 api에서 id 방식의 api로 변경중이다.
	var id string
	if q.Get("slug") != "" {
		id = q.Get("slug")
	} else if q.Get("id") != "" {
		id = q.Get("id")
	}
	if id == "" {
		fmt.Fprintf(w, "{\"error\":\"%s\"}\n", "id 또는 slug가 빈 문자열입니다.")
		return
	}
	item, err := getItem(session, project, id)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	return
}

// handleAPISetThummov 함수는 아이템의 Thummov 값을 설정한다.
func handleAPISetThummov(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Name    string `json:"name"`
		Path    string `json:"path"`
		UserID  string `json:"userid"`
		Error   string `json:"error"`
	}
	rcp := Recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.UserID, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	for key, values := range r.PostForm {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			rcp.Project = v
		case "name":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Name = v
		case "userid":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if rcp.UserID == "unknown" && v != "" {
				rcp.UserID = v
			}
		case "path":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Path = v
		}
	}
	err = SetThummov(session, rcp.Project, rcp.Name, rcp.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// log
	err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Set Thumbnail: %s", rcp.Path), rcp.Project, rcp.Name, "csi3", rcp.UserID, 180)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Set Thumbnail: %s\nProject: %s, Name: %s, Author: %s", rcp.Path, rcp.Project, rcp.Name, rcp.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, _ := json.Marshal(rcp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIRenderSize 함수는 아이템에 RenderSize를 설정한다.
func handleAPISetRenderSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Name    string `json:"name"`
		ID      string `json:"id"`
		Size    string `json:"size"`
		UserID  string `json:"userid"`
		Error   string `json:"error"`
	}
	rcp := Recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.UserID, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	for key, values := range r.PostForm {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Project = v
		case "name":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Name = v
		case "userid":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if rcp.UserID == "unknown" && v != "" {
				rcp.UserID = v
			}
		case "size":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if !regexpImageSize.MatchString(v) {
				http.Error(w, "2048x1152 형태로 입력해주세요", http.StatusBadRequest)
				return
			}
			rcp.Size = v
		}
	}
	id, err := SetImageSize(session, rcp.Project, rcp.Name, "rendersize", rcp.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.ID = id
	// log
	err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Set Rendersize: %s", rcp.Size), rcp.Project, rcp.Name, "csi3", rcp.UserID, 180)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Set Rendersize: %s\nProject: %s, Name: %s, Author: %s", rcp.Size, rcp.Project, rcp.Name, rcp.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, _ := json.Marshal(rcp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIShots 함수는 project, seq를 입력받아서 샷정보를 출력한다.
func handleAPIShots(w http.ResponseWriter, r *http.Request) {
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
	if project == "" {
		fmt.Fprintln(w, "{\"error\":\"project 정보가 없습니다.\"}")
		return
	}
	seq := q.Get("seq")
	if seq == "" {
		fmt.Fprintln(w, "{\"error\":\"seq 정보가 없습니다.\"}")
		return
	}
	shots, err := Shots(session, project, seq)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	type recipe struct {
		Data []string `json:"data"`
	}
	rcp := recipe{}
	rcp.Data = shots
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetRnum 함수는 아이템에 롤넘버를 설정합니다.
func handleAPISetRnum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Name    string `json:"name"`
		ID      string `json:"id"`
		Rnum    string `json:"rnum"`
		UserID  string `json:"userid"`
		Error   string `json:"error"`
	}
	rcp := Recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.UserID, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	for key, values := range r.PostForm {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Project = v
		case "name":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Name = v
		case "userid":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if rcp.UserID == "unknown" && v != "" {
				rcp.UserID = v
			}
		case "rnum":
			if len(values) == 0 {
				rcp.Rnum = ""
			} else {
				rcp.Rnum = values[0]
			}
		}
	}
	if rcp.Rnum != "" && !regexpRnum.MatchString(rcp.Rnum) {
		http.Error(w, rcp.Rnum+"값은 A0001 형식이 아닙니다.", http.StatusBadRequest)
		return
	}
	typ, err := Type(session, rcp.Project, rcp.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = SetRnum(session, rcp.Project, rcp.Name+"_"+typ, rcp.Rnum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.ID = rcp.Name + "_" + typ
	// log
	err = dilog.Add(*flagDBIP, host, "Set Rnum: "+rcp.Rnum, rcp.Project, rcp.Name, "csi3", rcp.UserID, 180)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Set Rnum: %s\nProject: %s, Name: %s, Author: %s", rcp.Rnum, rcp.Project, rcp.Name, rcp.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, _ := json.Marshal(rcp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
