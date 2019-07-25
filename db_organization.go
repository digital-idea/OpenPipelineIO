package main

import (
	"errors"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// addDivision 함수는 Division을 추가하는 함수이다.
func addDivision(session *mgo.Session, d Division) error {
	if d.ID == "" {
		err := errors.New("ID가 빈 문자열입니다. Division 을 생성할 수 없습니다")
		log.Println(err)
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("divisions")

	num, err := c.Find(bson.M{"id": d.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		err = errors.New(d.ID + " ID를 가진 Division이 이미 DB에 존재합니다.")
		log.Println(err)
		return err
	}
	err = c.Insert(d)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// addDepartment 함수는 Department 를 추가하는 함수이다.
func addDepartment(session *mgo.Session, d Department) error {
	if d.ID == "" {
		err := errors.New("ID가 빈 문자열입니다. Department 를 생성할 수 없습니다")
		log.Println(err)
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("departments")

	num, err := c.Find(bson.M{"id": d.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		err = errors.New(d.ID + " ID를 가진 Department 가 이미 DB에 존재합니다.")
		log.Println(err)
		return err
	}
	err = c.Insert(d)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// addTeam 함수는 Team을 추가하는 함수이다.
func addTeam(session *mgo.Session, t Team) error {
	if t.ID == "" {
		err := errors.New("ID가 빈 문자열입니다. Team을 생성할 수 없습니다")
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

// addRole 함수는 Role을 추가하는 함수이다.
func addRole(session *mgo.Session, r Role) error {
	if r.ID == "" {
		err := errors.New("ID가 빈 문자열입니다. Role 을 생성할 수 없습니다")
		log.Println(err)
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("roles")

	num, err := c.Find(bson.M{"id": r.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		err = errors.New(r.ID + " ID를 가진 Role이 이미 DB에 존재합니다.")
		log.Println(err)
		return err
	}
	err = c.Insert(r)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// addPosition 함수는 Position을 추가하는 함수이다.
func addPosition(session *mgo.Session, p Position) error {
	if p.ID == "" {
		err := errors.New("ID가 빈 문자열입니다. Position 을 생성할 수 없습니다")
		log.Println(err)
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("positions")

	num, err := c.Find(bson.M{"id": p.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		err = errors.New(p.ID + " ID를 가진 Position이 이미 DB에 존재합니다.")
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
