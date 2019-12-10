package main

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

// SetNotes 함수는 item에 작업,현장내용을 추가한다.  //legacy
func SetNotes(session *mgo.Session, project, name, userID string, texts []string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	i, err := getItem(session, project, id)
	if err != nil {
		return err
	}
	// 이 부분은 나중에 좋은 구조로 다시 바꾸어야 한다. 호환성을 위해서 현재는 CSI1의 구조로 현장노트를 입력한다.
	i.Onsetnote = []string{}
	for _, text := range texts {
		if text == "" {
			continue
		}
		note := fmt.Sprintf("%s;restapi;%s;%s", time.Now().Format(time.RFC3339), userID, text)
		i.Onsetnote = append(i.Onsetnote, note)
	}
	err = setItem(session, project, i)
	if err != nil {
		return err
	}
	return nil
}

// AddNote 함수는 item에 작업,현장내용을 추가한다.  //legacy
func AddNote(session *mgo.Session, project, name, userID, text string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	i, err := getItem(session, project, id)
	if err != nil {
		return err
	}
	// 이 부분은 나중에 좋은 구조로 다시 바꾸어야 한다. 호환성을 위해서 현재는 CSI1의 구조로 현장노트를 입력한다.
	note := fmt.Sprintf("%s;restapi;%s;%s", time.Now().Format(time.RFC3339), userID, text)
	i.Onsetnote = append(i.Onsetnote, note)
	err = setItem(session, project, i)
	if err != nil {
		return err
	}
	return nil
}

// RmNote 함수는 item에 작업내용을 삭제한다. //legacy
func RmNote(session *mgo.Session, project, name, userID, text string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	i, err := getItem(session, project, id)
	if err != nil {
		return err
	}
	// 이 부분은 나중에 좋은 구조로 다시 바꾸어야 한다. 호환성을 위해서 현재는 CSI1의 구조로 현장노트를 입력한다.
	var notes []string
	for _, note := range i.Onsetnote {
		i := strings.LastIndex(note, ";")
		if i == -1 {
			continue
		}
		if note[i+1:] == text {
			continue
		}
		notes = append(notes, note)
	}
	i.Onsetnote = notes
	err = setItem(session, project, i)
	if err != nil {
		return err
	}
	return nil
}
