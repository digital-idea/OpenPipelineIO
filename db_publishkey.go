package main

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetPublishKey 함수는 PublishKey를 DB에 저장한다.
func SetPublishKey(session *mgo.Session, key PublishKey) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("publishkey")
	err := c.Update(bson.M{"id": key.ID}, key)
	if err != nil {
		if err == mgo.ErrNotFound {
			err = c.Insert(key)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

// GetPublishKey 함수는 PublishKey를 DB에서 가지고 온다.
func GetPublishKey(session *mgo.Session, id string) (PublishKey, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("publishkey")
	key := PublishKey{}
	err := c.Find(bson.M{"id": id}).One(&key)
	if err != nil {
		return key, err
	}
	return key, nil
}

// AddPublishKey 함수는 PublishKey를 DB에 추가한다.
func AddPublishKey(session *mgo.Session, key PublishKey) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("publishkey")
	n, err := c.Find(bson.M{"id": key.ID}).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New(key.ID + " PublishKey가 이미 존재합니다")
	}
	err = c.Insert(key)
	if err != nil {
		return err
	}
	return nil
}

// RmPublishKey 함수는 PublishKey를 DB에서 삭제한다.
func RmPublishKey(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("publishkey")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// AllPublishKeys 함수는 모든 PublishKey 값을 DB에서 가지고 온다.
func AllPublishKeys(session *mgo.Session) ([]PublishKey, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("publishkey")
	keys := []PublishKey{}
	err := c.Find(bson.M{}).Sort("id").All(&keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
