package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2"
)

// 이 파일은 restAPI가 설정되는 페이지이다.

// handleAPIItem 함수는 아이템 자료구조를 불러온다.
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
	err = TokenHandler(r, session)
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

// handleAPIRmItem 함수는 아이템을 삭제한다.
func handleAPIRmItem(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}

	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var typ string
	info := r.PostForm
	for key, value := range info {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "type":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			typ = v
		default:
			fmt.Fprintf(w, "{\"error\":\"%s\"}\n", key+"키는 사용할 수 없습니다.(project, name, type 키값만 사용가능합니다.)")
			return
		}
	}
	err = rmItem(session, project, name, typ)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	fmt.Fprintf(w, "{\"error\":\"%s\"}\n", "")
}

// handleAPISearchname 함수는 입력 문자열을 포함하는 샷,에셋 정보를 검색한다.
func handleAPISearchname(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	name := q.Get("name")
	items, err := SearchName(session, project, name)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	type recipe struct {
		Data []Item `json:"data"`
	}
	rcp := recipe{}
	rcp.Data = items
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPI2Items 함수는 아이템을 검색한다.
func handleAPI2Items(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	type recipe struct {
		Data []Item `json:"data"`
	}
	rcp := recipe{}
	q, err := URLUnescape(r.URL)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	op := SearchOption{
		Project:    q.Get("project"),
		Searchword: q.Get("searchword"),
		Sortkey:    q.Get("sortkey"),
		Assign:     str2bool(q.Get("assign")),
		Ready:      str2bool(q.Get("ready")),
		Wip:        str2bool(q.Get("wip")),
		Confirm:    str2bool(q.Get("confirm")),
		Done:       str2bool(q.Get("done")),
		Omit:       str2bool(q.Get("omit")),
		Hold:       str2bool(q.Get("hold")),
		Out:        str2bool(q.Get("out")),
		None:       str2bool(q.Get("none")),
		Shot:       str2bool(q.Get("shot")),
		Assets:     str2bool(q.Get("shot")),
		Type3d:     str2bool(q.Get("type3d")),
		Type2d:     str2bool(q.Get("type2d")),
	}
	result, err := Search(session, op)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	rcp.Data = result
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIShot 함수는 project, name을 받아서 shot을 반환한다.
func handleAPIShot(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	name := q.Get("name")
	type recipe struct {
		Data Item `json:"data"`
	}
	rcp := recipe{}
	result, err := Shot(session, project, name)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	rcp.Data = result
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISeqs 함수는 프로젝트의 시퀀스를 가져온다.
func handleAPISeqs(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
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
	seqs, err := Seqs(session, project)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	type recipe struct {
		Data []string `json:"data"`
	}
	rcp := recipe{}
	rcp.Data = seqs
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
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
	err = TokenHandler(r, session)
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

// handleAPISetTaskMov 함수는 Task에 mov를 설정한다.
func handleAPISetTaskMov(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var task string
	var mov string
	info := r.PostForm
	for key, value := range info {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name", "shot", "asset":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "task":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			err = validTask(v)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%s\"}\n", err.Error())
			}
			task = v
		case "mov": // 앞뒤샷 포함 여러개의 mov를 등록할 수 있다.
			mov = strings.Join(value, ";")
		default:
			fmt.Fprintf(w, "{\"error\":\"%s\"}\n", key+"키는 사용할 수 없습니다.(project, shot, asset, task, mov 키값만 사용가능합니다.)")
			return
		}
	}
	if mov == "" {
		fmt.Fprintf(w, "{\"error\":\"%s\"}\n", "mov 키값을 설정해주세요.")
		return
	}
	if *flagDebug {
		fmt.Println(project, name, task, mov)
	}
	err = setTaskMov(session, project, name, task, mov)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIRenderSize 함수는 아이템에 RenderSize를 설정한다.
func handleAPISetRenderSize(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var size string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "size":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			if !regexpImageSize.MatchString(v) {
				fmt.Fprintf(w, "{\"error\":\"%s\"}\n", "size 는 2048x1152 형태여야 합니다.")
				return
			}
			size = v
		}
	}
	err = SetImageSize(session, project, name, "rendersize", size)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIDistortionSize 함수는 아이템의 DistortionSize를 설정한다.
func handleAPISetDistortionSize(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var size string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "size":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			if !regexpImageSize.MatchString(v) {
				fmt.Fprintf(w, "{\"error\":\"%s\"}\n", "size 는 2048x1152 형태여야 합니다.")
				return
			}
			size = v
		}
	}
	err = SetImageSize(session, project, name, "dsize", size)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetJustIn 함수는 아이템에 JustIn 값을 설정한다.
