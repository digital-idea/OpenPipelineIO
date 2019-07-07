package main

import (
	"fmt"
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
	"gopkg.in/mgo.v2"
)

// handleSearch 함수는 검색결과를 반환하는 페이지다.
func handleSearch(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		ID          string
		Projectlist []string
		Project     Project
		Searchop    SearchOption
		Infobarnum  Infobarnum
		Searchnum   Infobarnum
		Items       []Item
		Ddline3d    []string
		Ddline2d    []string
		Tags        []string
		Assettags   []string
		Dilog       string
		Wfs         string
		MailDNS     string
	}
	rcp := recipe{}
	// 쿠키값을 rcp로 보낸다.
	ssid, err = GetSessionID(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.ID = ssid.ID
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Dilog = *flagDILOG
	rcp.Wfs = *flagWFS
	rcp.MailDNS = *flagMailDNS

	//옵션불러오기.
	rcp.Searchop.Project = r.FormValue("Project")
	rcp.Searchop.Searchword = r.FormValue("Searchword")
	rcp.Searchop.Assign = str2bool(r.FormValue("Assign"))
	rcp.Searchop.Ready = str2bool(r.FormValue("Ready"))
	rcp.Searchop.Wip = str2bool(r.FormValue("Wip"))
	rcp.Searchop.Confirm = str2bool(r.FormValue("Confirm"))
	rcp.Searchop.Done = str2bool(r.FormValue("Done"))
	rcp.Searchop.Omit = str2bool(r.FormValue("Omit"))
	rcp.Searchop.Hold = str2bool(r.FormValue("Hold"))
	rcp.Searchop.Out = str2bool(r.FormValue("Out"))
	rcp.Searchop.None = str2bool(r.FormValue("None"))
	rcp.Searchop.Sortkey = r.FormValue("Sortkey")
	rcp.Searchop.Template = r.FormValue("Template")

	// 마지막으로 검색한 프로젝트를 쿠키에 저장한다.
	// 이 정보는 index 페이지 접근시 프로젝트명으로 사용된다.
	cookie := http.Cookie{
		Name:   "Project",
		Value:  rcp.Searchop.Project,
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)

	if rcp.Searchop.Project != "" && rcp.Searchop.Searchword != "" {
		rcp.Items, err = Search(session, rcp.Searchop)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	rcp.Ddline3d, err = DistinctDdline(session, rcp.Searchop.Project, "ddline3d")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Ddline2d, err = DistinctDdline(session, rcp.Searchop.Project, "ddline2d")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Tags, err = Distinct(session, rcp.Searchop.Project, "tag")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Assettags, err = Distinct(session, rcp.Searchop.Project, "assettags")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Infobarnum, err = Resultnum(session, rcp.Searchop.Project)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Searchnum, err = Searchnum(rcp.Searchop.Project, rcp.Items)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Project, err = getProject(session, r.FormValue("Project"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(rcp.Items) == 0 {
		err = t.ExecuteTemplate(w, "searchno", rcp)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	err = t.ExecuteTemplate(w, rcp.Searchop.Template, rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAssettags는 에셋태그 페이지이다.
func handleAssettags(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		ID          string
		Projectlist []string
		Project     Project
		Searchop    SearchOption
		Infobarnum  Infobarnum
		Searchnum   Infobarnum
		Items       []Item
		Ddline3d    []string
		Ddline2d    []string
		Tags        []string
		Assettags   []string
		Dilog       string
		Wfs         string
		MailDNS     string
	}
	rcp := recipe{}
	rcp.Dilog = *flagDILOG
	rcp.Wfs = *flagWFS
	rcp.MailDNS = *flagMailDNS
	rcp.Searchop.Template = "csi3"
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 옵션불러오기 및 초기화.
	rcp.Project, err = getProject(session, strings.Split(r.URL.Path, "/")[2])
	mode := strings.Split(r.URL.Path, "/")[1] // assettags 또는 assettree 문자열이다.

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.Searchop.Project = strings.Split(r.URL.Path, "/")[2]
	rcp.Searchop.Searchword = strings.Split(r.URL.Path, "/")[3]
	rcp.Searchop.setStatusDefault()
	rcp.Searchop.Sortkey = "slug"
	rcp.Searchop.Template = "csi3"

	if rcp.Searchop.Project != "" && rcp.Searchop.Searchword != "" {
		rcp.Searchop.setStatusAll() // 에셋태그 검색시 전체상태를 검색한다.
		if mode == "assettree" {
			rcp.Items, err = SearchAssetTree(session, rcp.Searchop)
		} else {
			rcp.Items, err = SearchAssettags(session, rcp.Searchop)
		}
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Ddline3d, err = DistinctDdline(session, rcp.Searchop.Project, "ddline3d")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Ddline2d, err = DistinctDdline(session, rcp.Searchop.Project, "ddline2d")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Tags, err = Distinct(session, rcp.Searchop.Project, "tag")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Assettags, err = Distinct(session, rcp.Searchop.Project, "assettags")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Infobarnum, err = Resultnum(session, rcp.Searchop.Project)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Searchnum, err = Searchnum(rcp.Searchop.Project, rcp.Items)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	rcp.Searchop.Searchword = "#" + rcp.Searchop.Searchword // 태그검색은 #으로 시작한다.
	rcp.Searchop.setStatusDefault()                         // 검색이후 상태를 기본형으로 바꾸어 놓는다.
	err = t.ExecuteTemplate(w, rcp.Searchop.Template, rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleEdit 함수는 Item 편집페이지이다.
func handleEdit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	}
	rcp := recipe{}
	defer session.Close()
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
	err = t.ExecuteTemplate(w, "editItem", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleTags 함수는 태그 클릭시 출력되는 페이지이다.
func handleTags(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		ID          string
		Projectlist []string
		Project     Project
		Searchop    SearchOption
		Infobarnum  Infobarnum
		Searchnum   Infobarnum
		Items       []Item
		Ddline3d    []string
		Ddline2d    []string
		Tags        []string
		Assettags   []string
		Dilog       string
		Wfs         string
		MailDNS     string
	}
	rcp := recipe{}
	rcp.Dilog = *flagDILOG
	rcp.Wfs = *flagWFS
	rcp.MailDNS = *flagMailDNS
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 옵션불러오기 및 초기화.
	rcp.Project, err = getProject(session, strings.Split(r.URL.Path, "/")[2])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Searchop.Project = strings.Split(r.URL.Path, "/")[2]
	rcp.Searchop.Searchword = strings.Split(r.URL.Path, "/")[3]
	rcp.Searchop.setStatusDefault()
	rcp.Searchop.Sortkey = "slug"
	rcp.Searchop.Template = "csi3"
	if rcp.Searchop.Project != "" && rcp.Searchop.Searchword != "" {
		rcp.Searchop.setStatusAll() //태그검색은 전체상태를 검색한다.
		rcp.Items, err = SearchTag(session, rcp.Searchop)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Ddline3d, err = DistinctDdline(session, rcp.Searchop.Project, "ddline3d")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Ddline2d, err = DistinctDdline(session, rcp.Searchop.Project, "ddline2d")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Tags, err = Distinct(session, rcp.Searchop.Project, "tag")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Assettags, err = Distinct(session, rcp.Searchop.Project, "assettags")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Infobarnum, err = Resultnum(session, rcp.Searchop.Project)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Searchnum, err = Searchnum(rcp.Searchop.Project, rcp.Items)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 태그는 전체 검색을 하고 #으로 시작한다.
	// Search 박스의 상태 옵션은 기본형이어야 한다
	rcp.Searchop.Searchword = "#" + rcp.Searchop.Searchword
	rcp.Searchop.setStatusDefault()
	err = t.ExecuteTemplate(w, rcp.Searchop.Template, rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleEditItemSubmit 함수는 Item의 수정사항을 처리하는 페이지이다.
func handleEditItemSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	//var logstring string
	//기존 Item의 값을 가지고 온다.
	project := r.FormValue("project")
	slug := r.FormValue("slug")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	CurrentItem, err := getItem(session, project, slug)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//과거의 값과 현재값을 비교하고 다르면 셋팅한다.
	NewItem := CurrentItem //과거값을 먼저 복사한다.
	NewItem.Rollmedia = r.FormValue("Rollmedia")
	NewItem.Platesize = r.FormValue("Platesize")
	NewItem.Dsize = r.FormValue("Dsize")
	NewItem.Rendersize = r.FormValue("Rendersize")
	NewItem.Thummov = r.FormValue("thummov")
	NewItem.Beforemov = r.FormValue("Beforemov")
	NewItem.Aftermov = r.FormValue("Aftermov")
	NewItem.Retimeplate = r.FormValue("retimeplate")
	if CurrentItem.Type == "org" || CurrentItem.Type == "left" { // 일반상황과 입체상황을 체크한다.
		NewItem.Rnum = strings.ToUpper(r.FormValue("Rnum"))
		NewItem.OCIOcc = r.FormValue("OCIOcc")
		NewItem.PlateIn, err = strconv.Atoi(r.FormValue("PlateIn"))
		if err != nil {
			NewItem.PlateIn = CurrentItem.PlateIn // 에러가 나면 과거값으로 놔둔다.
		}
		NewItem.PlateOut, err = strconv.Atoi(r.FormValue("PlateOut"))
		if err != nil {
			NewItem.PlateOut = CurrentItem.PlateOut // 에러가 나면 과거값으로 놔둔다.
		}
		NewItem.JustIn, err = strconv.Atoi(r.FormValue("JustIn"))
		if err != nil {
			NewItem.JustIn = CurrentItem.JustIn // 에러가 나면 과거값으로 놔둔다.
		}
		NewItem.JustOut, err = strconv.Atoi(r.FormValue("JustOut"))
		if err != nil {
			NewItem.JustOut = CurrentItem.JustOut // 에러가 나면 과거값으로 놔둔다.
		}
	}
	NewItem.ObjectidIn, err = strconv.Atoi(r.FormValue("ObjectidIn"))
	if err != nil {
		NewItem.ObjectidIn = CurrentItem.ObjectidIn // 에러가 나면 과거값으로 놔둔다.
	}
	NewItem.ObjectidOut, err = strconv.Atoi(r.FormValue("ObjectidOut"))
	if err != nil {
		NewItem.ObjectidOut = CurrentItem.ObjectidOut // 에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Ddline2d = ToFullTime(r.FormValue("Ddline2d"))
	NewItem.Ddline3d = ToFullTime(r.FormValue("Ddline3d"))
	NewItem.Tag = Str2Tags(r.FormValue("Tag"))
	NewItem.Assettags = Str2Tags(r.FormValue("Assettags"))
	// 카메라 퍼블리쉬 관련 셋팅
	NewItem.ProductionCam.PubTask = r.FormValue("ProductionCamPubTask")
	NewItem.ProductionCam.PubPath = r.FormValue("ProductionCamPubPath")
	NewItem.ProductionCam.Projection = str2bool(r.FormValue("ProductionCamProjection"))
	//concept
	NewItem.Concept.Status = r.FormValue("ConceptStatus")
	NewItem.Concept.User = r.FormValue("ConceptUser")
	NewItem.Concept.UserNote = r.FormValue("ConceptUserNote")
	NewItem.Concept.Startdate = ToFullTime(r.FormValue("ConceptStartdate"))
	NewItem.Concept.Due, err = strconv.Atoi(r.FormValue("ConceptDue"))
	if err != nil {
		NewItem.Concept.Due = CurrentItem.Concept.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Concept.Predate = ToFullTime(r.FormValue("ConceptPredate"))
	NewItem.Concept.Date = ToFullTime(r.FormValue("ConceptDate"))
	if *flagCompany == "digitalidea" {
		// Concept Team은 윈도우즈를 사용한다.
		NewItem.Concept.Mov = dipath.Win2lin(r.FormValue("ConceptMov"))
	} else {
		NewItem.Concept.Mov = r.FormValue("ConceptMov")
	}

	//model
	NewItem.Model.Status = r.FormValue("ModelStatus")
	NewItem.Model.User = r.FormValue("ModelUser")
	NewItem.Model.UserNote = r.FormValue("ModelUserNote")
	NewItem.Model.Startdate = ToFullTime(r.FormValue("ModelStartdate"))
	NewItem.Model.Due, err = strconv.Atoi(r.FormValue("ModelDue"))
	if err != nil {
		NewItem.Model.Due = CurrentItem.Model.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Model.Predate = ToFullTime(r.FormValue("ModelPredate"))
	NewItem.Model.Date = ToFullTime(r.FormValue("ModelDate"))
	if *flagCompany == "digitalidea" {
		// 모델링팀 일부는 ZBrush로 인해 윈도우즈를 사용한다.
		NewItem.Model.Mov = dipath.Win2lin(r.FormValue("ModelMov"))
	} else {
		NewItem.Model.Mov = r.FormValue("ModelMov")
	}
	//mm
	NewItem.Mm.Status = r.FormValue("MmStatus")
	NewItem.Mm.User = r.FormValue("MmUser")
	NewItem.Mm.UserNote = r.FormValue("MmUserNote")
	NewItem.Mm.Startdate = ToFullTime(r.FormValue("MmStartdate"))
	NewItem.Mm.Due, err = strconv.Atoi(r.FormValue("MmDue"))
	if err != nil {
		NewItem.Mm.Due = CurrentItem.Mm.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Mm.Predate = ToFullTime(r.FormValue("MmPredate"))
	NewItem.Mm.Date = ToFullTime(r.FormValue("MmDate"))
	NewItem.Mm.Mov = r.FormValue("MmMov")
	//layout
	NewItem.Layout.Status = r.FormValue("LayoutStatus")
	NewItem.Layout.User = r.FormValue("LayoutUser")
	NewItem.Layout.UserNote = r.FormValue("LayoutUserNote")
	NewItem.Layout.Startdate = ToFullTime(r.FormValue("LayoutStartdate"))
	NewItem.Layout.Due, err = strconv.Atoi(r.FormValue("LayoutDue"))
	if err != nil {
		NewItem.Layout.Due = CurrentItem.Layout.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Layout.Predate = ToFullTime(r.FormValue("LayoutPredate"))
	NewItem.Layout.Date = ToFullTime(r.FormValue("LayoutDate"))
	NewItem.Layout.Mov = r.FormValue("LayoutMov")
	//ani
	NewItem.Ani.Status = r.FormValue("AniStatus")
	NewItem.Ani.User = r.FormValue("AniUser")
	NewItem.Ani.UserNote = r.FormValue("AniUserNote")
	NewItem.Ani.Startdate = ToFullTime(r.FormValue("AniStartdate"))
	NewItem.Ani.Due, err = strconv.Atoi(r.FormValue("AniDue"))
	if err != nil {
		NewItem.Ani.Due = CurrentItem.Ani.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Ani.Predate = ToFullTime(r.FormValue("AniPredate"))
	NewItem.Ani.Date = ToFullTime(r.FormValue("AniDate"))
	NewItem.Ani.Mov = r.FormValue("AniMov")
	//fx
	NewItem.Fx.Status = r.FormValue("FxStatus")
	NewItem.Fx.User = r.FormValue("FxUser")
	NewItem.Fx.UserNote = r.FormValue("FxUserNote")
	NewItem.Fx.Startdate = ToFullTime(r.FormValue("FxStartdate"))
	NewItem.Fx.Due, err = strconv.Atoi(r.FormValue("FxDue"))
	if err != nil {
		NewItem.Fx.Due = CurrentItem.Fx.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Fx.Predate = ToFullTime(r.FormValue("FxPredate"))
	NewItem.Fx.Date = ToFullTime(r.FormValue("FxDate"))
	NewItem.Fx.Mov = r.FormValue("FxMov")
	//mg
	NewItem.Mg.Status = r.FormValue("MgStatus")
	NewItem.Mg.User = r.FormValue("MgUser")
	NewItem.Mg.UserNote = r.FormValue("MgUserNote")
	NewItem.Mg.Startdate = ToFullTime(r.FormValue("MgStartdate"))
	NewItem.Mg.Due, err = strconv.Atoi(r.FormValue("MgDue"))
	if err != nil {
		NewItem.Mg.Due = CurrentItem.Mg.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Mg.Predate = ToFullTime(r.FormValue("MgPredate"))
	NewItem.Mg.Date = ToFullTime(r.FormValue("MgDate"))
	if *flagCompany == "digitalidea" {
		// 모션그래픽팀은 에프터이펙트로 인해 윈도우즈를 사용한다.
		NewItem.Mg.Mov = dipath.Win2lin(r.FormValue("MgMov"))
	} else {
		NewItem.Mg.Mov = r.FormValue("MgMov")
	}

	//temp1
	NewItem.Temp1.Title = r.FormValue("Temp1Title")
	NewItem.Temp1.Status = r.FormValue("Temp1Status")
	NewItem.Temp1.User = r.FormValue("Temp1User")
	NewItem.Temp1.UserNote = r.FormValue("Temp1UserNote")
	NewItem.Temp1.Startdate = ToFullTime(r.FormValue("Temp1Startdate"))
	NewItem.Temp1.Due, err = strconv.Atoi(r.FormValue("Temp1Due"))
	if err != nil {
		NewItem.Temp1.Due = CurrentItem.Temp1.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Temp1.Predate = ToFullTime(r.FormValue("Temp1Predate"))
	NewItem.Temp1.Date = ToFullTime(r.FormValue("Temp1Date"))
	NewItem.Temp1.Mov = r.FormValue("Temp1Mov")
	//previz
	NewItem.Previz.Status = r.FormValue("PrevizStatus")
	NewItem.Previz.User = r.FormValue("PrevizUser")
	NewItem.Previz.UserNote = r.FormValue("PrevizUserNote")
	NewItem.Previz.Startdate = ToFullTime(r.FormValue("PrevizStartdate"))
	NewItem.Previz.Due, err = strconv.Atoi(r.FormValue("PrevizDue"))
	if err != nil {
		NewItem.Previz.Due = CurrentItem.Previz.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Previz.Predate = ToFullTime(r.FormValue("PrevizPredate"))
	NewItem.Previz.Date = ToFullTime(r.FormValue("PrevizDate"))
	NewItem.Previz.Mov = r.FormValue("PrevizMov")
	//fur
	NewItem.Fur.Status = r.FormValue("FurStatus")
	NewItem.Fur.User = r.FormValue("FurUser")
	NewItem.Fur.UserNote = r.FormValue("FurUserNote")
	NewItem.Fur.Startdate = ToFullTime(r.FormValue("FurStartdate"))
	NewItem.Fur.Due, err = strconv.Atoi(r.FormValue("FurDue"))
	if err != nil {
		NewItem.Fur.Due = CurrentItem.Fur.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Fur.Predate = ToFullTime(r.FormValue("FurPredate"))
	NewItem.Fur.Date = ToFullTime(r.FormValue("FurDate"))
	NewItem.Fur.Mov = r.FormValue("FurMov")
	//sim
	NewItem.Sim.Status = r.FormValue("SimStatus")
	NewItem.Sim.User = r.FormValue("SimUser")
	NewItem.Sim.UserNote = r.FormValue("SimUserNote")
	NewItem.Sim.Startdate = ToFullTime(r.FormValue("SimStartdate"))
	NewItem.Sim.Due, err = strconv.Atoi(r.FormValue("SimDue"))
	if err != nil {
		NewItem.Sim.Due = CurrentItem.Sim.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Sim.Predate = ToFullTime(r.FormValue("SimPredate"))
	NewItem.Sim.Date = ToFullTime(r.FormValue("SimDate"))
	NewItem.Sim.Mov = r.FormValue("SimMov")
	//crowd
	NewItem.Crowd.Status = r.FormValue("CrowdStatus")
	NewItem.Crowd.User = r.FormValue("CrowdUser")
	NewItem.Crowd.UserNote = r.FormValue("CrowdUserNote")
	NewItem.Crowd.Startdate = ToFullTime(r.FormValue("CrowdStartdate"))
	NewItem.Crowd.Due, err = strconv.Atoi(r.FormValue("CrowdDue"))
	if err != nil {
		NewItem.Crowd.Due = CurrentItem.Crowd.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Crowd.Predate = ToFullTime(r.FormValue("CrowdPredate"))
	NewItem.Crowd.Date = ToFullTime(r.FormValue("CrowdDate"))
	NewItem.Crowd.Mov = r.FormValue("CrowdMov")
	//light
	NewItem.Light.Status = r.FormValue("LightStatus")
	NewItem.Light.User = r.FormValue("LightUser")
	NewItem.Light.UserNote = r.FormValue("LightUserNote")
	NewItem.Light.Startdate = ToFullTime(r.FormValue("LightStartdate"))
	NewItem.Light.Due, err = strconv.Atoi(r.FormValue("LightDue"))
	if err != nil {
		NewItem.Light.Due = CurrentItem.Light.Due //에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Light.Predate = ToFullTime(r.FormValue("LightPredate"))
	NewItem.Light.Date = ToFullTime(r.FormValue("LightDate"))
	NewItem.Light.Mov = r.FormValue("LightMov")
	//comp
	NewItem.Comp.Status = r.FormValue("CompStatus")
	NewItem.Comp.User = r.FormValue("CompUser")
	NewItem.Comp.UserNote = r.FormValue("CompUserNote")
	NewItem.Comp.Startdate = ToFullTime(r.FormValue("CompStartdate"))
	NewItem.Comp.Due, err = strconv.Atoi(r.FormValue("CompDue"))
	if err != nil {
		NewItem.Comp.Due = CurrentItem.Comp.Due // 에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Comp.Predate = ToFullTime(r.FormValue("CompPredate"))
	NewItem.Comp.Date = ToFullTime(r.FormValue("CompDate"))
	NewItem.Comp.Mov = r.FormValue("CompMov")
	//matte
	NewItem.Matte.Status = r.FormValue("MatteStatus")
	NewItem.Matte.User = r.FormValue("MatteUser")
	NewItem.Matte.UserNote = r.FormValue("MatteUserNote")
	NewItem.Matte.Startdate = ToFullTime(r.FormValue("MatteStartdate"))
	NewItem.Matte.Due, err = strconv.Atoi(r.FormValue("MatteDue"))
	if err != nil {
		NewItem.Matte.Due = CurrentItem.Matte.Due // 에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Matte.Predate = ToFullTime(r.FormValue("MattePredate"))
	NewItem.Matte.Date = ToFullTime(r.FormValue("MatteDate"))
	if *flagCompany == "digitalidea" {
		// 매트팀은 포토샵등의 툴로 인해 윈도우즈를 사용한다.
		NewItem.Matte.Mov = dipath.Win2lin(r.FormValue("MatteMov"))
	} else {
		NewItem.Matte.Mov = r.FormValue("MatteMov")
	}
	//env
	NewItem.Env.Status = r.FormValue("EnvStatus")
	NewItem.Env.User = r.FormValue("EnvUser")
	NewItem.Env.UserNote = r.FormValue("EnvUserNote")
	NewItem.Env.Startdate = ToFullTime(r.FormValue("EnvStartdate"))
	NewItem.Env.Due, err = strconv.Atoi(r.FormValue("EnvDue"))
	if err != nil {
		NewItem.Env.Due = CurrentItem.Env.Due // 에러가 나면 과거값으로 놔둔다.
	}
	NewItem.Env.Predate = ToFullTime(r.FormValue("EnvPredate"))
	NewItem.Env.Date = ToFullTime(r.FormValue("EnvDate"))
	if *flagCompany == "digitalidea" {
		// 환경팀은 간혹 윈도우즈를 사용한다.
		NewItem.Env.Mov = dipath.Win2lin(r.FormValue("EnvMov"))
	} else {
		NewItem.Env.Mov = r.FormValue("EnvMov")
	}

	file, fileHandler, fileErr := r.FormFile("Thumbnail")
	if fileErr == nil {
		if !(fileHandler.Header.Get("Content-Type") == "image/jpeg" || fileHandler.Header.Get("Content-Type") == "image/png") {
			log.Println("업로드 파일이 jpeg 또는 png 파일이 아닙니다.")
			err = t.ExecuteTemplate(w, "edited", nil)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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
		thumbnailPath := fmt.Sprintf("%s/%s/%s.jpg", *flagThumbPath, project, NewItem.Slug)
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

	//DB업데이트
	NewItem.updateMdate(&CurrentItem) //팀의 mov를 비교하고 달라졌다면 mov 업데이트 날짜를 변경한다.
	err = setItem(session, project, NewItem)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//로그를 추후 처리한다.
	// logging(CurrentItem, NewItem)
	err = t.ExecuteTemplate(w, "edited", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleDdline 함수는 데드라인 클릭시 출력되는 페이지이다.
func handleDdline(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		ID          string
		Projectlist []string
		Project     Project
		Searchop    SearchOption
		Infobarnum  Infobarnum
		Searchnum   Infobarnum
		Items       []Item
		Ddline3d    []string
		Ddline2d    []string
		Tags        []string
		Assettags   []string
		Dilog       string
		Wfs         string
		MailDNS     string
	}
	rcp := recipe{}
	rcp.Dilog = *flagDILOG
	rcp.Wfs = *flagWFS
	rcp.MailDNS = *flagMailDNS
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//옵션불러오기.
	rcp.Project, err = getProject(session, strings.Split(r.URL.Path, "/")[3]) //project
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	part := strings.Split(r.URL.Path, "/")[2]                   //2d,3d,_
	rcp.Searchop.Project = strings.Split(r.URL.Path, "/")[3]    //project
	rcp.Searchop.Searchword = strings.Split(r.URL.Path, "/")[4] //date
	rcp.Searchop.setStatusDefault()
	rcp.Searchop.Sortkey = "slug"
	rcp.Searchop.Template = "csi3"
	if rcp.Searchop.Project != "" && rcp.Searchop.Searchword != "" {
		rcp.Searchop.setStatusAll() // 데드라인 클릭시 전체 상태를 검색한다.
		rcp.Items, err = SearchDdline(session, rcp.Searchop, part)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Ddline3d, err = DistinctDdline(session, rcp.Searchop.Project, "ddline3d")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Ddline2d, err = DistinctDdline(session, rcp.Searchop.Project, "ddline2d")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Tags, err = Distinct(session, rcp.Searchop.Project, "tag")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Assettags, err = Distinct(session, rcp.Searchop.Project, "assettags")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Infobarnum, err = Resultnum(session, rcp.Searchop.Project)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Searchnum, err = Searchnum(rcp.Searchop.Project, rcp.Items)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	rcp.Searchop.setStatusDefault() // 페이지 렌더링시 상태는 기본형으로 바꾼다.
	err = t.ExecuteTemplate(w, rcp.Searchop.Template, rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleIndex 함수는 index 페이지이다.
func handleIndex(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		ID          string
		Projectlist []string
		Searchop    SearchOption
	}
	rcp := recipe{}
	rcp.ID = ssid.ID

	// 쿠키에 저장된 값이 있다면 rcp에 저장한다.
	for _, cookie := range r.Cookies() {
		if cookie.Name == "Project" {
			rcp.Searchop.Project = cookie.Value
		}
	}
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//검색바 초기셋팅
	rcp.Searchop.Assign = true
	rcp.Searchop.Ready = true
	rcp.Searchop.Wip = true
	rcp.Searchop.Confirm = true
	rcp.Searchop.Sortkey = "slug"
	rcp.Searchop.Template = "csi3"
	err = t.ExecuteTemplate(w, "index", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	}
	rcp := recipe{}
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
	err = t.ExecuteTemplate(w, "addShot", rcp)
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
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_'
	}
	names := strings.FieldsFunc(name, f)
	type Shot struct {
		Name  string
		Error string
	}
	var success []Shot
	var fails []Shot
	for _, n := range names {
		if n == " " || n == "" { // 사용자가 실수로 여러개의 스페이스를 추가할 수 있다.
			continue
		}
		s := Shot{}
		s.Name = n
		if !regexpShotname.MatchString(n) {
			s.Error = "SS_0010 형식의 이름이 아닙니다"
			fails = append(fails, s)
			continue
		}
		i := Item{}
		i.Name = n
		if stereo == "true" {
			i.Type = "left"
		} else {
			i.Type = "org"
		}
		i.Project = project
		i.Slug = i.Name + "_" + i.Type
		i.ID = i.Name + "_" + i.Type
		i.Status = ASSIGN
		i.Thumpath = fmt.Sprintf("/%s/%s_%s.jpg", i.Project, i.Name, i.Type)
		i.Platepath = fmt.Sprintf("/show/%s/seq/%s/%s/plate/", i.Project, strings.Split(i.Name, "_")[0], i.Name)
		i.Thummov = fmt.Sprintf("/show/%s/seq/%s/%s/plate/%s_%s.mov", i.Project, strings.Split(i.Name, "_")[0], i.Name, i.Name, i.Type)
		i.Scantime = time.Now().Format(time.RFC3339)
		i.Updatetime = time.Now().Format(time.RFC3339)
		err = addItem(session, project, i)
		if err != nil {
			s.Error = err.Error()
			fails = append(fails, s)
			continue
		}
		success = append(success, s)
	}

	type recipe struct {
		Success []Shot
		Fails   []Shot
		User    User
	}
	rcp := recipe{}
	rcp.Success = success
	rcp.Fails = fails
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u

	w.Header().Set("Content-Type", "text/html")
	err = t.ExecuteTemplate(w, "addShot_success", rcp)
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
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	}
	rcp := recipe{}
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
	err = t.ExecuteTemplate(w, "addAsset", rcp)
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
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_'
	}
	names := strings.FieldsFunc(name, f)
	type Asset struct {
		Name  string
		Error string
	}
	var success []Asset
	var fails []Asset
	for _, n := range names {
		if n == " " || n == "" { // 사용자가 실수로 여러개의 스페이스를 추가할 수 있다.
			continue
		}
		a := Asset{}
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
		i.Slug = i.Name + "_" + i.Type
		i.ID = i.Name + "_" + i.Type
		if assettype == "comp" {
			i.Comp.Status = ASSIGN
		} else if assettype == "env" {
			i.Env.Status = ASSIGN
			i.Matte.Status = ASSIGN
		} else {
			i.Model.Status = ASSIGN
		}
		i.Status = ASSIGN
		i.Updatetime = time.Now().Format(time.RFC3339)
		i.Assettype = assettype
		i.Assettags = []string{assettype, construction}
		err = addItem(session, project, i)
		if err != nil {
			a.Error = err.Error()
			fails = append(fails, a)
			continue
		}
		success = append(success, a)
	}

	type recipe struct {
		Success []Asset
		Fails   []Asset
		User    User
	}
	rcp := recipe{}
	rcp.Success = success
	rcp.Fails = fails
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u

	w.Header().Set("Content-Type", "text/html")
	err = t.ExecuteTemplate(w, "addAsset_success", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
