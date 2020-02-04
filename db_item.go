package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/digital-idea/ditime"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func addItem(session *mgo.Session, project string, i Item) error {
	session.SetMode(mgo.Monotonic, true)
	// 프로젝트가 존재하는지 체크합니다.
	c := session.DB("projectinfo").C(project)
	num, err := c.Find(bson.M{"id": project}).Count()
	if err != nil {
		log.Println(err)
		return err
	}
	if num != 1 {
		log.Println("프로젝트가 존제하지 않습니다.")
		return errors.New("프로젝트가 존재하지 않습니다")
	}
	//문서의 중복이 있는지 체크합니다.
	c = session.DB("project").C(project)
	num, err = c.Find(bson.M{"name": i.Name, "type": i.Type}).Count()
	if err != nil {
		return err
	}
	if num != 0 {
		return fmt.Errorf("%s 프로젝트에 이미 %s_%s 샷은 존재합니다", project, i.Name, i.Type)
	}
	err = c.Insert(i)
	if err != nil {
		return err
	}
	return nil
}

func setItem(session *mgo.Session, project string, i Item) error {
	session.SetMode(mgo.Monotonic, true)
	i.Updatetime = time.Now().Format(time.RFC3339)
	i.updateStatus() // Task상태 업데이트
	i.setRnumTag()   // 롤넘버에 따른 테그 셋팅
	c := session.DB("project").C(project)
	err := c.Update(bson.M{"id": i.ID}, i)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func getItem(session *mgo.Session, project, id string) (Item, error) {
	if project == "" || id == "" {
		return Item{}, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	var result Item
	err := c.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		log.Println(err)
		return Item{}, err
	}
	return result, nil
}

// Shot 함수는 프로젝트명, 샷이름을 이용해서 샷정보를 반환한다.
func Shot(session *mgo.Session, project string, name string) (Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	var result Item
	// org, left 갯수를 구한다.
	orgnum, err := c.Find(bson.M{"id": name + "_org"}).Count()
	if err != nil {
		return Item{}, err
	}
	leftnum, err := c.Find(bson.M{"id": name + "_left"}).Count()
	if err != nil {
		return Item{}, err
	}
	q := bson.M{"id": name + "_org"}
	if leftnum == 1 && orgnum == 0 {
		q = bson.M{"id": name + "_left"}
	}
	err = c.Find(q).One(&result)
	if err != nil {
		log.Println(err)
		return Item{}, err
	}
	return result, nil
}

// SearchName 함수는 입력된 문자열이 'name'키 값에 포함되어 있다면 해당 아이템을 반환한다.
func SearchName(session *mgo.Session, project string, name string) ([]Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	results := []Item{}
	err := c.Find(bson.M{"name": &bson.RegEx{Pattern: name, Options: "i"}}).Sort("name").All(&results)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return results, nil
}

// Seqs 함수는 프로젝트 이름을 받아서 seq 리스트를 반환한다.
func Seqs(session *mgo.Session, project string) ([]string, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	var results []Item
	err := c.Find(bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}}).Select(bson.M{"seq": 1}).All(&results)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	keys := make(map[string]bool)
	for _, result := range results {
		seq := result.Seq
		keys[seq] = true
	}
	seqs := []string{}
	for k := range keys {
		seqs = append(seqs, k)
	}
	sort.Strings(seqs)
	return seqs, nil
}

// Shots 함수는 프로젝트 이름과 입력된 시퀀스가 'name'키 값에 포함되어 있다면 shots 리스트를 반환한다.
func Shots(session *mgo.Session, project string, seq string) ([]string, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	var results []Item
	query := bson.M{"name": &bson.RegEx{Pattern: seq, Options: "i"}, "type": bson.M{"$in": []string{"org", "left"}}}
	err := c.Find(query).Select(bson.M{"name": 1}).All(&results)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	keys := make(map[string]bool)
	for _, result := range results {
		shot := strings.Split(result.Name, "_")[1]
		keys[shot] = true
	}
	shots := []string{}
	for k := range keys {
		shots = append(shots, k)
	}
	sort.Strings(shots)
	return shots, nil
}

func rmItem(session *mgo.Session, project, name, usertyp string) error {
	session.SetMode(mgo.Monotonic, true)
	var typ string
	if usertyp == "" {
		t, err := Type(session, project, name)
		if err != nil {
			return err
		}
		typ = t
	} else {
		typ = usertyp
	}
	c := session.DB("project").C(project)
	num, err := c.Find(bson.M{"name": name, "type": typ}).Count()
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("삭제할 아이템이 없습니다")
	}
	err = c.Remove(bson.M{"name": name, "type": typ})
	if err != nil {
		return err
	}
	return nil
}

func rmItemID(session *mgo.Session, project, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

func rmItemAndType(session *mgo.Session, project, name, typ string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	num, err := c.Find(bson.M{"name": name, "type": typ}).Count()
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("삭제할 아이템이 없습니다")
	}
	err = c.Remove(bson.M{"name": name, "type": typ})
	if err != nil {
		return err
	}
	return nil
}

