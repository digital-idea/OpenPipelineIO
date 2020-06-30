package main

import (
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
			orQueries = append(orQueries, bson.M{"createtime": &bson.RegEx{Pattern: strings.TrimPrefix(word, "daily:")}})
		} else if strings.HasPrefix(word, "status:") {
			orQueries = append(orQueries, bson.M{"status": &bson.RegEx{Pattern: strings.TrimPrefix(word, "status:")}})
		} else if strings.HasPrefix(word, "project:") {
			orQueries = append(orQueries, bson.M{"project": &bson.RegEx{Pattern: strings.TrimPrefix(word, "project:")}})
		} else if strings.HasPrefix(word, "name:") {
			orQueries = append(orQueries, bson.M{"name": &bson.RegEx{Pattern: strings.TrimPrefix(word, "name:")}})
		} else {
			orQueries = append(orQueries, bson.M{"project": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"name": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"task": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"createtime": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"updatetime": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"author": &bson.RegEx{Pattern: word}})
			orQueries = append(orQueries, bson.M{"path": &bson.RegEx{Pattern: word}})
		}
		allQueries = append(allQueries, bson.M{"$or": orQueries})
	}
	q := bson.M{"$and": allQueries}
	err := c.Find(q).Sort("createtime").All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}