func handleAPISetJustIn(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var frame int
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "frame":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			n, err := strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			frame = n
		}
	}
	err = SetFrame(session, project, name, "justin", frame)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetJustOut 함수는 아이템에 JustOut 값을 설정한다.
func handleAPISetJustOut(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var frame int
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "frame":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			n, err := strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			frame = n
		}
	}
	err = SetFrame(session, project, name, "justout", frame)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIPlateSize 함수는 아이템의 PlateSize를 설정한다.
func handleAPISetPlateSize(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var size string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "size":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			if !regexpImageSize.MatchString(v) {
				fmt.Fprintf(w, "{\"error\":\"%s\"}\n", "size 는 2048x1152 형태여야 합니다")
				return
			}
			size = v
		}
	}
	err = SetImageSize(session, project, name, "platesize", size)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// PostFormValueInList 는 PostForm 쿼리시 Value값이 1개라면 값을 리턴한다.
func PostFormValueInList(key string, values []string) (string, error) {
	if len(values) != 1 {
		return "", errors.New(key + "값이 여러개 입니다")
	}
	if key == "startdate" && values[0] == "" { // Task 시작일은 빈 문자를 허용한다.
		return "", nil
	}
	if key == "predate" && values[0] == "" { // 1차마감일은 빈 문자를 허용한다.
		return "", nil
	}
	if key == "date" && values[0] == "" { // 2차마감일은 빈 문자를 허용한다.
		return "", nil
	}
	if values[0] == "" {
		return "", errors.New(key + "값이 빈 문자입니다")
	}
	return values[0], nil
}

// handleAPISetCameraPubPath 함수는 아이템의 Camera PubPath를 설정한다.
func handleAPISetCameraPubPath(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var path string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "path":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			path = v
		}
	}
	err = SetCameraPubPath(session, project, name, path)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetCameraPubTask 함수는 아이템의 Camera PubTask를 설정한다.
func handleAPISetCameraPubTask(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var task string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "task":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			task = v
		}
	}
	err = SetCameraPubTask(session, project, name, task)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetCameraProjection 함수는 아이템의 Camera Projection 여부를 설정한다.
func handleAPISetCameraProjection(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var projection bool
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "projection":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			if strings.ToLower(v) == "true" {
				projection = true
			}
		}
	}
	err = SetCameraProjection(session, project, name, projection)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetThummov 함수는 아이템의 Thummov 값을 설정한다.
func handleAPISetThummov(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var path string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "path":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			path = v
		}
	}
	err = SetThummov(session, project, name, path)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetTaskStatus 함수는 아이템의 task에 대한 상태를 설정한다.
func handleAPISetTaskStatus(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var task string
	var status string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "task":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			task = v
		case "status":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			status = v
		}
	}
	err = SetTaskStatus(session, project, name, task, status)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetUser 함수는 아이템의 task에 대한 유저를 설정한다.
func handleAPISetTaskUser(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var task string
	var user string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "task":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			task = v
		case "user":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			user = v
		}
	}
	err = SetTaskUser(session, project, name, task, user)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetTaskStartdate 함수는 아이템의 task에 대한 시작일을 설정한다.
func handleAPISetTaskStartdate(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var task string
	var date string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "task":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			task = v
		case "date", "startdate":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			date = v
		}
	}
	err = SetTaskStartdate(session, project, name, task, date)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetTaskPredate 함수는 아이템의 task에 대한 1차마감일을 설정한다.
func handleAPISetTaskPredate(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var task string
	var date string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "task":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			task = v
		case "date", "predate":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			date = v
		}
	}
	err = SetTaskPredate(session, project, name, task, date)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetTaskDate 함수는 아이템의 task에 대한 최종마감일을 설정한다.
func handleAPISetTaskDate(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var task string
	var date string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "task":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			task = v
		case "date":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			date = v
		}
	}
	err = SetTaskDate(session, project, name, task, date)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetShotType 함수는 아이템의 shot type을 설정한다.
func handleAPISetShotType(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var typ string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "type", "shottype":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			typ = v
		}
	}
	err = SetShotType(session, project, name, typ)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPISetOutputName 함수는 아이템의 shot의 아웃풋 이름을 설정합니다.
func handleAPISetOutputName(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var outputname string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "outputname":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			outputname = v
		}
	}
	err = SetOutputName(session, project, name, outputname)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	fmt.Fprintf(w, "{\"error\":\"\"}\n")
}

// handleAPISetAssetType 함수는 아이템의 shot type을 설정한다.
func handleAPISetAssetType(w http.ResponseWriter, r *http.Request) {
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
	err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var typ string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "type", "assettype":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			typ = v
		}
	}
	err = SetAssetType(session, project, name, typ)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}
