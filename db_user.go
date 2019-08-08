// 이 코드는 사용자와 관련된 DBapi가 모여있는 파일입니다.

package main

import (
	"encoding/base64"
	"errors"
	"log"
	"sort"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	c := session.DB("user").C("users")

	num, err := c.Find(bson.M{"id": u.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		err = errors.New(u.ID + " ID를 가진 사용자가 이미 DB에 존재합니다.")
		log.Println(err)
		return err
	}
	u.Createtime = time.Now().Format(time.RFC3339)
	err = c.Insert(u)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// addToken 함수는 사용자정보로 token을 추가하는 함수이다.
func addToken(session *mgo.Session, u User) error {
	c := session.DB("user").C("token")
	num, err := c.Find(bson.M{"token": u.Token}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		err = errors.New(u.Token + " 키가 이미 DB에 존재합니다.")
		log.Println(err)
		return err
	}
	t := Token{
		Token:       u.Token,
		AccessLevel: u.AccessLevel,
		ID:          u.ID,
	}
	err = c.Insert(t)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// setToken 함수는 사용자 정보를 업데이트하는 함수이다.
func setToken(session *mgo.Session, t Token) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("token")
	err := c.Update(bson.M{"id": t.ID}, t)
	if err != nil {
		return err
	}
	return nil
}

// validToken 함수는 Token이 유효한지 체크한다.
func validToken(session *mgo.Session, token string) (Token, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("token")
	t := Token{}
	err := c.Find(bson.M{"token": token}).One(&t)
	if err != nil {
		return t, errors.New("authorization failed")
	}
	return t, nil
}

// getUser 함수는 사용자를 가지고오는 함수이다.
func getUser(session *mgo.Session, id string) (User, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	u := User{}
	err := c.Find(bson.M{"id": id}).One(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// getToken 함수는 사용자의 토큰을 가지고오는 함수이다.
func getToken(session *mgo.Session, id string) (Token, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("token")
	t := Token{}
	err := c.Find(bson.M{"id": id}).One(&t)
	if err != nil {
		return t, err
	}
	return t, nil
}

// rmUser 함수는 사용자를 삭제하는 함수이다.
func rmUser(session *mgo.Session, u User) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	err := c.Remove(bson.M{"id": u.ID})
	if err != nil {
		return err
	}
	return nil
}

// rmToken 함수는 token 키를 삭제하는 함수이다.
func rmToken(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("token")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// setUser 함수는 사용자 정보를 업데이트하는 함수이다.
func setUser(session *mgo.Session, u User) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	num, err := c.Find(bson.M{"id": u.ID}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	u.Updatetime = time.Now().Format(time.RFC3339)
	err = c.Update(bson.M{"id": u.ID}, u)
	if err != nil {
		return err
	}
	return nil
}

// initPassUser 함수는 사용자 정보를 수정하는 함수이다.
func initPassUser(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	num, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	q := bson.M{"id": id}
	encryptPass, err := Encrypt(*flagInitPass)
	if err != nil {
		log.Println(err)
		return err
	}
	change := bson.M{
		"$set": bson.M{
			"password":        encryptPass,
			"passwordattempt": 0,
			"updatetime":      time.Now().Format(time.RFC3339),
			"token":           base64.StdEncoding.EncodeToString([]byte(encryptPass)),
		},
	}
	err = c.Update(q, change)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// updatePasswordUser 함수는 사용자 패스워드를 수정하는 함수이다.
func updatePasswordUser(session *mgo.Session, id, pw, newPw string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	num, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	q := bson.M{"id": id}
	// 과거의 패스워드로 로그인가능했는지 체크한다.
	err = vaildUser(session, id, pw)
	if err != nil {
		return err
	}
	// 새로운 패스워드로 업데이트 한다.
	encryptPass, err := Encrypt(newPw)
	if err != nil {
		return err
	}
	change := bson.M{
		"$set": bson.M{
			"password":   encryptPass,
			"updatetime": time.Now().Format(time.RFC3339),
			"token":      base64.StdEncoding.EncodeToString([]byte(encryptPass)),
		},
	}
	err = c.Update(q, change)
	if err != nil {
		return err
	}
	return nil
}

// allUsers 함수는 DB에서 전체 사용자 정보를 가지고오는 함수입니다.
func allUsers(session *mgo.Session) ([]User, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	var result []User
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// searchUsers 함수는 검색을 입력받고 해당 검색어가 있는 사용자 정보를 가지고 옵니다.
func searchUsers(session *mgo.Session, words []string) ([]User, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	var searchwords []string
	// 사람 이름을 가지고 검색을 자주한다.
	for _, word := range words {
		if isASCII(word) {
			searchwords = append(searchwords, word)
			continue
		}
		if strings.HasPrefix(word, "tag:") {
			searchwords = append(searchwords, word)
			continue
		}
		r := []rune(word)
		if len(r) == 2 { // 이름이 2자리 일경우 "김웅"
			searchwords = append(searchwords, string(r[0])) // 성을 추가한다.
			searchwords = append(searchwords, string(r[1])) // 이름을 추가한다.
			continue
		} else if len(r) == 3 { // 이름일 확률이 높다.
			searchwords = append(searchwords, string(r[0]))  // 성을 추가한다.
			searchwords = append(searchwords, string(r[1:])) // 이름을 추가한다.
			continue
		} else if len(r) == 4 { // 이름이 4자리 일경우가 있다. 예)독고영재
			searchwords = append(searchwords, string(r[2:])) // 이름이 고영재" 또는 "영재" 일 수 있다. 이름을 위주로 검색시킨다.
			continue
		}
		searchwords = append(searchwords, word)
	}

	wordsQueries := []bson.M{}
	if *flagDebug {
		log.Println(searchwords)
	}
	for _, word := range searchwords {
		wordQueries := []bson.M{}
		if strings.HasPrefix(word, "tag:") {
			wordQueries = append(wordQueries, bson.M{"tags": &bson.RegEx{Pattern: strings.Trim(word, "tag:")}})
		} else {
			wordQueries = append(wordQueries, bson.M{"id": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"firstnamekor": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"lastnamekor": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"firstnameeng": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"lastnameeng": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"firstnamechn": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"lastnamechn": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"email": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"emailexternal": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"phone": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"hotline": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"location": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"tags": &bson.RegEx{Pattern: word}})
			wordQueries = append(wordQueries, bson.M{"lastip": &bson.RegEx{Pattern: word}})
		}
		wordsQueries = append(wordsQueries, bson.M{"$or": wordQueries})
	}
	q := bson.M{"$and": wordsQueries}
	var results []User
	err := c.Find(q).Sort("id").All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// vaildUser 함수는 사용자의 id, pw를 받아서 유효한 사용자인지 체크한다.
func vaildUser(session *mgo.Session, id, pw string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	q := bson.M{"id": id}
	num, err := c.Find(q).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	u := User{}
	err = c.Find(q).One(&u)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	if err != nil {
		return err
	}
	return nil
}

// addPasswordAttempt 함수는 사용자의 id를 받아서 패스워드 시도횟수를 추가한다.
func addPasswordAttempt(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	num, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$inc": bson.M{"passwordattempt": 1}})
	if err != nil {
		return err
	}
	return nil
}

// resetPasswordAttempt 함수는 사용자의 id를 받아서 패스워드 시도횟수를 초기화 한다.
func resetPasswordAttempt(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	num, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"passwordattempt": 0}})
	if err != nil {
		return err
	}
	return nil
}

// UserTags 함수는 전체 사용자에 등록된 Tags를 분석하여 태그리스트를 반환합니다.
func UserTags(session *mgo.Session) ([]string, error) {
	var tags []string
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	err := c.Find(bson.M{}).Distinct("tags", &tags)
	if err != nil {
		return nil, err
	}
	sort.Strings(tags)
	return tags, nil
}

// ReplaceTags 함수는 전체 사용자에 등록된 태그의 이름을 변경한다.
func ReplaceTags(session *mgo.Session, old, new string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("user").C("users")
	var users []User
	q := bson.M{"tags": &bson.RegEx{Pattern: old}}
	err := c.Find(q).All(&users)
	if err != nil {
		return err
	}
	for _, u := range users {
		var newTags []string
		for _, tag := range u.Tags {
			if tag == old {
				newTags = append(newTags, new)
				continue
			}
			newTags = append(newTags, tag)
		}
		u.Tags = newTags
		err = c.Update(bson.M{"id": u.ID}, u)
		if err != nil {
			return err
		}
	}
	// 각 유저를 체크하면서 태그이름을 변경한다.
	return nil
}
