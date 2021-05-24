package main

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetImageSize 함수는 해당 샷의 이미지 사이즈를 설정한다. // legacy
// key 설정값 : platesize, undistortionsize, rendersize
func SetImageSize(session *mgo.Session, project, name, key, size string) (string, error) {
	if !(key == "platesize" || key == "dsize" || key == "undistortionsize" || key == "rendersize") {
		return "", errors.New("잘못된 key값입니다")
	}
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return "", err
	}
	id := name + "_" + typ
	c := session.DB("project").C(project)
	if key == "dsize" || key == "undistortionsize" {
		err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"dsize": size, "undistortionsize": size, "updatetime": time.Now().Format(time.RFC3339)}})
		if err != nil {
			return id, err
		}
	} else {
		err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{key: size, "updatetime": time.Now().Format(time.RFC3339)}})
		if err != nil {
			return id, err
		}
	}
	return id, nil
}
