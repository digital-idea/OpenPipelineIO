package main

import (
	"errors"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetelliteItems 함수는 프로젝트와 롤미디어 문자를 받아서 관련된 Setellite 현장자료를 반환한다.
func SetelliteItems(session *mgo.Session, project, rollmedia string) ([]Setellite, error) {
	if project == "" {
		return nil, errors.New("프로젝트를 설정해주세요")
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setellite").C(project)
	var results []Setellite
	err := c.Find(bson.M{"rollmedia": &bson.RegEx{Pattern: rollmedia, Options: "i"}}).All(&results)
	if err != nil {
		if err == mgo.ErrNotFound {
			return results, errors.New(project + "프로젝트에" + rollmedia + "이 존재하지 않습니다.")
		}
		log.Println(err)
		return results, err
	}
	return results, nil
}

func hasSetelliteItems(session *mgo.Session, project, rollmedia string) bool {
	if project == "" || rollmedia == "" {
		return false
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setellite").C(project)
	num, err := c.Find(bson.M{"rollmedia": &bson.RegEx{Pattern: rollmedia, Options: "i"}}).Count()
	if err != nil {
		return false
	}
	if num == 0 {
		return false
	}
	return true
}

// SetelliteSearch 함수는 프로젝트와 검색어를 입력받아서 검색어가 포함된 현장정보를 반환한다.
func SetelliteSearch(session *mgo.Session, project, word string) ([]Setellite, error) {
	if project == "" {
		return nil, errors.New("프로젝트를 설정해주세요")
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("setellite").C(project)
	var results []Setellite
	querys := []bson.M{}
	querys = append(querys, bson.M{"timestamp": &bson.RegEx{Pattern: word}})
	querys = append(querys, bson.M{"timestampstart": &bson.RegEx{Pattern: word}})
	querys = append(querys, bson.M{"episode": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"scenenumber": &bson.RegEx{Pattern: word}})
	querys = append(querys, bson.M{"storyboardnumber": &bson.RegEx{Pattern: word}})
	querys = append(querys, bson.M{"unit": &bson.RegEx{Pattern: word}})
	querys = append(querys, bson.M{"shootday": &bson.RegEx{Pattern: word}})
	querys = append(querys, bson.M{"vfxshot": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"scriptlocation": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"shotdescription": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"setlocation": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"elementtype": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"setmedia": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"wrangler": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"hdrifilename": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"cameraname": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"cameragrip": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"cameramodel": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"camerahead": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"lensmodel": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"lenstype": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"stereocameramodel": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"stereorigorientation": &bson.RegEx{Pattern: word}})
	querys = append(querys, bson.M{"stereolenstype": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"notes": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"rollmedia": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"shutterangle": &bson.RegEx{Pattern: word}})
	querys = append(querys, bson.M{"filter": &bson.RegEx{Pattern: word, Options: "i"}})
	querys = append(querys, bson.M{"stereoconvergence": &bson.RegEx{Pattern: word}})
	querys = append(querys, bson.M{"stereoia": &bson.RegEx{Pattern: word}})

	queries := []bson.M{bson.M{"$or": querys}}
	q := bson.M{"$and": queries}
	err := c.Find(q).Sort("timestamp").All(&results)
	if err != nil {
		if err == mgo.ErrNotFound {
			return results, errors.New(project + "프로젝트에 해당 검색어에 대한 아이템이 존재하지 않습니다.")
		}
		log.Println(err)
		return nil, err
	}
	return results, nil
}

// addSetellite함수는 Setellite 자료구조를 DB에 넣는다.
func addSetellite(session *mgo.Session, project string, item Setellite, overwrite bool) error {
	c := session.DB("setellite").C(project)
	num, err := c.Find(bson.M{"id": item.ID}).Count()
	if err != nil {
		return err
	}
	// 이미 존재하고 덮어쓰기 옵션이 꺼있을 때
	if num != 0 && !overwrite {
		return errors.New("데이터가 이미 존재합니다")
	}
	// 데이터가 존재하고 덮어쓰기 할 때
	if num != 0 && overwrite {
		target := bson.M{"id": item.ID}
		err = c.Update(target, item)
		if err != nil {
			return err
		}
		return nil
	}
	// DB에 존재하지 않아서 추가할 때
	err = c.Insert(item)
	if err != nil {
		return err
	}
	return nil
}