// Distinct 함수는 프로젝트, dict key를 받아서 key에 사용되는 모든 문자열을 반환한다. 예) 태그
func Distinct(session *mgo.Session, project string, key string) ([]string, error) {
	var result []string
	if project == "" || key == "" {
		return result, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	err := c.Find(bson.M{}).Distinct(key, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sort.Strings(result)
	return result, nil
}

// DistinctDdline 함수는 프로젝트, dict key를 받아서 key에 사용되는 모든 마감일을 반환한다. 예) 태그
func DistinctDdline(session *mgo.Session, project string, key string) ([]string, error) {
	var result []string
	if project == "" || key == "" {
		return result, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	err := c.Find(bson.M{}).Distinct(key, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//result로 datelist를 만든다.
	sort.Strings(result)
	if *flagDebug {
		fmt.Println("DB에서 가지고온 마감일 리스트")
		fmt.Println(result)
		fmt.Println()
	}
	var before string
	var datelist []string
	for _, r := range result {
		if r != "" {
			date := ToNormalTime(r)
			if date == before {
				break
			} else {
				datelist = append(datelist, date)
			}
			before = date
		}
	}
	sort.Strings(datelist) //기존 CSI2의 4자리 수를 위하여 정렬한다. 추후 이 줄은 사라진다.
	if *flagDebug {
		fmt.Println("마감일을 Tag형태로 바꾼 리스트")
		fmt.Println(datelist)
		fmt.Println()
	}
	return datelist, nil
}

// SearchAllShot 함수는 모든 Shot 데이터를 반환한다. Excel을 뽑을 때 사용한다.
func SearchAllShot(session *mgo.Session, project, sortkey string) ([]Item, error) {
	results := []Item{}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	queries := []bson.M{}
	queries = append(queries, bson.M{"type": "org"})
	queries = append(queries, bson.M{"type": "left"})
	q := bson.M{"$or": queries}
	err := c.Find(q).Sort(sortkey).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SearchAllAsset 함수는 모든 Asset 데이터를 반환한다. Excel을 뽑을 때 사용한다.
func SearchAllAsset(session *mgo.Session, project, sortkey string) ([]Item, error) {
	results := []Item{}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	err := c.Find(bson.M{"type": "asset"}).Sort(sortkey).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SearchAll 함수는 모든 Shot, Asset 데이터를 반환한다. Excel을 뽑을 때 사용한다.
func SearchAll(session *mgo.Session, project, sortkey string) ([]Item, error) {
	results := []Item{}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	queries := []bson.M{}
	queries = append(queries, bson.M{"type": "org"})
	queries = append(queries, bson.M{"type": "left"})
	queries = append(queries, bson.M{"type": "asset"})
	q := bson.M{"$or": queries}
	err := c.Find(q).Sort(sortkey).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SearchTag 함수는 태그를 검색할때 사용한다.
// SearchTags라고 이름을 붙히지 않은 이유는 CSI 자료구조의 필드명이 Tag이기 때문이다.
// 미래에 Tag필드를  Tags 필드로 바꾼후 이 함수의 이름을 SearchTags로 바꿀 예정이다.
func SearchTag(session *mgo.Session, op SearchOption) ([]Item, error) {
	return SearchKey(session, op, "tag")
}

// SearchAssettags 함수는 검색옵션으로 에셋태그를 검색할때 사용한다.
func SearchAssettags(session *mgo.Session, op SearchOption) ([]Item, error) {
	return SearchKey(session, op, "assettags")
}

// SearchAssetTree 함수는 에셋이름을 하나 받아서 관련된 모든 에셋을 구하는 함수이다.
func SearchAssetTree(session *mgo.Session, op SearchOption) ([]Item, error) {
	if op.Searchword == "" { // 검색어가 없다면 종료한다.
		return []Item{}, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(op.Project)
	var allAssets []Item
	err := c.Find(bson.M{"type": "asset"}).Sort(op.Sortkey).All(&allAssets)
	if err != nil {
		if err == mgo.ErrNotFound {
			return []Item{}, nil
		}
		log.Println("DB Find Err : ", err)
		return nil, err
	}

	var result []Item
	var waittags []string
	var donetags []string
	waittags = append(waittags, op.Searchword) // 최초 검색단어를 대기리스트에 넣는다.
	for len(waittags) > 0 {
		// fmt.Println("-------------------------------------------------------------") // 디버그 구분선
		hasDone := false
		for _, donetag := range donetags {
			if donetag == waittags[0] {
				// fmt.Printf("%s는 이미 과거에 처리되었습니다.\n", waittags[0]) // 디버그를 위해서 놔둔다.
				hasDone = true
			}
		}
		for _, item := range allAssets {
			if item.Name == waittags[0] {
				// fmt.Printf("%s은 자기자신입니다.\n", item.Name) // 디버그를 위해서 놔둔다.
				donetags = append(donetags, item.Name)
				if !hasDone {
					result = append(result, item)
				}
			} else {
				for _, tag := range item.Assettags {
					if tag == waittags[0] {
						// fmt.Printf("%s는 %s와 연결되어 있습니다.\n", item.Name, tag) // 디버그를 위해서 놔둔다.
						waittags = append(waittags, item.Name)
					}
				}
			}
		}
		// fmt.Println("디버그:", waittags) // 디버그를 위해서 놔둔다.
		// fmt.Println("디버그:", donetags) // 디버그를 위해서 놔둔다.
		waittags = waittags[1:] // 처리한 값은 뺀다.
	}
	return result, nil
}

// SearchKey 함수는 Item.{key} 필드의 값과 검색어가 정확하게 일치하는 항목들만 검색한다.
func SearchKey(session *mgo.Session, op SearchOption, key string) ([]Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(op.Project)
	query := []bson.M{}
	// 아래 단어는 CSI에 버튼으로 되어있는 태그단어를 클릭시 작동되는 예약어이다.
	switch op.Searchword {
	case "2d", "2D":
		query = append(query, bson.M{"shottype": &bson.RegEx{Pattern: "2d", Options: "i"}})
	case "3d", "3D":
		query = append(query, bson.M{"shottype": &bson.RegEx{Pattern: "3d", Options: "i"}})
	case "asset":
		query = append(query, bson.M{"type": &bson.RegEx{Pattern: "asset", Options: "i"}})
	case "assign":
		query = append(query, bson.M{"status": ASSIGN})
	case "ready":
		query = append(query, bson.M{"status": READY})
	case "wip":
		query = append(query, bson.M{"status": WIP})
	case "confirm":
		query = append(query, bson.M{"status": CONFIRM})
	case "done":
		query = append(query, bson.M{"status": DONE})
	case "omit":
		query = append(query, bson.M{"status": OMIT})
	case "hold":
		query = append(query, bson.M{"status": HOLD})
	case "out":
		query = append(query, bson.M{"status": OUT})
	case "none":
		query = append(query, bson.M{"status": NONE})
	default:
		query = append(query, bson.M{key: op.Searchword})
	}

	var results []Item
	if !op.Assign && !op.Ready && !op.Wip && !op.Confirm && !op.Done && !op.Omit && !op.Hold && !op.Out && !op.None {
		// 체크박스가 아무것도 켜있지 않다면 바로 빈 값을 리턴한다.
		return results, nil
	}
	status := []bson.M{}
	if op.Assign {
		status = append(status, bson.M{"status": ASSIGN})
	}
	if op.Ready {
		status = append(status, bson.M{"status": READY})
	}
	if op.Wip {
		status = append(status, bson.M{"status": WIP})
	}
	if op.Confirm {
		status = append(status, bson.M{"status": CONFIRM})
	}
	if op.Done {
		status = append(status, bson.M{"status": DONE})
	}
	if op.Omit {
		status = append(status, bson.M{"status": OMIT})
	}
	if op.Hold {
		status = append(status, bson.M{"status": HOLD})
	}
	if op.Out {
		status = append(status, bson.M{"status": OUT})
	}
	if op.None {
		status = append(status, bson.M{"status": NONE})
	}

	q := bson.M{"$and": []bson.M{
		bson.M{"$or": query},
		bson.M{"$or": status},
	}}
	err := c.Find(q).Sort(op.Sortkey).All(&results)
	if err != nil {
		log.Println("DB Find Err : ", err)
		return nil, err
	}
	return results, nil
}

// SearchDdline 함수는 검색옵션, 파트정보(2d,3d)를 받아서 쿼리한다.
func SearchDdline(session *mgo.Session, op SearchOption, part string) ([]Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(op.Project)
	query := []bson.M{}
	switch part {
	case "2d":
		query = append(query, bson.M{"ddline2d": &bson.RegEx{Pattern: "....-" + op.Searchword[0:2] + "-" + op.Searchword[2:4]}})
	case "3d":
		query = append(query, bson.M{"ddline3d": &bson.RegEx{Pattern: "....-" + op.Searchword[0:2] + "-" + op.Searchword[2:4]}})
	default:
		query = append(query, bson.M{})
	}

	status := []bson.M{}
	if op.Assign {
		status = append(status, bson.M{"status": ASSIGN})
	}
	if op.Ready {
		status = append(status, bson.M{"status": READY})
	}
	if op.Wip {
		status = append(status, bson.M{"status": WIP})
	}
	if op.Confirm {
		status = append(status, bson.M{"status": CONFIRM})
	}
	if op.Done {
		status = append(status, bson.M{"status": DONE})
	}
	if op.Omit {
		status = append(status, bson.M{"status": OMIT})
	}
	if op.Hold {
		status = append(status, bson.M{"status": HOLD})
	}
	if op.Out {
		status = append(status, bson.M{"status": OUT})
	}
	if op.None {
		status = append(status, bson.M{"status": NONE})
	}

	q := bson.M{"$and": []bson.M{
		bson.M{"$or": query},
		bson.M{"$or": status},
	}}
	var results []Item
	err := c.Find(q).Sort(op.Sortkey).All(&results)
	if err != nil {
		log.Println("DB Find Err : ", err)
		return nil, err
	}
	return results, nil
}

// Searchnum 함수는 검색된 결과에 대한 상태별 갯수를 검색한다.
func Searchnum(project string, items []Item) (Infobarnum, error) {
	var results Infobarnum
	results.Search = len(items)
	for _, item := range items {
		if item.Shottype == "2D" || item.Shottype == "2d" {
			results.Shot2d++
		}
		if item.Shottype == "3D" || item.Shottype == "3d" {
			results.Shot3d++
		}
		if item.Type == "asset" {
			results.Assets++
		}
		if item.Type == "org" || item.Type == "left" {
			results.Shot++
		}
		switch item.Status {
		case ASSIGN:
			results.Assign++
		case READY:
			results.Ready++
		case WIP:
			results.Wip++
		case CONFIRM:
			results.Confirm++
		case DONE:
			results.Done++
		case OMIT:
			results.Omit++
		case HOLD:
			results.Hold++
		case OUT:
			results.Out++
		case NONE:
			results.None++
		}
	}
	return results, nil
}

// Totalnum 함수는 프로젝트의 전체샷에 대한 상태 갯수를 검색한다.
func Totalnum(session *mgo.Session, project string) (Infobarnum, error) {
	if project == "" {
		return Infobarnum{}, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)

	var results Infobarnum
	//진행률 출력.
	assign := bson.M{}
	assign = bson.M{"$and": []bson.M{
		bson.M{"status": ASSIGN},
		bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}},
	}}
	var assignnum int
	assignnum, err := c.Find(assign).Count()
	if err != nil {
		log.Println("DB Find Err : ", err)
		return Infobarnum{}, err
	}

	ready := bson.M{}
	ready = bson.M{"$and": []bson.M{
		bson.M{"status": READY},
		bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}},
	}}
	var readynum int
	readynum, err = c.Find(ready).Count()
	if err != nil {
		log.Println("DB Find Err : ", err)
		return Infobarnum{}, err
	}

	wip := bson.M{"$and": []bson.M{
		bson.M{"status": WIP},
		bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}},
	}}
	var wipnum int
	wipnum, err = c.Find(wip).Count()
	if err != nil {
		log.Println("DB Find Err : ", err)
		return Infobarnum{}, err
	}

	confirm := bson.M{}
	confirm = bson.M{"$and": []bson.M{
		bson.M{"status": CONFIRM},
		bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}},
	}}
	var confirmnum int
	confirmnum, err = c.Find(confirm).Count()
	if err != nil {
		log.Println("DB Find Err : ", err)
		return Infobarnum{}, err
	}

	done := bson.M{"$and": []bson.M{
		bson.M{"status": DONE},
		bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}},
	}}
	var donenum int
	donenum, err = c.Find(done).Count()
	if err != nil {
		log.Println("DB Find Err : ", err)
		return Infobarnum{}, err
	}

	omit := bson.M{"$and": []bson.M{
		bson.M{"status": OMIT},
		bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}},
	}}
	var omitnum int
	omitnum, err = c.Find(omit).Count()
	if err != nil {
		log.Println("DB Find Err : ", err)
		return Infobarnum{}, err
	}

	hold := bson.M{"$and": []bson.M{
		bson.M{"status": HOLD},
		bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}},
	}}
	var holdnum int
	holdnum, err = c.Find(hold).Count()
	if err != nil {
		log.Println("DB Find Err : ", err)
		return Infobarnum{}, err
	}

	out := bson.M{"$and": []bson.M{
		bson.M{"status": OUT},
		bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}},
	}}
	var outnum int
	outnum, err = c.Find(out).Count()
	if err != nil {
		log.Println("DB Find Err : ", err)
		return Infobarnum{}, err
	}

	none := bson.M{"$and": []bson.M{
		bson.M{"status": NONE},
		bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}},
	}}
	var nonenum int
	nonenum, err = c.Find(none).Count()
	if err != nil {
		log.Println("DB Find Err : ", err)
		return Infobarnum{}, err
	}

	var totalnum int
	totalnum, err = c.Find(bson.M{"$or": []bson.M{bson.M{"type": "org"}, bson.M{"type": "left"}}}).Count()
	if err != nil {
		log.Println("DB Find Err : ", err)
		return Infobarnum{}, err
	}

	results.Assign = assignnum
	results.Ready = readynum
	results.Wip = wipnum
	results.Confirm = confirmnum
	results.Done = donenum
	results.Omit = omitnum
	results.Hold = holdnum
	results.Out = outnum
	results.None = nonenum
	results.Total = totalnum

	return results, nil
}

