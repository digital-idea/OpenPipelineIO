// 이 코드는 사용자와 관련된 DBapi가 모여있는 파일입니다.

package main

import (
	"errors"
	"log"

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
	c := session.DB("user").C("active") // inactive(비활성화계정), active(활성화계정)

	num, err := c.Find(bson.M{"id": u.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		err = errors.New(u.ID + " ID가 DB에 이미 존재합니다.")
		log.Println(err)
		return err
	}
	u.Createtime = Now()
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
	c := session.DB("user").C("acitive")
	u := User{}
	err := c.Find(bson.M{"id": id}).One(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// rmUser 함수는 사용자를 삭제하는 함수이다.
// 삭제 과정은 실제, 사용자 문서를 active 컬렉션에서 inactive 컬렉션으로 옮깁니다.
func rmUser(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	u := User{}
	// 사용자 정보를 active 컬렉션에서 볼러온다.
	c := session.DB("user").C("acitive")
	err := c.Find(bson.M{"id": id}).One(&u)
	if err != nil {
		return err
	}
	err = c.Remove(bson.M{"id": id})
	if err != nil {
		log.Println(err)
		return err
	}
	// 불러온 사용자를 inactive에 추가한다.
	c = session.DB("user").C("inactive") // 비활성화계정으로 커서 변경.
	// inactive 컬렉션에 넣기전 사용자의 모든 권한을 해제한다.
	u.AccessLevel = UnknownAccessLevel
	err = c.Insert(u)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// setUser 함수는 사용자 정보를 수정하는 함수이다.
func setUser(session *mgo.Session, u User) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("active")
	num, err := c.Find(bson.M{"id": u.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	u.Updatetime = Now()
	err = c.Update(bson.M{"id": u.ID}, u)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 전체 유저 정보를 가지고오는 함수입니다.
func getUsers(session *mgo.Session) ([]User, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("active")
	var result []User
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
