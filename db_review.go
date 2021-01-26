package main

import (
	"errors"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func addReview(session *mgo.Session, r Review) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.Insert(r)
	if err != nil {
		return err
	}
	return nil
}

// GetWaitProcessStatusReview 함수는 processstatus 가 wait 값인 아이템을 하나 반환하고 상태를 processing으로 바꾼다.
func GetWaitProcessStatusReview() (Review, error) {
	var review Review
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		return review, err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err = c.Find(bson.M{"processstatus": "wait"}).One(&review)
	if err != nil {
		return review, err
	}
	// 참고: 아래 부분은 추가 mongo-driver로 바꾸면 한번에 처리할 수 있다.
	err = setReviewProcessStatus(session, review.ID.Hex(), "processing")
	if err != nil {
		return review, err
	}
	review.ProcessStatus = "processing"
	return review, nil
}

func getReview(session *mgo.Session, id string) (Review, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	r := Review{}
	err := c.FindId(bson.ObjectIdHex(id)).One(&r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func setReviewStatus(session *mgo.Session, id, status string) error {
	if !(status == "approve" || status == "wait" || status == "comment") {
		return errors.New("wait, approve, comment 상태만 사용가능합니다")
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}
	return nil
}

func setReviewProcessStatus(session *mgo.Session, id, status string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"processstatus": status}})
	if err != nil {
		return err
	}
	return nil
}

// setErrReview 함수는 id와 log를 입력받아서 에러상태 변경 및 로그를 기록한다.
func setErrReview(session *mgo.Session, id, log string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"processstatus": "error", "log": log}})
	if err != nil {
		return err
	}
	return nil
}

// searchReview 함수는 Review를 검색한다.
func searchReview(session *mgo.Session, searchword string) ([]Review, error) {
	var results []Review
	if searchword == "" {
		return results, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	allQueries := []bson.M{}
	for _, word := range strings.Split(searchword, " ") {
		orQueries := []bson.M{}
		if strings.HasPrefix(word, "daily:") {
			if strings.TrimPrefix(word, "daily:") == "" {
				return results, nil
			}
			orQueries = append(orQueries, bson.M{"createtime": &bson.RegEx{Pattern: strings.TrimPrefix(word, "daily:")}})
		} else if strings.HasPrefix(word, "status:") {
			orQueries = append(orQueries, bson.M{"status": &bson.RegEx{Pattern: strings.TrimPrefix(word, "status:")}})
		} else if strings.HasPrefix(word, "project:") {
			orQueries = append(orQueries, bson.M{"project": &bson.RegEx{Pattern: strings.TrimPrefix(word, "project:")}})
		} else if strings.HasPrefix(word, "name:") {
			orQueries = append(orQueries, bson.M{"name": &bson.RegEx{Pattern: strings.TrimPrefix(word, "name:")}})
		} else if strings.HasPrefix(word, "task:") {
			orQueries = append(orQueries, bson.M{"task": &bson.RegEx{Pattern: strings.TrimPrefix(word, "task:")}})
		} else {
			orQueries = append(orQueries, bson.M{"project": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"name": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"task": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"createtime": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"updatetime": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"author": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"path": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"description": &bson.RegEx{Pattern: word}})
		}
		allQueries = append(allQueries, bson.M{"$or": orQueries})
	}
	q := bson.M{"$and": allQueries}
	err := c.Find(q).Sort("-createtime").All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// addReviewComment 함수는 Review에 Comment를 추가한다.
func addReviewComment(session *mgo.Session, id string, cmt Comment) error {
	if cmt.Text == "" {
		return errors.New("comment가 빈 문자열입니다")
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$push": bson.M{"comments": cmt}})
	if err != nil {
		return err
	}
	return nil
}

// EditReviewComment 함수는 review에 comment를 수정합니다.
func EditReviewComment(session *mgo.Session, id, date, text string) error {
	session.SetMode(mgo.Monotonic, true)
	reviewItem, err := getReview(session, id)
	if err != nil {
		return err
	}
	var newComments []Comment
	for _, comment := range reviewItem.Comments {
		if comment.Date == date {
			comment.Text = text
		}
		newComments = append(newComments, comment)
	}
	reviewItem.Comments = newComments
	err = setReviewItem(session, reviewItem)
	if err != nil {
		return err
	}
	return nil
}

// RmReviewComment 함수는 review에 comment를 삭제합니다.
func RmReviewComment(session *mgo.Session, id, date string) error {
	session.SetMode(mgo.Monotonic, true)
	reviewItem, err := getReview(session, id)
	if err != nil {
		return err
	}
	var newComments []Comment
	for _, comment := range reviewItem.Comments {
		if comment.Date == date {
			continue
		}
		newComments = append(newComments, comment)
	}
	reviewItem.Comments = newComments
	err = setReviewItem(session, reviewItem)
	if err != nil {
		return err
	}
	return nil
}

// setReviewItem은 Review 자료구조를 새로운 Review로 설정한다.
func setReviewItem(session *mgo.Session, r Review) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(r.ID, r)
	if err != nil {
		return err
	}
	return nil
}

// RmReview 함수는 Review를 DB에서 삭제한다.
func RmReview(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}
	return nil
}

// RmProjectReview 함수는 해당 프로젝트의 Review 데이터를 DB에서 삭제한다.
func RmProjectReview(session *mgo.Session, project string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	_, err := c.RemoveAll(bson.M{"project": &bson.RegEx{Pattern: project}})
	if err != nil {
		return err
	}
	return nil
}

// SetReviewProject 함수는 Review에 Project를 설정한다.
func SetReviewProject(session *mgo.Session, id string, project string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"project": project}})
	if err != nil {
		return err
	}
	return nil
}

// SetReviewTask 함수는 Review에 Task를 설정한다.
func SetReviewTask(session *mgo.Session, id string, task string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"task": task}})
	if err != nil {
		return err
	}
	return nil
}

// SetReviewName 함수는 Review에 Name을 설정한다.
func SetReviewName(session *mgo.Session, id string, name string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		return err
	}
	return nil
}

// SetReviewPath 함수는 Review에 Path를 설정한다.
func SetReviewPath(session *mgo.Session, id string, path string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"path": path}})
	if err != nil {
		return err
	}
	return nil
}

// SetReviewMainVersion 함수는 Review에 MainVersion을 설정한다.
func SetReviewMainVersion(session *mgo.Session, id string, mainversion int) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"mainversion": mainversion}})
	if err != nil {
		return err
	}
	return nil
}

// SetReviewSubVersion 함수는 Review에 SubVersion을 설정한다.
func SetReviewSubVersion(session *mgo.Session, id string, subversion int) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"subversion": subversion}})
	if err != nil {
		return err
	}
	return nil
}

// SetReviewFps 함수는 Review에 Fps를 설정한다.
func SetReviewFps(session *mgo.Session, id string, fps float64) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"fps": fps}})
	if err != nil {
		return err
	}
	return nil
}

// SetReviewDescription 함수는 Review에 Description을 설정한다.
func SetReviewDescription(session *mgo.Session, id string, description string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"description": description}})
	if err != nil {
		return err
	}
	return nil
}

// SetReviewCameraInfo 함수는 Review에 CameraInfo를 설정한다.
func SetReviewCameraInfo(session *mgo.Session, id string, camerainfo string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("csi").C("review")
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"camerainfo": camerainfo}})
	if err != nil {
		return err
	}
	return nil
}
