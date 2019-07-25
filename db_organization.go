package main

import (
	"errors"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// addTeam 함수는 Part를 추가하는 함수이다.
func addTeam(session *mgo.Session, t Team) error {
	if t.ID == "" {
		err := errors.New("ID가 빈 문자열입니다. 유저를 생성할 수 없습니다")
		log.Println(err)
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("teams")

	num, err := c.Find(bson.M{"id": t.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		err = errors.New(t.ID + " ID를 가진 Team이 이미 DB에 존재합니다.")
		log.Println(err)
		return err
	}
	err = c.Insert(t)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// allDivisions 함수는 모든 Division을 반환한다.
func allDivisions(session *mgo.Session) ([]Division, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("divisions")
	var result []Division
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// allDepartments 함수는 모든 Department를 반환한다.
func allDepartments(session *mgo.Session) ([]Department, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("departments")
	var result []Department
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// allTeams 함수는 모든 Team을 반환한다.
func allTeams(session *mgo.Session) ([]Team, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("teams")
	var result []Team
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// allRoles 함수는 모든 Role을 반환한다.
func allRoles(session *mgo.Session) ([]Role, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("roles")
	var result []Role
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// allPositions 함수는 모든 Position 을 반환한다.
func allPositions(session *mgo.Session) ([]Position, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("positions")
	var result []Position
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// getTeam 함수는 Team 아이디를 받아서 Team 자료구조를 반환하는 함수이다.
func getTeam(session *mgo.Session, id string) (Team, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("teams")
	t := Team{}
	err := c.Find(bson.M{"id": id}).One(&t)
	if err != nil {
		return t, err
	}
	return t, nil
}
