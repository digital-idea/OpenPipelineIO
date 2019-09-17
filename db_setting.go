package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetAdminSetting 함수는 adminsetting을 DB에 저장한다.
func SetAdminSetting(session *mgo.Session, s Setting) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("admin")
	err := c.Update(bson.M{"id": "admin"}, s)
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

// GetAdminSetting 함수는 adminsetting을 DB에서 가지고 온다.
func GetAdminSetting(session *mgo.Session) (Setting, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("admin")
	s := Setting{}
	err := c.Find(bson.M{"id": "admin"}).One(&s)
	if err != nil {
		if err == mgo.ErrNotFound {
			s.ID = "admin"
			err = c.Insert(s)
			if err != nil {
				return s, err
			}
			return s, nil
		}
		return s, err
	}
	return s, nil
}
