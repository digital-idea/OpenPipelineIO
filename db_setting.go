package main

import (
	"errors"

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

// AddTaskSetting 함수는 tasksetting을 DB에 추가한다.
func AddTaskSetting(session *mgo.Session, t Tasksetting) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("tasksetting")
	num, err := c.Find(bson.M{"id": t.ID}).Count()
	if err != nil {
		return err
	}
	if num > 0 {
		return errors.New("이미 Tasksetting이 존재합니다")
	}
	err = c.Insert(t)
	if err != nil {
		return err
	}
	return nil
}

// RmTaskSetting 함수는 tasksetting을 DB에 추가한다.
func RmTaskSetting(session *mgo.Session, name, typ string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("tasksetting")
	err := c.Remove(bson.M{"name": name, "type": typ})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskSetting 함수는 Tasksetting 값을 바꾼다.
func SetTaskSetting(session *mgo.Session, t Tasksetting) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("tasksetting")
	err := c.Update(bson.M{"id": t.ID}, t)
	if err != nil {
		return err
	}
	return nil
}

// AllTaskSettings 함수는 모든 tasksetting값을 가지고 온다.
func AllTaskSettings(session *mgo.Session) ([]Tasksetting, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("tasksetting")
	results := []Tasksetting{}
	err := c.Find(bson.M{}).Sort("id").All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// getTaskSetting 함수는 id를 입력받아서 tasksetting값을 가지고 온다.
func getTaskSetting(session *mgo.Session, id string) (Tasksetting, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("tasksetting")
	result := Tasksetting{}
	err := c.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// getShotTaskSetting 함수는 type이 shot인 tasksetting값을 가지고 온다.
func getShotTaskSetting(session *mgo.Session) ([]Tasksetting, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("tasksetting")
	results := []Tasksetting{}
	err := c.Find(bson.M{"type": "shot"}).Sort("name").All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// getAssetTaskSetting 함수는 type이 asset인 tasksetting값을 가지고 온다.
func getAssetTaskSetting(session *mgo.Session) ([]Tasksetting, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("tasksetting")
	results := []Tasksetting{}
	err := c.Find(bson.M{"type": "asset"}).Sort("name").All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// TasksettingNames 함수는 Tasksetting 이름을 수집하여 반환한다.
func TasksettingNames(session *mgo.Session) ([]string, error) {
	var results []string
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setting").C("tasksetting")
	err := c.Find(bson.M{}).Sort("name").Distinct("name", &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
