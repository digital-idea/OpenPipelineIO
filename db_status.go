package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetStatus 함수는 Status를 DB에 저장한다.
func SetStatus(session *mgo.Session, s StatusV2) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("status")
	err := c.Update(bson.M{"id": s.ID}, s)
	if err != nil {
		if err == mgo.ErrNotFound {
			err = c.Insert(s)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

// GetStatus 함수는 Status를 DB에서 가지고 온다.
func GetStatus(session *mgo.Session, id string) (StatusV2, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("status")
	s := StatusV2{}
	err := c.Find(bson.M{"id": id}).One(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// AddStatus 함수는 tasksetting을 DB에 추가한다.
func AddStatus(session *mgo.Session, s StatusV2) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("status")
	err := c.Insert(s)
	if err != nil {
		return err
	}
	return nil
}

// RmStatus 함수는 Status를 DB에서 삭제한다.
func RmStatus(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("status")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// AllStatus 함수는 모든 Status값을 DB에서 가지고 온다.
func AllStatus(session *mgo.Session) ([]StatusV2, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("status")
	results := []StatusV2{}
	err := c.Find(bson.M{}).Sort("order").All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
