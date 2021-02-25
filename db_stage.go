package main

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetStage 함수는 Stage를 DB에 저장한다.
func SetStage(session *mgo.Session, s Stage) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("stage")
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

// GetStage 함수는 Stage를 DB에서 가지고 온다.
func GetStage(session *mgo.Session, id string) (Stage, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("stage")
	s := Stage{}
	err := c.Find(bson.M{"id": id}).One(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// GetInitStageID 함수는 초기에 설정해야하는 Stage를 DB에서 가지고 온다.
func GetInitStageID(session *mgo.Session) (string, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("stage")
	s := Stage{}
	n, err := c.Find(bson.M{"initstage": true}).Count()
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", errors.New("초기 상태값 설정이 필요합니다")
	}
	if n != 1 {
		return "", errors.New("초기 상태 설정값이 1개가 아닙니다")
	}
	err = c.Find(bson.M{"initstage": true}).One(&s)
	if err != nil {
		return "", err
	}
	return s.ID, nil
}

// AddStage 함수는 Stage를 DB에 추가한다.
func AddStage(session *mgo.Session, s Stage) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("stage")
	n, err := c.Find(bson.M{"id": s.ID}).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New(s.ID + " Stage가 이미 존재합니다")
	}
	err = c.Insert(s)
	if err != nil {
		return err
	}
	return nil
}

// RmStage 함수는 Stage를 DB에서 삭제한다.
func RmStage(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("stage")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// AllStages 함수는 모든 Stage값을 DB에서 가지고 온다.
func AllStages(session *mgo.Session) ([]Stage, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("stage")
	results := []Stage{}
	err := c.Find(bson.M{}).Sort("-order").All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
