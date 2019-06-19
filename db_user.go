// 이 코드는 사용자와 관련된 DBapi가 모여있는 파일입니다.

package main

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// addUser 함수는 사용자를 추가하는 함수이다.
func addUser(session *mgo.Session, u User) error {
	if u.ID == "" {
		err := errors.New("ID가 빈 문자열입니다. 유저를 생성할 수 없습니다")
		log.Println(err)
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")

	num, err := c.Find(bson.M{"id": u.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		err = errors.New(u.ID + " ID를 가진 사용자가 이밎 DB에 존재합니다.")
		log.Println(err)
		return err
	}
	u.Createtime = time.Now().Format(time.RFC3339)
	err = c.Insert(u)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// getUser 함수는 사용자를 가지고오는 함수이다.
func getUser(session *mgo.Session, id string) (User, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	u := User{}
	err := c.Find(bson.M{"id": id}).One(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// rmUser 함수는 사용자를 삭제하는 함수이다.
func rmUser(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// setUser 함수는 사용자 정보를 업데이트하는 함수이다.
func setUser(session *mgo.Session, u User) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	num, err := c.Find(bson.M{"id": u.ID}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	u.Updatetime = time.Now().Format(time.RFC3339)
	err = c.Update(bson.M{"id": u.ID}, u)
	if err != nil {
		return err
	}
	return nil
}

// initPassUser 함수는 사용자 정보를 수정하는 함수이다.
func initPassUser(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	num, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	q := bson.M{"id": id}
	encryptPass, err := Encrypt("Welcome2csi!")
	if err != nil {
		return err
	}
	change := bson.M{"$set": bson.M{"password": encryptPass, "updatetime": time.Now().Format(time.RFC3339)}}
	err = c.Update(q, change)
	if err != nil {
		return err
	}
	return nil
}

// updatePasswordUser 함수는 사용자 패스워드를 수정하는 함수이다.
func updatePasswordUser(session *mgo.Session, id, pw, newPw string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	num, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	q := bson.M{"id": id}
	// 과거의 패스워드로 로그인가능했는지 체크한다.
	err = vaildUser(session, id, pw)
	if err != nil {
		return err
	}
	// 새로운 패스워드로 업데이트 한다.
	encryptPass, err := Encrypt(newPw)
	if err != nil {
		return err
	}
	change := bson.M{"$set": bson.M{"password": encryptPass, "updatetime": time.Now().Format(time.RFC3339)}}
	err = c.Update(q, change)
	if err != nil {
		return err
	}
	return nil
}

// getUsers 함수는 DB에서 전체 사용자 정보를 가지고오는 함수입니다.
func getUsers(session *mgo.Session) ([]User, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	var result []User
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// vaildUser 함수는 사용자의 id, pw를 받아서 유효한 사용자인지 체크한다.
func vaildUser(session *mgo.Session, id, pw string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	q := bson.M{"id": id}
	num, err := c.Find(q).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	u := User{}
	err = c.Find(q).One(&u)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	if err != nil {
		return err
	}
	return nil
}
