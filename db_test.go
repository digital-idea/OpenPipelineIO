package main

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"
)

func Test_projeclist(t *testing.T) {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	// TravisCI 테스트용 프로젝트를 생성한다.
	testName := "TravisCI"
	p := *NewProject(testName)
	err = addProject(session, p)
	if err != nil {
		log.Fatal(err)
	}
	// 프로젝트 리스트를 가지고 온다.
	plist, err := Projectlist(session)
	if err != nil {
		t.Fatal(err)
	}
	if len(plist) == 0 {
		t.Fatal("프로젝트가 존재하지 않습니다. 테스트를 할 수 없습니다.")
	}
	// TravisCI 테스트가 끝나면 다시 프로젝트를 제거한다.
	err = rmProject(session, testName)
	if err != nil {
		log.Fatal(err)
	}
	return
}
