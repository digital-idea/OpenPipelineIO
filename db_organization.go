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

// allDivisions 함수는 DB에서 전체 사용자 정보를 가지고오는 함수입니다.
func allDivisions(session *mgo.Session) ([]Division, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("divisions")
	var result []Division
	err := c.Find(bson.M{}).Sort("id").All(&result)
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
	err := c.Find(bson.M{}).Sort("id").All(&result)
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
	err := c.Find(bson.M{}).Sort("id").All(&result)
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
	err := c.Find(bson.M{}).Sort("id").All(&result)
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
	err := c.Find(bson.M{}).Sort("id").All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// getDivision 함수는 본부를 가지고오는 함수이다.
func getDivision(session *mgo.Session, id string) (Division, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("divisions")
	d := Division{}
	err := c.Find(bson.M{"id": id}).One(&d)
	if err != nil {
		return d, err
	}
	return d, nil
}

// getDepartment 함수는 Department 아이디를 받아서 Department 자료구조를 반환하는 함수이다.
func getDepartment(session *mgo.Session, id string) (Department, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("departments")
	d := Department{}
	err := c.Find(bson.M{"id": id}).One(&d)
	if err != nil {
		return d, err
	}
	return d, nil
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

// getRole 함수는 Role 아이디를 받아서 Role 자료구조를 반환하는 함수이다.
func getRole(session *mgo.Session, id string) (Role, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("roles")
	r := Role{}
	err := c.Find(bson.M{"id": id}).One(&r)
	if err != nil {
		return r, err
	}
	return r, nil
}

// getPosition 함수는 Position 아이디를 받아서 Position 자료구조를 반환하는 함수이다.
func getPosition(session *mgo.Session, id string) (Position, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("positions")
	p := Position{}
	err := c.Find(bson.M{"id": id}).One(&p)
	if err != nil {
		return p, err
	}
	return p, nil
}

// setDivision 함수는 Division 정보를 수정하는 함수입니다.
func setDivision(session *mgo.Session, d Division) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("divisions")
	num, err := c.Find(bson.M{"id": d.ID}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 아이템이 존재하지 않습니다")
	}
	err = c.Update(bson.M{"id": d.ID}, d)
	if err != nil {
		return err
	}
	return nil
}

// setDepartment 함수는 Department 정보를 수정하는 함수입니다.
func setDepartment(session *mgo.Session, d Department) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("departments")
	num, err := c.Find(bson.M{"id": d.ID}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 아이템이 존재하지 않습니다")
	}
	err = c.Update(bson.M{"id": d.ID}, d)
	if err != nil {
		return err
	}
	return nil
}

func setTeam(session *mgo.Session, t Team) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("teams")
	num, err := c.Find(bson.M{"id": t.ID}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 아이템이 존재하지 않습니다")
	}
	err = c.Update(bson.M{"id": t.ID}, t)
	if err != nil {
		return err
	}
	return nil
}

func setRole(session *mgo.Session, r Role) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("roles")
	num, err := c.Find(bson.M{"id": r.ID}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 아이템이 존재하지 않습니다")
	}
	err = c.Update(bson.M{"id": r.ID}, r)
	if err != nil {
		return err
	}
	return nil
}

func setPosition(session *mgo.Session, p Position) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("positions")
	num, err := c.Find(bson.M{"id": p.ID}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 아이템이 존재하지 않습니다")
	}
	err = c.Update(bson.M{"id": p.ID}, p)
	if err != nil {
		return err
	}
	return nil
}

// rmDivision 함수는 Division을 삭제하는 함수이다.
func rmDivision(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("divisions")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// rmDepartment 함수는 Department를 삭제하는 함수이다.
func rmDepartment(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("departments")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// rmTeam 함수는 Team을 삭제하는 함수이다.
func rmTeam(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("teams")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// rmRole 함수는 Role을 삭제하는 함수이다.
func rmRole(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("roles")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// rmPosition 함수는 Position을 삭제하는 함수이다.
func rmPosition(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("organization").C("positions")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}