// setTaskMov함수는 해당 샷에 mov를 설정하는 함수이다.
func setTaskMov(session *mgo.Session, project, name, task, mov string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".mov": mov, task + ".mdate": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// setTaskDue함수는 해당 샷에 mov를 설정하는 함수이다.
func setTaskDue(session *mgo.Session, project, name, task string, due int) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".due": due, task + ".mdate": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// setTaskLevel함수는 해당 샷에 level를 설정하는 함수이다.
func setTaskLevel(session *mgo.Session, project, name, task, level string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)

	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	l, err := strconv.Atoi(level)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".tasklevel": TaskLevel(l)}})
	if err != nil {
		return err
	}
	return nil
}

// Type 함수는 이름을 이용해서 Type값을 반환한다.
// 일반샷은 org를 반환한다.
// 입체샷인경우 left를 반환한다.
// 에셋은 asset을 반환한다.
func Type(session *mgo.Session, project, name string) (string, error) {
	c := session.DB("project").C(project)
	var items []Item
	err := c.Find(bson.M{"$or": []bson.M{bson.M{"name": name, "type": "org"}, bson.M{"name": name, "type": "left"}, bson.M{"name": name, "type": "asset"}}}).All(&items)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", errors.New(name + "에 해당하는 org,left,asset 타입을 DB에서 찾을 수 없습니다.")
	}
	if len(items) != 1 {
		return "", errors.New(name + "값이 DB에서 고유하지 않습니다.")
	}
	return items[0].Type, nil
}

