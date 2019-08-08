package main

import (
	"fmt"
	"log"
	"os/user"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

func addProjectCmd(name string) {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	p := *NewProject(name)
	err = addProject(session, p)
	if err != nil {
		log.Fatal(err)
	}
	// 최초 프로젝트가 생성될 때 temp_asset 더미 아이템이 생성된다.
	// 이 자료구조로 비상시 DB에 접근, 구조를 디버깅할 때 사용한다.
	i := Item{
		Project:    name,
		Name:       "temp",
		Type:       "asset",
		Slug:       "temp_asset",
		ID:         "temp_asset",
		Updatetime: time.Now().Format(time.RFC3339),
		Status:     NONE,
	}
	err = addItem(session, name, i)
	if err != nil {
		log.Fatal(err)
	}
}

func rmProjectCmd(name string) {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if user.Username != "root" {
		log.Fatal("루트계정이 아닙니다.")
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	err = rmProject(session, name)
	if err != nil {
		log.Fatal(err)
	}
}

func addShotItemCmd(project, name, typ, platesize, scanname, scantimecodein, scantimecodeout, justtimecodein, justtimecodeout string, scanframe, scanin, scanout, platein, plateout, justin, justout int) {
	if !regexpShotname.MatchString(name) {
		log.Fatal("샷 이름 규칙이 아닙니다.")
	}
	seq := strings.Split(name, "_")[0]
	now := time.Now().Format(time.RFC3339)
	i := Item{
		Project:    project,
		Name:       name,
		Seq:        seq,
		Type:       typ,
		Slug:       name + "_" + typ,
		ID:         name + "_" + typ,
		Status:     ASSIGN,
		Thumpath:   fmt.Sprintf("/%s/%s_%s.jpg", project, name, typ),
		Platepath:  fmt.Sprintf("/show/%s/seq/%s/%s/plate/", project, seq, name),
		Thummov:    fmt.Sprintf("/show/%s/seq/%s/%s/plate/%s_%s.mov", project, seq, name, name, typ),
		Dataname:   scanname, // 일반적인 프로젝트는 스캔네임과 데이터네임이 같다. PM의 노가다를 줄이기 위해서 기본적으로 같은값이 들어가고 추후 수동처리해야하는 부분은 손으로 수정한다.
		Scanname:   scanname,
		Scantime:   now,
		Platesize:  platesize,
		Updatetime: now,
	}
	i.Comp.Status = ASSIGN // 샷의 경우 합성팀을 무조건 거쳐야 한다. Assign상태로 만든다.

	if scanframe != 0 {
		i.ScanFrame = scanframe
	}
	if scantimecodein != "" {
		i.ScanTimecodeIn = scantimecodein
	}
	if scantimecodeout != "" {
		i.ScanTimecodeOut = scantimecodeout
	}
	if justtimecodein != "" {
		i.JustTimecodeIn = justtimecodein
	}
	if justtimecodeout != "" {
		i.JustTimecodeOut = justtimecodeout
	}
	if scanin != -1 {
		i.ScanIn = scanin
	}
	if scanout != -1 {
		i.ScanOut = scanout
	}
	if platein != -1 && justin == -1 {
		i.PlateIn = platein
		i.JustIn = platein
	}
	if plateout != -1 && justout == -1 {
		i.PlateOut = plateout
		i.JustOut = plateout
	}
	if justin != -1 {
		i.JustIn = justin
	}
	if justout != -1 {
		i.JustOut = justout
	}
	i.Project = project
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// 현장데이터가 존재하는지 체크한다.
	rollmedia := Scanname2RollMedia(scanname)
	if hasSetelliteItems(session, project, rollmedia) {
		i.Rollmedia = rollmedia
	}

	err = addItem(session, project, i)
	if err != nil {
		log.Fatal(err)
	}
}

func addAssetItemCmd(project, name, typ, assettype, assettags string) {
	if assettype == "" {
		log.Fatal("assettype을 입력해주세요.")
	}
	// 유효한 에셋타입인지 체크.
	err := validAssettype(assettype)
	if err != nil {
		log.Fatal(err)
	}
	i := Item{
		Project:    project,
		Name:       name,
		Type:       typ,
		Slug:       name + "_" + typ,
		ID:         name + "_" + typ,
		Status:     NONE,
		Updatetime: time.Now().Format(time.RFC3339),
		Assettype:  assettype,
		Assettags:  []string{},
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	if assettags == "" {
		log.Fatal("에셋 생성시 assettags가 필요합니다.")
	}
	for _, tag := range Str2Tags(assettags) {
		if tag == "assembly" {
			i.Assettags = append(i.Assettags, name) //에셈블리 추가시 자기 자신도 태그로 포함되어야 한다.
		}
		i.Assettags = append(i.Assettags, tag)
	}
	err = addItem(session, project, i)
	if err != nil {
		log.Fatal(err)
	}
}

// addOtherItemCmd함수는 Shot, Asset 이 아닌 나머지 아이템을 추가하는 함수이다.
func addOtherItemCmd(project, name, typ, platesize, scanname, scantimecodein, scantimecodeout, justtimecodein, justtimecodeout string, scanframe, scanin, scanout, platein, plateout, justin, justout int) {
	if !regexpShotname.MatchString(name) {
		log.Fatal("소스, 재스캔 이름 규칙이 아닙니다.")
	}
	seq := strings.Split(name, "_")[0]
	now := time.Now().Format(time.RFC3339)
	i := Item{
		Project:    project,
		Name:       name,
		Seq:        seq,
		Type:       typ,
		Slug:       name + "_" + typ,
		ID:         name + "_" + typ,
		Platepath:  fmt.Sprintf("/show/%s/seq/%s/%s/plate/", project, seq, name),
		Thummov:    fmt.Sprintf("/show/%s/seq/%s/%s/plate/%s_%s.mov", project, seq, name, name, typ),
		Status:     NONE,
		Dataname:   scanname, // 일반적인 프로젝트는 스캔네임과 데이터네임이 같다. PM의 노가다를 줄이기 위해서 기본적으로 같은값이 들어가고 추후 수동처리해야하는 부분은 손으로 수정한다.
		Scanname:   scanname,
		Scantime:   now,
		Updatetime: now,
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	// org1, src1 같은 아이템도 키코드가 들어가야 한다.
	if scanframe != 0 {
		i.ScanFrame = scanframe
	}
	if scantimecodein != "" {
		i.ScanTimecodeIn = scantimecodein
	}
	if scantimecodeout != "" {
		i.ScanTimecodeOut = scantimecodeout
	}
	if justtimecodein != "" {
		i.JustTimecodeIn = justtimecodein
	}
	if justtimecodeout != "" {
		i.JustTimecodeOut = justtimecodeout
	}
	if scanin != -1 {
		i.ScanIn = scanin
	}
	if scanout != -1 {
		i.ScanOut = scanout
	}
	if platein != -1 && justin == -1 {
		i.PlateIn = platein
		i.JustIn = platein
	}
	if plateout != -1 && justout == -1 {
		i.PlateOut = plateout
		i.JustOut = plateout
	}
	if justin != -1 {
		i.JustIn = justin
	}
	if justout != -1 {
		i.JustOut = justout
	}

	// 현장데이터가 존재하는지 체크한다.
	rollmedia := Scanname2RollMedia(scanname)
	if hasSetelliteItems(session, project, rollmedia) {
		i.Rollmedia = rollmedia
	}

	err = addItem(session, project, i)
	if err != nil {
		log.Fatal(err)
	}
}

func rmItemCmd(project, name, typ string) {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if user.Username != "root" {
		log.Fatal("루트계정이 아닙니다.")
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	err = rmItem(session, project, name, typ)
	if err != nil {
		log.Fatal(err)
	}
}
