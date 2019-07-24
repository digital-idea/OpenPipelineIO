package main

import (
	"errors"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// addPart 함수는 Part를 추가하는 함수이다.
func addPart(session *mgo.Session, p Part) error {
	if p.ID == "" {
		err := errors.New("ID가 빈 문자열입니다. 유저를 생성할 수 없습니다")
		log.Println(err)
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("part").C("parts")

	num, err := c.Find(bson.M{"id": p.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		err = errors.New(p.ID + " ID를 가진 Part가 이미 DB에 존재합니다.")
		log.Println(err)
		return err
	}
	err = c.Insert(p)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// allParts 함수는 모든 Part를 반환한다.
func allParts(session *mgo.Session) ([]Part, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("part").C("parts")
	var result []Part
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
