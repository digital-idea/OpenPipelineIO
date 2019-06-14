package main

import (
	"errors"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Projectlist 함수는 프로젝트 리스트를 출력하는 함수입니다.
func Projectlist(session *mgo.Session) ([]string, error) {
	session.SetMode(mgo.Monotonic, true)
	Projectlist, err := session.DB("projectinfo").CollectionNames()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var results []string
	for _, project := range Projectlist {
		switch project {
		case "system.indexes":
			break //mongodb의 기본 컬렉션이다. 제외한다.
		default:
			results = append(results, project)
		}
	}
	return results, nil
}

// 프로젝트를 추가하는 함수입니다.
func addProject(session *mgo.Session, p Project) error {
	if p.ID == "" {
		return errors.New("빈 문자열입니다. 프로젝트를 생성할 수 없습니다")
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("projectinfo").C(p.ID)
	num, err := c.Find(bson.M{"id": p.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 0 {
		log.Println("같은 프로젝트가 존재하여 프로젝트를 생성할 수 없습니다.")
		return errors.New("같은 프로젝트가 존재해서 프로젝트를 생성할 수 없습니다")
	}
	err = c.Insert(p)
	if err != nil {
		log.Println("DB Insert Err : ", err)
		return err
	}
	// project DB에 "slug" 필드의 인덱스를 생성한다.
	c = session.DB("project").C(p.ID)
	err = c.EnsureIndexKey("slug")
	if err != nil {
		return err
	}
	return nil
}

func getProject(session *mgo.Session, id string) (Project, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("projectinfo").C(id)
	var p Project
	err := c.Find(bson.M{"id": id}).One(&p)
	if err != nil {
		p := Project{}
		p.ID = id
		if err == mgo.ErrNotFound {
			return p, errors.New(id + " 프로젝트가 존재하지 않습니다.")
		}
		return p, err
	}
	return p, nil
}

// 전체 프로젝트 정보를 가지고오는 함수입니다.
func getProjects(session *mgo.Session) ([]Project, error) {
	session.SetMode(mgo.Monotonic, true)
	projects, err := session.DB("projectinfo").CollectionNames()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var result []Project
	for _, project := range projects {
		if project == "system.indexes" { // 몽고DB가 자동으로 생성하는 컬렉션 입니다.
			continue
		}
		c := session.DB("projectinfo").C(project)
		var p Project
		err := c.Find(bson.M{"id": project}).One(&p)
		// 프로젝트코드의 문자열과 같은 id를 가진 문서를 DB에서 찾아 p값에 넣습니다.
		// 몽고DB에서 Find 메소드만 사용하면 documents 형태로 리턴되지 않기 때문에
		// One이라는 함수를 사용해서 바로 document로 값을 받을 수 있도록 처리했습니다.
		if err != nil {
			if err == mgo.ErrNotFound {
				log.Println(project + " 프로젝트가 존재하지 않습니다.")
				continue
			}
			log.Println(err)
			continue
		}
		result = append(result, p)
	}
	return result, nil
}

func rmProject(session *mgo.Session, project string) error {
	session.SetMode(mgo.Monotonic, true)
	//프로젝트 정보제거
	err := session.DB("projectinfo").C(project).DropCollection()
	if err != nil {
		log.Println(err)
		return err
	}
	//프로젝트 샷,에셋 제거
	err = session.DB("project").C(project).DropCollection()
	if err != nil {
		log.Println(err)
		return err
	}
	// 삭제 프로젝트 현장데이터가 존재하면 제거한다.
	collections, err := session.DB("setellite").CollectionNames()
	if err != nil {
		log.Println(err)
		return err
	}
	for _, c := range collections {
		if project != c {
			continue
		}
		err = session.DB("setellite").C(project).DropCollection()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

// setProject 함수는 프로젝트 정보를 수정하는 함수입니다.
func setProject(session *mgo.Session, p Project) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("projectinfo").C(p.ID)
	num, err := c.Find(bson.M{"id": p.ID}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 1 {
		return errors.New("해당 아이템이 존재하지 않습니다")
	}
	p.Updatetime = Now()
	err = c.Update(bson.M{"id": p.ID}, p)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// HasProject 함수는 프로젝트가 존재하는지 체크한다.
func HasProject(session *mgo.Session, project string) error {
	plist, err := Projectlist(session)
	if err != nil {
		return err
	}
	for _, p := range plist {
		if project == p {
			return nil
		}
	}
	return errors.New(project + " 프로젝트가 존재하지 않습니다.")
}
