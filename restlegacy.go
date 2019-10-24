package main

import (
	"fmt"
	"net/http"
	"strings"
	"gopkg.in/mgo.v2"
)

// handleAPISetPmnotes 함수는 아이템에 수정사항 리스트를 교체합니다.
func handleAPISetComments(w http.ResponseWriter, r *http.Request) {
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
	userID, _, err := TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var text string
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
		case "text":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			text = v
		}
	}
	var texts []string
	// 문장의 끝을 의미하는 마침표를 이용하여 각 문장을 분리한다.
	for _, text := range strings.Split(text, ".") {
		if text == "" {
			continue
		}
		texts = append(texts, strings.TrimSpace(text)+".")
	}
	// 문장이 통째로 바뀔 때는 순서대로 보여야 한다. Reverse 한다.
	for i, j := 0, len(texts)-1; i < j; i, j = i+1, j-1 {
		texts[i], texts[j] = texts[j], texts[i]
	}

	err = SetComments(session, project, name, userID, texts)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	fmt.Fprintf(w, "{\"error\":\"\"}\n")
}
