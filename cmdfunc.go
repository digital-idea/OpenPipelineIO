package main

import (
	"bytes"
	"html/template"
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
	now := time.Now().Format(time.RFC3339)
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	initStatusID, err := GetInitStatusID(session)
	if err != nil {
		log.Fatal(err)
	}
	admin, err := GetAdminSetting(session)
	if err != nil {
		log.Fatal(err)
	}
	i := Item{
		Project:    project,
		Name:       name,
		Type:       typ,
		ID:         name + "_" + typ,
		Status:     ASSIGN, // legacy
		StatusV2:   initStatusID,
		Scanname:   scanname,
		Dataname:   scanname, // 보통 스캔네임과 데이터네임은 같다. 데이터 입력자의 노동을 줄이기 위해 기본적으로 동일값을 넣고, 필요시 수정한다.
		Scantime:   now,
		Platesize:  platesize,
		Updatetime: now,
		UseType:    typ, // 최초 생성시 사용타입은 자신의 Type과 같다.
		Season:     *flagSeason,
		Episode:    *flagEpisode,
	}
	i.Tasks = make(map[string]Task)
	i.SetSeq()
	i.SetCut()
	// 썸네일 경로 자동설정
	if *flagThumbnailImagePath != "" {
		i.Thumpath = *flagThumbnailImagePath
	} else {
		// 만약 빈값이라면 adminSetting의 설정값을 이용해서 설정한다.
		var thumbnailImagePath bytes.Buffer
		thumbnailImagePathTmpl, err := template.New("thumbnailImagePath").Parse(admin.ThumbnailImagePath)
		if err != nil {
			log.Fatal(err)
		}
		err = thumbnailImagePathTmpl.Execute(&thumbnailImagePath, i)
		if err != nil {
			log.Fatal(err)
		}
		i.Thumpath = thumbnailImagePath.String()
	}
	// 썸네일Mov 경로 자동설정
	if *flagThumbnailMovPath != "" {
		i.Thummov = *flagThumbnailMovPath
	} else {
		// 만약 빈값이라면 adminSetting의 설정값을 이용해서 설정한다.
		var thumbnailMovPath bytes.Buffer
		thumbnailMovPathTmpl, err := template.New("thumbnailMovPath").Parse(admin.ThumbnailMovPath)
		if err != nil {
			log.Fatal(err)
		}
		err = thumbnailMovPathTmpl.Execute(&thumbnailMovPath, i)
		if err != nil {
			log.Fatal(err)
		}
		i.Thummov = thumbnailMovPath.String()
	}
	// 플레이트 경로 자동설정
	if *flagPlatePath != "" {
		i.Platepath = *flagPlatePath
	} else {
		// 만약 빈값이라면 adminSetting의 설정값을 이용해서 설정한다.
		var platePath bytes.Buffer
		platePathTmpl, err := template.New("platePath").Parse(admin.PlatePath)
		if err != nil {
			log.Fatal(err)
		}
		err = platePathTmpl.Execute(&platePath, i)
		if err != nil {
			log.Fatal(err)
		}
		i.Platepath = platePath.String()
	}
	tasks, err := AllTaskSettings(session)
	if err != nil {
		log.Fatal(err)
	}
	for _, task := range tasks {
		if !task.InitGenerate {
			continue
		}
		if task.Type != "shot" {
			continue
		}
		t := Task{
			Title:    task.Name,
			Status:   ASSIGN, // 샷의 경우 합성팀을 무조건 거쳐야 한다. Assign상태로 만든다. // legacy
			StatusV2: initStatusID,
		}
		i.Tasks[task.Name] = t
	}
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
	if platein != -1 {
		i.PlateIn = platein
		i.JustIn = platein
	}
	if plateout != -1 {
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
	_, err := validAssettype(assettype)
	if err != nil {
		log.Fatal(err)
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	initStatusID, err := GetInitStatusID(session)
	if err != nil {
		log.Fatal(err)
	}
	i := Item{
		Project:    project,
		Name:       name,
		Type:       typ,
		ID:         name + "_" + typ,
		Status:     ASSIGN, // legacy
		StatusV2:   initStatusID,
		Updatetime: time.Now().Format(time.RFC3339),
		Assettype:  assettype,
		Assettags:  []string{},
		Season:     *flagSeason,
		Episode:    *flagEpisode,
	}

	tasks, err := AllTaskSettings(session)
	if err != nil {
		log.Fatal(err)
	}
	i.Tasks = make(map[string]Task)
	for _, task := range tasks {
		if !task.InitGenerate {
			continue
		}
		if task.Type != "asset" {
			continue
		}
		t := Task{
			Title:    task.Name,
			Status:   ASSIGN, // 샷의 경우 합성팀을 무조건 거쳐야 한다. Assign상태로 만든다. // legacy
			StatusV2: initStatusID,
		}
		i.Tasks[task.Name] = t
	}
	if assettags == "" {
		log.Fatal("에셋 생성시 assettags가 필요합니다.")
	}
	for _, tag := range Str2List(assettags) {
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
	now := time.Now().Format(time.RFC3339)
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	admin, err := GetAdminSetting(session)
	if err != nil {
		log.Fatal(err)
	}
	i := Item{
		Project:    project,
		Name:       name,
		Type:       typ,
		ID:         name + "_" + typ,
		Status:     NONE,
		StatusV2:   "none",
		Dataname:   scanname, // 일반적인 프로젝트는 스캔네임과 데이터네임이 같다. PM의 노가다를 줄이기 위해서 기본적으로 같은값이 들어가고 추후 수동처리해야하는 부분은 손으로 수정한다.
		Scanname:   scanname,
		Scantime:   now,
		Updatetime: now,
		Season:     *flagSeason,
		Episode:    *flagEpisode,
	}
	i.SetSeq()
	platePath := *flagPlatePath
	if *flagPlatePath != "" {
		i.Platepath = *flagPlatePath
	} else {
		// 만약 빈값이라면 adminSetting의 설정값을 이용해서 설정한다.
		var platePath bytes.Buffer
		platePathTmpl, err := template.New("platePath").Parse(admin.PlatePath)
		if err != nil {
			log.Fatal(err)
		}
		err = platePathTmpl.Execute(&platePath, i)
		if err != nil {
			log.Fatal(err)
		}
		i.Platepath = platePath.String()
	}
	if *flagThumbnailMovPath != "" {
		i.Thummov = *flagThumbnailMovPath
	} else {
		// 만약 빈값이라면 adminSetting의 설정값을 이용해서 설정한다.
		var thumbnailMovPath bytes.Buffer
		thumbnailMovPathTmpl, err := template.New("thumbnailMovPath").Parse(admin.ThumbnailMovPath)
		if err != nil {
			log.Fatal(err)
		}
		err = thumbnailMovPathTmpl.Execute(&thumbnailMovPath, i)
		if err != nil {
			log.Fatal(err)
		}
		i.Thummov = thumbnailMovPath.String()
	}
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
	if platein != -1 {
		i.PlateIn = platein
		i.JustIn = platein
	}
	if plateout != -1 {
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
	// 만약 typ에 src 문자로 시작하면 소스로 판단하고 기존 item에 자동으로 소스 등록을 진행한다.
	if strings.HasPrefix(typ, "src") {
		_, err = AddSource(session, project, name, "scantool", name+"_"+typ, platePath)
		if err != nil {
			log.Println(err)
		}
	}
	// 만약 typ이 ref 문자로 시작하면 레퍼런스로 판단하고 기존 item에 자동으로 레퍼런스 등록을 진행한다.
	if strings.HasPrefix(typ, "ref") {
		_, err = AddReference(session, project, name, "scantool", name+"_"+typ, platePath)
		if err != nil {
			log.Println(err)
		}
	}
	// org1, left1 형태의 아이템이 처리되면 org, left 아이템의 .UseType을 추가해준다.
	// 이 값은 썸네일을 업데이트하고, 아티스트가 재스캔 되었을 때 사용할 타입의 알람으로 사용된다.
	if strings.Contains(typ, "org") || strings.Contains(typ, "left") {
		err = SetUseType(session, project, i.ID, typ)
		if err != nil {
			log.Println(err)
		}
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
