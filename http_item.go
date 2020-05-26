package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/digital-idea/dipath"
	"github.com/disintegration/imaging"
	"golang.org/x/sys/unix"
	"gopkg.in/mgo.v2"
)

// handleSearchSubmit 함수는 검색창의 옵션을 파싱하고 검색 URI로 리다이렉션 한다. // legacy
func handleSearchSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	project := r.FormValue("Project")
	searchword := r.FormValue("Searchword")
	sortkey := r.FormValue("Sortkey")
	assign := str2bool(r.FormValue("Assign"))   // legacy
	ready := str2bool(r.FormValue("Ready"))     // legacy
	wip := str2bool(r.FormValue("Wip"))         // legacy
	confirm := str2bool(r.FormValue("Confirm")) // legacy
	done := str2bool(r.FormValue("Done"))       // legacy
	omit := str2bool(r.FormValue("Omit"))       // legacy
	hold := str2bool(r.FormValue("Hold"))       // legacy
	out := str2bool(r.FormValue("Out"))         // legacy
	none := str2bool(r.FormValue("None"))       // legacy
	template := r.FormValue("Template")
	searchbarTemplate := r.FormValue("SearchbarTemplate")
	task := r.FormValue("Task")
	truestatus := r.FormValue("truestatus")
	// 아래 코드는 임시로 사용한다.
	if truestatus == "" {
		var statuslist []string
		if assign {
			statuslist = append(statuslist, "assign")
		}
		if ready {
			statuslist = append(statuslist, "ready")
		}
		if wip {
			statuslist = append(statuslist, "wip")
		}
		if confirm {
			statuslist = append(statuslist, "confirm")
		}
		if done {
			statuslist = append(statuslist, "done")
		}
		if omit {
			statuslist = append(statuslist, "omit")
		}
		if out {
			statuslist = append(statuslist, "out")
		}
		if hold {
			statuslist = append(statuslist, "hold")
		}
		if none {
			statuslist = append(statuslist, "none")
		}
		truestatus = strings.Join(statuslist, ",")
	}
	redirectURL := fmt.Sprintf(`/inputmode?project=%s&searchword=%s&sortkey=%s&assign=%t&ready=%t&wip=%t&confirm=%t&done=%t&omit=%t&hold=%t&out=%t&none=%t&template=%s&task=%s&searchbartemplate=%s&truestatus=%s`,
		project,
		searchword,
		sortkey,
		assign,
		ready,
		wip,
		confirm,
		done,
		omit,
		hold,
		out,
		none,
		template,
		task,
		searchbarTemplate,
		truestatus,
	)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// handleSearchSubmitV2 함수는 검색창의 옵션을 파싱하고 검색 URI로 리다이렉션 한다.
func handleSearchSubmitV2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	project := r.FormValue("Project")
	searchword := r.FormValue("Searchword")
	sortkey := r.FormValue("Sortkey")
	template := r.FormValue("Template")
	searchbarTemplate := r.FormValue("SearchbarTemplate")
	task := r.FormValue("Task")
	// status를 체크할 때 마다 truestatus form에 값이 추가되어야 한다.
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	statuslist, err := AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var truestatus []string
	for _, status := range statuslist {
		if "true" == r.FormValue(status.ID) {
			truestatus = append(truestatus, status.ID)
		}
	}
	redirectURL := fmt.Sprintf(`/inputmode?project=%s&searchword=%s&sortkey=%s&template=%s&task=%s&searchbartemplate=%s&truestatus=%s`,
		project,
		searchword,
		sortkey,
		template,
		task,
		searchbarTemplate,
		strings.Join(truestatus, ","),
	)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// handleItemDetail 함수는 아이템 디테일 페이지를 출력한다.
func handleItemDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	id := q.Get("id")
	searchbarTemplate := q.Get("searchbartemplate")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		User        User
		Projectlist []string
		Devmode     bool
		Projectinfo Project
		SearchOption
		Dilog               string
		Wfs                 string
		Item                Item
		TasksettingOrderMap map[string]float64
		Status              []Status
	}
	rcp := recipe{}
	rcp.Wfs = *flagWFS
	rcp.Dilog = *flagDILOG
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Devmode = *flagDevmode
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rcp.SearchOption.Project != "" {
		rcp.Projectinfo, err = getProject(session, rcp.SearchOption.Project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	rcp.SearchOption.SearchbarTemplate = searchbarTemplate
	tasks, err := AllTaskSettings(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.TasksettingOrderMap = make(map[string]float64)
	for _, t := range tasks {
		rcp.TasksettingOrderMap[t.Name] = t.Order
	}
	rcp.Status, err = AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Item, err = getItem(session, project, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "detail", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleEditItem 함수는 Item 편집페이지이다.
func handleEditItem(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	q := r.URL.Query()
	project := q.Get("project")
	slug := q.Get("slug")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		ID      string
		Project Project
		Item    Item
		Devmode bool
		User
		SearchOption
	}
	rcp := recipe{}
	rcp.Devmode = *flagDevmode
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.SearchOption = handleRequestToSearchOption(r)
	rcp.Project, err = getProject(session, project)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Item, err = getItem(session, project, slug)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "edititem", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleEditedItem 함수는 Item 이 수정된 이후 이동하는 페이지이다.
func handleEditedItem(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "editeditem", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleIndex 함수는 index 페이지이다.
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // 등록되지 않은 URL은 NotFound 처리한다.
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	type recipe struct {
		SearchOption
	}
	rcp := recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 아무 상태도 선택되어있지 않다면 기본 상태설정으로 변경한다. // legacy
	if rcp.SearchOption.isStatusOff() {
		err := rcp.SearchOption.setStatusDefaultV1()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 아무 상태도 선택되어있지 않다면 기본 상태설정으로 변경한다. // V2
	if len(rcp.SearchOption.TrueStatus) == 0 {
		err := rcp.SearchOption.setStatusDefaultV2()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 만약 SearchbarTempate 값이 설정되어 있다면 해당 SearchbarTemplate으로 변경한다.
	rcp.SearchOption.SearchbarTemplate = "searchbarV1"
	q := r.URL.Query()
	searchbarTemplate := q.Get("searchbartemplate")
	if searchbarTemplate != "" {
		rcp.SearchOption.SearchbarTemplate = searchbarTemplate
	}

	url := fmt.Sprintf("/inputmode?project=%s&sortkey=%s&template=index&searchbartemplate=%s&endpoint=searchv2&assign=%t&ready=%t&wip=%t&confirm=%t&done=%t&omit=%t&hold=%t&out=%t&none=%t&task=%s&searchword=%s&truestatus=%s",
		rcp.SearchOption.Project,
		rcp.SearchOption.Sortkey,
		rcp.SearchOption.SearchbarTemplate,
		rcp.SearchOption.Assign,  // legacy
		rcp.SearchOption.Ready,   // legacy
		rcp.SearchOption.Wip,     // legacy
		rcp.SearchOption.Confirm, // legacy
		rcp.SearchOption.Done,    // legacy
		rcp.SearchOption.Omit,    // legacy
		rcp.SearchOption.Hold,    // legacy
		rcp.SearchOption.Out,     // legacy
		rcp.SearchOption.None,    // legacy
		rcp.SearchOption.Task,
		rcp.SearchOption.Searchword,
		strings.Join(rcp.SearchOption.TrueStatus, ","),
	)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// handleAddShot 함수는 shot을 추가하는 페이지이다.
func handleAddShot(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		User        User
		Projectlist []string
		Devmode     bool
		SearchOption
	}
	rcp := recipe{}
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Devmode = *flagDevmode
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "addShot", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddShotSubmit 함수는 shot을 생성한다.
func handleAddShotSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	project := r.FormValue("Project")
	name := r.FormValue("Name")
	stereo := r.FormValue("Stereo")
	mkdir := str2bool(r.FormValue("Mkdir"))
	setRendersize := str2bool(r.FormValue("SetRendersize"))
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_'
	}
	names := strings.FieldsFunc(name, f)
	type Shot struct {
		Project string
		Name    string
		Error   string
	}
	var success []Shot
	var fails []Shot
	tasks, err := AllTaskSettings(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	admin, err := GetAdminSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pinfo, err := getProject(session, project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, n := range names {
		if n == " " || n == "" { // 사용자가 실수로 여러개의 스페이스를 추가할 수 있다.
			continue
		}
		s := Shot{}
		s.Project = project
		s.Name = n
		if !regexpShotname.MatchString(n) {
			s.Error = "SS_0010 형식의 이름이 아닙니다"
			fails = append(fails, s)
			continue
		}
		i := Item{}
		i.Name = n
		i.Type = "org"
		if stereo == "true" {
			i.Type = "left"
		}
		i.Project = project
		i.ID = i.Name + "_" + i.Type
		i.Seq = strings.Split(i.Name, "_")[0]
		i.Cut = strings.Split(i.Name, "_")[1]
		i.Shottype = "2d"
		i.Thumpath = fmt.Sprintf("/%s/%s_%s.jpg", i.Project, i.Name, i.Type)
		i.Platepath = fmt.Sprintf("/show/%s/seq/%s/%s/plate/", i.Project, i.Seq, i.Name)
		i.Thummov = fmt.Sprintf("/show/%s/seq/%s/%s/plate/%s_%s.mov", i.Project, i.Seq, i.Name, i.Name, i.Type)
		i.Scantime = time.Now().Format(time.RFC3339)
		i.Updatetime = time.Now().Format(time.RFC3339)
		i.Status = ASSIGN
		i.StatusV2 = "assign"
		if setRendersize {
			width := int(float64(pinfo.PlateWidth) * admin.DefaultScaleRatioOfUndistortionPlate)
			height := int(float64(pinfo.PlateHeight) * admin.DefaultScaleRatioOfUndistortionPlate)
			i.Platesize = fmt.Sprintf("%dx%d", pinfo.PlateWidth, pinfo.PlateHeight)
			i.Dsize = fmt.Sprintf("%dx%d", width, height)
			i.Rendersize = fmt.Sprintf("%dx%d", width, height)
		}
		// 기본적으로 생성해야할 Task를 추가한다.
		i.Tasks = make(map[string]Task)
		for _, task := range tasks {
			if !task.InitGenerate {
				continue
			}
			if task.Type != "shot" {
				continue
			}
			t := Task{
				Title:    task.Name,
				Status:   ASSIGN, // legacy
				StatusV2: "assign",
			}
			i.Tasks[task.Name] = t
		}
		err = addItem(session, project, i)
		if err != nil {
			s.Error = err.Error()
			fails = append(fails, s)
			continue
		}
		// 폴더 생성 옵션을 체크하면 폴더를 생성한다.
		if mkdir {
			adminSetting, err := GetAdminSetting(session)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			var shotRootPath bytes.Buffer
			var seqPath bytes.Buffer
			var shotPath bytes.Buffer
			shotRootPathTmpl, err := template.New("shotRootPath").Parse(adminSetting.ShotRootPath)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			err = shotRootPathTmpl.Execute(&shotRootPath, i)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			seqPathTmpl, err := template.New("seqPath").Parse(adminSetting.SeqPath)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			err = seqPathTmpl.Execute(&seqPath, i)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			shotPathTmpl, err := template.New("shotPath").Parse(adminSetting.ShotPath)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			err = shotPathTmpl.Execute(&shotPath, i)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			// Umask를 셋팅한다.
			if adminSetting.Umask == "" {
				unix.Umask(0)
			} else {
				umask, err := strconv.Atoi(adminSetting.Umask)
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
				unix.Umask(umask)
			}
			// ShotRootPath를 점검하고 없다면 경로를 생성한다.
			if _, err := os.Stat(shotRootPath.String()); os.IsNotExist(err) {
				per, err := strconv.ParseInt(adminSetting.ShotRootPathPermission, 8, 64)
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
				err = os.MkdirAll(shotRootPath.String(), os.FileMode(per))
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
				// admin 셋팅에 UID와 GID가 선언되어 있다면, 설정을 해줍니다.
				if adminSetting.ShotRootPathUID != "" && adminSetting.ShotRootPathGID != "" {
					uid, err := strconv.Atoi(adminSetting.ShotRootPathUID)
					if err != nil {
						s.Error = err.Error()
						fails = append(fails, s)
						continue
					}
					gid, err := strconv.Atoi(adminSetting.ShotRootPathGID)
					if err != nil {
						s.Error = err.Error()
						fails = append(fails, s)
						continue
					}
					err = os.Chown(shotRootPath.String(), uid, gid)
					if err != nil {
						s.Error = err.Error()
						fails = append(fails, s)
						continue
					}
				}
			}
			// 개별 샷 경로를 생성한다.
			per, err := strconv.ParseInt(adminSetting.ShotPathPermission, 8, 64)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			err = os.MkdirAll(shotPath.String(), os.FileMode(per))
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			// admin 셋팅에 UID와 GID가 선언되어 있다면, 설정을 해줍니다.
			if adminSetting.ShotPathUID != "" && adminSetting.ShotPathGID != "" {
				uid, err := strconv.Atoi(adminSetting.ShotPathUID)
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
				gid, err := strconv.Atoi(adminSetting.ShotPathGID)
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
				err = os.Chown(shotPath.String(), uid, gid)
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
			}
			// Seq 경로의 Permission을 설정한다.
			if adminSetting.SeqPathPermission != "" {
				per, err := strconv.ParseInt(adminSetting.SeqPathPermission, 8, 64)
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
				err = os.Chmod(seqPath.String(), os.FileMode(per))
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
			}
			// Seq 경로의 UID, GID를 설정한다.
			if adminSetting.SeqPathUID != "" && adminSetting.SeqPathGID != "" {
				uid, err := strconv.Atoi(adminSetting.SeqPathUID)
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
				gid, err := strconv.Atoi(adminSetting.SeqPathGID)
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
				err = os.Chown(seqPath.String(), uid, gid)
				if err != nil {
					s.Error = err.Error()
					fails = append(fails, s)
					continue
				}
			}
		}
		success = append(success, s)
		// slack log
		err = slacklog(session, project, fmt.Sprintf("Add Shot: %s\nProject: %s, Author: %s", n, project, ssid.ID))
		if err != nil {
			log.Println(err)
			continue
		}
	}

	type recipe struct {
		Success []Shot
		Fails   []Shot
		User    User
		Devmode bool
		SearchOption
	}
	rcp := recipe{}
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Devmode = *flagDevmode
	rcp.Success = success
	rcp.Fails = fails
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u

	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addShot_success", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddAsset 함수는 asset을 추가하는 페이지이다.
func handleAddAsset(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		User        User
		Projectlist []string
		Devmode     bool
		SearchOption
	}
	rcp := recipe{}
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Devmode = *flagDevmode
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "addAsset", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddAssetSubmit 함수는 asset을 생성한다.
func handleAddAssetSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	project := r.FormValue("Project")
	name := r.FormValue("Name")
	assettype := r.FormValue("Assettype")
	construction := r.FormValue("Construction")
	crowdAsset := str2bool(r.FormValue("CrowdAsset"))
	mkdir := str2bool(r.FormValue("Mkdir"))
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_'
	}
	names := strings.FieldsFunc(name, f)
	type Asset struct {
		Project string
		Name    string
		Error   string
	}
	var success []Asset
	var fails []Asset
	tasks, err := AllTaskSettings(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, n := range names {
		if n == " " || n == "" { // 사용자가 실수로 여러개의 스페이스를 추가할 수 있다.
			continue
		}
		a := Asset{}
		a.Project = project
		a.Name = n
		if !regexpAssetname.MatchString(n) {
			a.Error = "stone 또는 stone01 형식의 이름이 아닙니다"
			fails = append(fails, a)
			continue
		}
		i := Item{}
		i.Name = n
		i.Type = "asset"
		i.Project = project
		i.ID = i.Name + "_" + i.Type
		i.Status = ASSIGN
		i.StatusV2 = "assign"
		i.Updatetime = time.Now().Format(time.RFC3339)
		i.Assettype = assettype
		i.Assettags = []string{assettype, construction}
		i.CrowdAsset = crowdAsset
		// 기본적으로 생성해야할 Task를 추가한다.
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
				Status:   ASSIGN, // legacy
				StatusV2: "assign",
			}
			i.Tasks[task.Name] = t
		}
		err = addItem(session, project, i)
		if err != nil {
			a.Error = err.Error()
			fails = append(fails, a)
			continue
		}
		// 폴더 생성 옵션을 체크하면 폴더를 생성한다.
		if mkdir {
			adminSetting, err := GetAdminSetting(session)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			var assetRootPath bytes.Buffer
			var assetTypePath bytes.Buffer
			var assetPath bytes.Buffer
			assetRootPathTmpl, err := template.New("assetRootPath").Parse(adminSetting.AssetRootPath)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			err = assetRootPathTmpl.Execute(&assetRootPath, i)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			AssetTypePathTmpl, err := template.New("assetTypePath").Parse(adminSetting.AssetTypePath)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			err = AssetTypePathTmpl.Execute(&assetTypePath, i)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			assetPathTmpl, err := template.New("assetPath").Parse(adminSetting.AssetPath)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			err = assetPathTmpl.Execute(&assetPath, i)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			// Umask를 셋팅한다.
			if adminSetting.Umask == "" {
				unix.Umask(0)
			} else {
				umask, err := strconv.Atoi(adminSetting.Umask)
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
				unix.Umask(umask)
			}
			// ShotRootPath를 점검하고 없다면 경로를 생성한다.
			if _, err := os.Stat(assetRootPath.String()); os.IsNotExist(err) {
				per, err := strconv.ParseInt(adminSetting.AssetRootPathPermission, 8, 64)
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
				err = os.MkdirAll(assetRootPath.String(), os.FileMode(per))
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
				// admin 셋팅에 UID와 GID가 선언되어 있다면, 설정을 해줍니다.
				if adminSetting.AssetRootPathUID != "" && adminSetting.AssetRootPathGID != "" {
					uid, err := strconv.Atoi(adminSetting.AssetRootPathUID)
					if err != nil {
						a.Error = err.Error()
						fails = append(fails, a)
						continue
					}
					gid, err := strconv.Atoi(adminSetting.AssetRootPathGID)
					if err != nil {
						a.Error = err.Error()
						fails = append(fails, a)
						continue
					}
					err = os.Chown(assetRootPath.String(), uid, gid)
					if err != nil {
						a.Error = err.Error()
						fails = append(fails, a)
						continue
					}
				}
			}
			// 개별 샷 경로를 생성한다.
			per, err := strconv.ParseInt(adminSetting.AssetPathPermission, 8, 64)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			err = os.MkdirAll(assetPath.String(), os.FileMode(per))
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			// admin 셋팅에 UID와 GID가 선언되어 있다면, 설정을 해줍니다.
			if adminSetting.AssetPathUID != "" && adminSetting.AssetPathGID != "" {
				uid, err := strconv.Atoi(adminSetting.AssetPathUID)
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
				gid, err := strconv.Atoi(adminSetting.AssetPathGID)
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
				err = os.Chown(assetPath.String(), uid, gid)
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
			}
			// AssetType 경로의 Permission을 설정한다.
			if adminSetting.AssetTypePathPermission != "" {
				per, err := strconv.ParseInt(adminSetting.AssetTypePathPermission, 8, 64)
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
				err = os.Chmod(assetTypePath.String(), os.FileMode(per))
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
			}
			// Seq 경로의 UID, GID를 설정한다.
			if adminSetting.AssetTypePathUID != "" && adminSetting.AssetTypePathGID != "" {
				uid, err := strconv.Atoi(adminSetting.AssetTypePathUID)
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
				gid, err := strconv.Atoi(adminSetting.AssetTypePathGID)
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
				err = os.Chown(assetTypePath.String(), uid, gid)
				if err != nil {
					a.Error = err.Error()
					fails = append(fails, a)
					continue
				}
			}
		}
		success = append(success, a)
	}

	type recipe struct {
		Success []Asset
		Fails   []Asset
		User    User
		Devmode bool
		SearchOption
	}
	rcp := recipe{}
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Devmode = *flagDevmode
	rcp.Success = success
	rcp.Fails = fails
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u

	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addAsset_success", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleEditItemSubmitv2 함수는 Item의 수정사항을 처리하는 페이지이다. legacy
func handleEditItemSubmitv2(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	//var logstring string
	//기존 Item의 값을 가지고 온다.
	project := r.FormValue("project")
	id := r.FormValue("id")
	file, fileHandler, fileErr := r.FormFile("Thumbnail")
	if fileErr == nil {
		if !(fileHandler.Header.Get("Content-Type") == "image/jpeg" || fileHandler.Header.Get("Content-Type") == "image/png") {
			http.Error(w, "업로드 파일이 jpeg 또는 png 파일이 아닙니다", http.StatusInternalServerError)
			return
		}
		//파일이 없다면 fileErr 값은 "http: no such file" 값이 된다.
		// 썸네일 파일이 존재한다면 아래 프로세스를 거친다.
		mediatype, fileParams, err := mime.ParseMediaType(fileHandler.Header.Get("Content-Disposition"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if *flagDebug {
			fmt.Println(mediatype)
			fmt.Println(fileParams)
			fmt.Println(fileHandler.Header.Get("Content-Type"))
			fmt.Println()
		}
		tempPath := os.TempDir() + fileHandler.Filename
		tempFile, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 사용자가 업로드한 파일을 tempFile에 복사한다.
		io.Copy(tempFile, io.LimitReader(file, MaxFileSize))
		tempFile.Close()
		defer os.Remove(tempPath)
		//fmt.Println(tempPath)
		thumbnailPath := fmt.Sprintf("%s/%s/%s.jpg", *flagThumbPath, project, id)
		thumbnailDir := filepath.Dir(thumbnailPath)
		// 썸네일을 생성할 경로가 존재하지 않는다면 생성한다.
		_, err = os.Stat(thumbnailDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(thumbnailDir, 0775)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// 디지털아이디어의 경우 스캔시스템에서 수동으로 이미지를 폴더에 생성하는 경우가 있다.
			if *flagCompany == "digitalidea" {
				err = dipath.Ideapath(thumbnailDir)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
		// 이미지변환
		src, err := imaging.Open(tempPath)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Resize the cropped image to width = 200px preserving the aspect ratio.
		dst := imaging.Fill(src, 410, 222, imaging.Center, imaging.Lanczos)
		err = imaging.Save(dst, thumbnailPath)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 디지털아이디어의 경우 스캔시스템에서 수동으로 이미지를 수정하는 경우도 있다.
		if *flagCompany == "digitalidea" {
			err = dipath.Ideapath(thumbnailPath)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	http.Redirect(w, r, "/editeditem", http.StatusSeeOther)
}