// SetImageSize 함수는 해당 샷의 이미지 사이즈를 설정한다.
// key 설정값 : platesize, distortionsize, rendersize
func SetImageSize(session *mgo.Session, project, name, key, size string) (string, error) {
	if !(key == "platesize" || key == "dsize" || key == "rendersize") {
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
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{key: size, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, err
	}
	return id, nil
}

// SetTimecode 함수는 item에 Timecode를 설정한다.
// ScanTimecodeIn,ScanTimecodeOut,JustTimecodeIn,JustTimecoeOut 문자를 key로 사용할 수 있다.
func SetTimecode(session *mgo.Session, project, name, key, timecode string) error {
	key = strings.ToLower(key)
	if !(key == "scantimecodein" ||
		key == "scantimecodeout" ||
		key == "justtimecodein" ||
		key == "justtimecodeout") {
		return errors.New("scantimecodein, scantimecodeout, justtimecodein, justtimecodeout 키값만 사용가능")
	}
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{key: timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	// 우리회사는 현재 timecode와 keycode를 혼용해서 사용중이다.
	// 원래는 Timecode가 맞지만 현재 DB가 keycode로 되어있어 아직은 아래줄이 필요하다.
	key = strings.Replace(key, "timecode", "keycode", -1)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{key: timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetUseType 함수는 item에 UseType string을 설정한다.
func SetUseType(session *mgo.Session, project, name, usetype string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"usetype": usetype, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetFrame 함수는 item에 프레임을 설정한다.
// ScanIn,ScanOut,ScanFrame,PlateIn,PlateOut,JustIn,JustOut,HandleIn,HandleOut 문자를 key로 사용할 수 있다.
func SetFrame(session *mgo.Session, project, name, key string, frame int) error {
	if frame == -1 {
		return nil
	}
	if !(key == "scanin" ||
		key == "scanout" ||
		key == "scanframe" ||
		key == "platein" ||
		key == "plateout" ||
		key == "justin" ||
		key == "justout" ||
		key == "handlein" ||
		key == "handleout") {
		return errors.New("scanin, scanout, scanframe, platein, plateout, justin, justout, handlein, handleout 키값만 사용가능합니다")
	}
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{key: frame, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetCameraPubPath 함수는 해당 카메라 퍼블리쉬 경로를 설정한다.
func SetCameraPubPath(session *mgo.Session, project, name, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"productioncam.pubpath": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetCameraPubTask 함수는 해당 카메라 퍼블리쉬 팀을 설정한다.
func SetCameraPubTask(session *mgo.Session, project, name, task string) error {
	if !(task == "" || task == "mm" || task == "layout" || task == "ani") {
		return errors.New("none(빈문자열), mm, layout, ani 팀만 카메라 publish가 가능합니다")
	}
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"productioncam.pubtask": task, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetCameraProjection 함수는 샷에 Projection 카메라 사용여부를 체크한다.
func SetCameraProjection(session *mgo.Session, project, name string, isProjection bool) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"productioncam.projection": isProjection, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetObjectID 함수는 Item에 Object In, Out 값을 설정한다.
func SetObjectID(session *mgo.Session, project, name string, in, out int) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ != "asset" {
		return errors.New("asset 타입이 아닙니다")
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"objectidin": in, "objectidout": out, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetThummov 함수는 item에 Thummov값을 셋팅한다.
func SetThummov(session *mgo.Session, project, name, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"thummov": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetBeforemov 함수는 item에 Before mov값을 셋팅한다.
func SetBeforemov(session *mgo.Session, project, name, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"beforemov": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetAftermov 함수는 item에 After mov값을 셋팅한다.
func SetAftermov(session *mgo.Session, project, name, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"aftermov": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskStatus 함수는 item에 task의 status 값을 셋팅한다.
func SetTaskStatus(session *mgo.Session, project, name, task, status string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	item, err := getItem(session, project, id)
	if err != nil {
		return err
	}
	statusNum := ""
	switch strings.ToLower(status) {
	case READY, "ready":
		statusNum = READY
	case ASSIGN, "assign":
		statusNum = ASSIGN
	case WIP, "wip":
		statusNum = WIP
	case CONFIRM, "confirm":
		statusNum = CONFIRM
	case DONE, "done":
		statusNum = DONE
	case OMIT, "omit":
		statusNum = OMIT
	case HOLD, "hold":
		statusNum = HOLD
	case OUT, "out":
		statusNum = OUT
	case NONE, "none":
		statusNum = NONE
	}
	if statusNum == "" {
		return errors.New("올바른 status가 아닙니다")
	}
	if _, found := item.Tasks[strings.ToLower(task)]; !found {
		return errors.New("task가 존재하지 않습니다")
	}
	t := item.Tasks[task]
	t.BeforeStatus = t.Status
	t.Status = statusNum
	item.Tasks[task] = t

	c := session.DB("project").C(project)
	item.Updatetime = time.Now().Format(time.RFC3339)
	item.updateStatus()
	err = c.Update(bson.M{"id": item.ID}, item)
	if err != nil {
		return err
	}
	return nil
}

// HasTask 함수는 item에 task가 존재하는 체크한다.
func HasTask(session *mgo.Session, project, name, task string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	item, err := getItem(session, project, id)
	if err != nil {
		return err
	}
	if _, found := item.Tasks[task]; !found {
		return fmt.Errorf("%s 프로젝트 %s에 %s Task가 존재하지 않습니다", project, name, task)
	}
	return nil
}

// SetAssignTask 함수는 item에 task의 assign을 셋팅한다.
func SetAssignTask(session *mgo.Session, project, name, taskname string, visable bool) (string, error) {
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
	item, err := getItem(session, project, id)
	if err != nil {
		return "", err
	}
	task := strings.ToLower(taskname)
	// 기존에 Task가 없다면 추가한다.
	if _, found := item.Tasks[task]; !found {
		t := Task{
			Title:  task,
			Status: ASSIGN,
		}
		item.Tasks[task] = t
	} else {
		// 이미 Task가 존재하면
		t := item.Tasks[task]
		t.Title = task
		t.Status = ASSIGN
		item.Tasks[task] = t
	}

	if !visable {
		delete(item.Tasks, task)
	}
	c := session.DB("project").C(project)
	item.Updatetime = time.Now().Format(time.RFC3339)
	item.updateStatus()
	err = c.Update(bson.M{"id": item.ID}, item)
	if err != nil {
		return "", err
	}
	return id, nil
}

// SetTaskUser 함수는 item에 task의 user 값을 셋팅한다.
func SetTaskUser(session *mgo.Session, project, name, task, user string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	item, err := getItem(session, project, id)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": item.ID}, bson.M{"$set": bson.M{"tasks." + task + ".user": user, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskDate 함수는 item에 task에 마감일을 셋팅한다.
func SetTaskDate(session *mgo.Session, project, name, task, date string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	item, err := getItem(session, project, id)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": item.ID}, bson.M{"$set": bson.M{"tasks." + task + ".date": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetDeadline2D 함수는 item에 2D마감일을 셋팅한다.
func SetDeadline2D(session *mgo.Session, project, name, date string) (string, error) {
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
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return id, err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"ddline2d": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, err
	}
	return id, nil
}

// SetDeadline3D 함수는 item에 3D마감일을 셋팅한다.
func SetDeadline3D(session *mgo.Session, project, name, date string) (string, error) {
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
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return id, err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"ddline3d": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, err
	}
	return id, nil
}

// SetTaskStartdate 함수는 item에 task의 startdate 값을 셋팅한다.
func SetTaskStartdate(session *mgo.Session, project, name, task, date string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	c := session.DB("project").C(project)
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".startdate": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskUserNote 함수는 item에 task의 user note 값을 셋팅한다.
func SetTaskUserNote(session *mgo.Session, project, name, task, usernote string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".usernote": usernote, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskPredate 함수는 item에 task의 predate 값을 셋팅한다.
func SetTaskPredate(session *mgo.Session, project, name, task, date string) (string, error) {
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
	item, err := getItem(session, project, id)
	if err != nil {
		return id, err
	}
	c := session.DB("project").C(project)
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return id, err
	}
	err = c.Update(bson.M{"id": item.ID}, bson.M{"$set": bson.M{"tasks." + task + ".predate": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, err
	}
	return id, nil
}

// SetShotType 함수는 item에 shot type을 셋팅한다.
func SetShotType(session *mgo.Session, project, name, shottype string) (string, error) {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return "", err
	}
	if typ == "asset" {
		return "", fmt.Errorf("%s 는 asset type 입니다. 변경할 수 없습니다", name)
	}
	id := name + "_" + typ
	err = validShottype(shottype)
	if err != nil {
		return id, err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"shottype": shottype, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, err
	}
	return id, nil
}

// SetOutputName 함수는 item에 Outputname 을 셋팅한다.
func SetOutputName(session *mgo.Session, project, name, outputname string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ == "asset" {
		return fmt.Errorf("%s 는 %s type 입니다. 변경할 수 없습니다", name, typ)
	}
	if outputname == "" {
		return errors.New("outputname 이 빈 문자열 입니다")
	}
	id := name + "_" + typ
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"outputname": outputname, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetRetimePlate 함수는 item에 RetimePlate를 셋팅한다.
func SetRetimePlate(session *mgo.Session, project, name, retimeplate string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ == "asset" {
		return fmt.Errorf("%s 는 %s type 입니다. retime plate를 설정할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"retimeplate": retimeplate, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetOCIOcc 함수는 item에 OCIO .cc를 셋팅한다.
func SetOCIOcc(session *mgo.Session, project, name, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ == "asset" {
		return fmt.Errorf("%s 는 %s type 입니다. 설정할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"ociocc": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetRollmedia 함수는 item에 Setellite Rollmedia를 셋팅한다.
func SetRollmedia(session *mgo.Session, project, name, rollmedia string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"rollmedia": rollmedia, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetRnum 함수는 샷에 롤넘버를 설정한다.
func SetRnum(session *mgo.Session, project, name, rnum string) (string, error) {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return "", err
	}
	if typ == "asset" {
		return "", fmt.Errorf("%s 는 %s type 입니다. 변경할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	item, err := getItem(session, project, id)
	if err != nil {
		return id, err
	}
	item.Rnum = rnum
	err = setItem(session, project, item)
	if err != nil {
		return id, err
	}
	return id, nil
}

// SetAssetType 함수는 item에 assettype을 셋팅한다.
func SetAssetType(session *mgo.Session, project, name, assettype string) (string, string, string, error) {
	_, err := validAssettype(assettype)
	if err != nil {
		return "", "", assettype, err
	}
	session.SetMode(mgo.Monotonic, true)
	err = HasProject(session, project)
	if err != nil {
		return "", "", assettype, err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return "", "", assettype, err
	}
	if typ != "asset" {
		return "", "", assettype, fmt.Errorf("%s 아이템은 %s 타입입니다. 처리할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	i, err := getItem(session, project, id)
	if err != nil {
		return id, "", assettype, err
	}
	var beforeType string
	beforeType = i.Assettype
	i.Assettype = assettype
	i.setAssettags()
	err = setItem(session, project, i)
	if err != nil {
		return id, beforeType, assettype, err
	}
	return id, beforeType, assettype, nil
}

// SetScanTimecodeIn 함수는 item에 Scan Timecode In을 셋팅한다.
func SetScanTimecodeIn(session *mgo.Session, project, name, timecode string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ == "asset" {
		return fmt.Errorf("%s 아이템은 %s 타입입니다. 처리할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"scantimecodein": timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetScanTimecodeOut 함수는 item에 Scan Timecode In을 셋팅한다.
func SetScanTimecodeOut(session *mgo.Session, project, name, timecode string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ == "asset" {
		return fmt.Errorf("%s 아이템은 %s 타입입니다. 처리할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"scantimecodeout": timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetJustTimecodeIn 함수는 item에 Just Timecode In을 셋팅한다.
func SetJustTimecodeIn(session *mgo.Session, project, name, timecode string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ == "asset" {
		return fmt.Errorf("%s 아이템은 %s 타입입니다. 처리할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"justtimecodein": timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetJustTimecodeOut 함수는 item에 Just Timecode In을 셋팅한다.
func SetJustTimecodeOut(session *mgo.Session, project, name, timecode string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ == "asset" {
		return fmt.Errorf("%s 아이템은 %s 타입입니다. 처리할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"justtimecodeout": timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetFinver 함수는 item에 파이널 버전을 셋팅한다.
func SetFinver(session *mgo.Session, project, name, version string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"finver": version, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetFindate 함수는 item에 최종 데이터 아웃풋 날짜를 셋팅한다.
func SetFindate(session *mgo.Session, project, name, date string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return err
	}
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"findate": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetCrowdAsset 함수는 item에 crowdtype을 설정한다.
func SetCrowdAsset(session *mgo.Session, project, name string) (string, bool, error) {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", false, err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return "", false, err
	}
	id := name + "_" + typ
	c := session.DB("project").C(project)
	item, err := getItem(session, project, id)
	if err != nil {
		return id, item.CrowdAsset, err
	}
	invertBool := !item.CrowdAsset
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"crowdasset": invertBool, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, invertBool, err
	}
	return id, invertBool, nil
}

// AddTag 함수는 item에 tag를 셋팅한다.
func AddTag(session *mgo.Session, project, name, inputTag string) (string, error) {
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
	i, err := getItem(session, project, id)
	if err != nil {
		return id, err
	}
	rmspaceTag := strings.Replace(inputTag, " ", "", -1) // 태그는 공백을 제거한다.
	for _, tag := range i.Tag {
		if rmspaceTag == tag {
			return id, errors.New(inputTag + "태그는 이미 존재하고 있습니다 추가할 수 없습니다")
		}
	}
	newTags := append(i.Tag, rmspaceTag)
	c := session.DB("project").C(project)
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tag": newTags, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, err
	}
	return id, nil
}

// RenameTag 함수는 item의 Tag를 리네임한다.
func RenameTag(session *mgo.Session, project, before, after string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("project").C(project)
	var items []Item
	err := c.Find(bson.M{}).All(&items)
	if err != nil {
		return err
	}
	for _, i := range items {
		beforeTags := i.Tag
		var newTags []string
		for _, t := range i.Tag {
			if t == before {
				newTags = append(newTags, after)
			} else {
				newTags = append(newTags, t)
			}
		}
		if !reflect.DeepEqual(beforeTags, newTags) {
			err = c.Update(bson.M{"id": i.ID}, bson.M{"$set": bson.M{"tag": newTags, "updatetime": time.Now().Format(time.RFC3339)}})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetTags 함수는 item에 tag를 교체한다.
func SetTags(session *mgo.Session, project, name string, tags []string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	i, err := getItem(session, project, id)
	if err != nil {
		return err
	}
	i.Tag = tags
	// 만약 태그에 권정보가 없더라도 권관련 태그는 날아가면 안된다. setItem을 이용한다.
	err = setItem(session, project, i)
	if err != nil {
		return err
	}
	return nil
}

// RmTag 함수는 item에 tag를 삭제한다.
func RmTag(session *mgo.Session, project, name string, inputTag string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	i, err := getItem(session, project, id)
	if err != nil {
		return err
	}
	var newTags []string
	for _, tag := range i.Tag {
		if inputTag == tag {
			continue
		}
		newTags = append(newTags, tag)
	}
	i.Tag = newTags
	// 만약 태그에 권정보가 없더라도 권관련 태그는 날아가면 안된다. setItem을 이용한다.
	err = setItem(session, project, i)
	if err != nil {
		return err
	}
	return nil
}

// SetNote 함수는 item에 작업내용을 추가한다. 노트내용과 에러를 반환한다.
func SetNote(session *mgo.Session, project, name, userID, text string, overwrite bool) (string, error) {
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
	var note string
	if overwrite {
		note = text
	} else {
		i, err := getItem(session, project, id)
		if err != nil {
			return "", err
		}
		if strings.HasSuffix(i.Note.Text, "\n") {
			note = text + i.Note.Text
		} else {
			note = text + "\n " + i.Note.Text
		}
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"note.text": note, "note.author": userID, "note.date": time.Now().Format(time.RFC3339), "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return "", err
	}
	return note, nil
}

// AddComment 함수는 item에 수정사항을 추가한다.
func AddComment(session *mgo.Session, project, name, userID, date, text, media string) (string, error) {
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
	i, err := getItem(session, project, id)
	if err != nil {
		return id, err
	}
	c := Comment{
		Date:   date,
		Author: userID,
		Text:   text,
		Media:  media,
	}
	i.Comments = append(i.Comments, c)
	err = setItem(session, project, i)
	if err != nil {
		return id, err
	}
	return id, nil
}

// RmComment 함수는 item에 수정사항을 삭제합니다. 로그처리를 위해서 삭제 내용을 반환합니다.
func RmComment(session *mgo.Session, project, name, userID, date string) (string, string, error) {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", "", err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return "", "", err
	}
	id := name + "_" + typ
	i, err := getItem(session, project, id)
	if err != nil {
		return id, "", err
	}
	var newComments []Comment
	var removeText string
	for _, comment := range i.Comments {
		if comment.Date == date {
			removeText = comment.Text
			continue
		}
		newComments = append(newComments, comment)
	}
	i.Comments = newComments
	err = setItem(session, project, i)
	if err != nil {
		return id, "", err
	}
	return id, removeText, nil
}

// AddSource 함수는 item에 소스링크를 추가한다.
func AddSource(session *mgo.Session, project, name, userID, title, path string) (string, error) {
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
	i, err := getItem(session, project, id)
	if err != nil {
		return id, err
	}
	for _, i := range i.Sources {
		if i.Title == title {
			return id, errors.New(title + "이 이미 존재합니다.")
		}
	}
	s := Source{}
	s.Date = time.Now().Format(time.RFC3339)
	s.Author = userID
	s.Title = title
	s.Path = path
	i.Sources = append(i.Sources, s)
	err = setItem(session, project, i)
	if err != nil {
		return id, err
	}
	return id, nil
}

// AddReference 함수는 item에 소스링크를 추가한다.
func AddReference(session *mgo.Session, project, name, userID, title, path string) (string, error) {
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
	i, err := getItem(session, project, id)
	if err != nil {
		return id, err
	}
	r := Source{}
	r.Date = time.Now().Format(time.RFC3339)
	r.Author = userID
	r.Title = title
	r.Path = path
	i.References = append(i.References, r)
	err = setItem(session, project, i)
	if err != nil {
		return id, err
	}
	return id, nil
}

// RmSource 함수는 item에서 소스를 삭제합니다.
func RmSource(session *mgo.Session, project, name, title string) (string, error) {
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
	i, err := getItem(session, project, id)
	if err != nil {
		return id, err
	}
	var newSources []Source
	for _, source := range i.Sources {
		if source.Title == title {
			continue
		}
		newSources = append(newSources, source)
	}
	i.Sources = newSources
	err = setItem(session, project, i)
	if err != nil {
		return id, err
	}
	return id, nil
}

// RmReference 함수는 item에서 레퍼런스를 삭제합니다.
func RmReference(session *mgo.Session, project, name, title string) (string, error) {
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
	i, err := getItem(session, project, id)
	if err != nil {
		return id, err
	}
	var newReferences []Source
	for _, ref := range i.References {
		if ref.Title == title {
			continue
		}
		newReferences = append(newReferences, ref)
	}
	i.References = newReferences
	err = setItem(session, project, i)
	if err != nil {
		return id, err
	}
	return id, nil
}

// GetTask 함수는 item의 Task 정보를 반환한다.
func GetTask(session *mgo.Session, project, name, task string) (string, Task, error) {
	session.SetMode(mgo.Monotonic, true)
	typ, err := Type(session, project, name)
	if err != nil {
		return "", Task{}, err
	}
	id := name + "_" + typ
	i, err := getItem(session, project, id)
	if err != nil {
		return id, Task{}, err
	}
	if _, found := i.Tasks[task]; !found {
		return id, Task{}, errors.New("task가 존재하지 않습니다")
	}
	return id, i.Tasks[task], nil
}

// GetShottype 함수는 item의 Shottype 정보를 반환한다.
func GetShottype(session *mgo.Session, project, name string) (string, error) {
	session.SetMode(mgo.Monotonic, true)
	typ, err := Type(session, project, name)
	if err != nil {
		return "", err
	}
	id := name + "_" + typ
	i, err := getItem(session, project, id)
	if err != nil {
		return "", err
	}
	return i.Shottype, nil
}
