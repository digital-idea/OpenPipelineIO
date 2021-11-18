package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/digital-idea/ditime"
	"gopkg.in/mgo.v2"
)

// handleAddProject 함수는 프로젝트를 추가하는 페이지이다.
func handleAddProject(w http.ResponseWriter, r *http.Request) {
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
	type recipe struct {
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
	u, err := getUser(session, ssid.ID)
	rcp.User = u

	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addProject", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddProjectSubmit 함수는 사용자로부터 프로젝트 id를 받아서 프로젝트를 생성한다.
func handleAddProjectSubmit(w http.ResponseWriter, r *http.Request) {
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
	id := r.FormValue("ID")
	p := *NewProject(id)
	err = addProject(session, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/projectinfo", http.StatusSeeOther)
}

// handleProjectinfo 함수는 프로젝트 자료구조 페이지이다.
func handleProjectinfo(w http.ResponseWriter, r *http.Request) {
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
	status := q.Get("status")
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		Projects []Project
		MailDNS  string
		User
		Devmode bool
		SearchOption
	}
	rcp := recipe{}
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	rcp.MailDNS = *flagMailDNS
	rcp.Devmode = *flagDevmode
	if status != "" {
		rcp.Projects, err = getStatusProjects(session, ToProjectStatus(status))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		rcp.Projects, err = getProjects(session)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 만약 사용자에게 AccessProjects가 설정되어있다면 해당프로젝트만 보여야한다.
	if len(rcp.User.AccessProjects) != 0 {
		var accessProjects []Project
		for _, i := range rcp.Projects {
			for _, j := range rcp.User.AccessProjects {
				if i.ID != j {
					continue
				}
				accessProjects = append(accessProjects, i)
			}
		}
		rcp.Projects = accessProjects
	}

	err = TEMPLATES.ExecuteTemplate(w, "projectinfo", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ToProjectStatus 함수는 문자를 받아서 ProjectStatus 형으로 변환합니다.
func ToProjectStatus(s string) ProjectStatus {
	switch s {
	case "pre", "ready", "준비":
		return PreProjectStatus
	case "post", "wip":
		return PostProjectStatus
	case "layover", "중단":
		return LayoverProjectStatus
	case "backup", "백업":
		return BackupProjectStatus
	case "archive", "done", "종료":
		return ArchiveProjectStatus
	case "lawsuit", "소송":
		return LawsuitProjectStatus
	default:
		return TestProjectStatus
	}
}

// handleEditProjectSubmit 함수는 Projectinfo의  수정정보를 처리하는 페이지이다.
func handleEditProjectSubmit(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	current, err := getProject(session, r.FormValue("Id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renewal := current //과거 프로젝트 값을 셋팅한다.
	if current.Name != r.FormValue("Name") {
		renewal.Name = r.FormValue("Name")
	}
	if current.NetflixShowID != r.FormValue("NetflixShowID") {
		renewal.NetflixShowID = r.FormValue("NetflixShowID")
	}
	if current.MailHead != r.FormValue("MailHead") {
		renewal.MailHead = r.FormValue("MailHead")
	}
	if current.Style != r.FormValue("Style") {
		renewal.Style = r.FormValue("Style")
	}
	if current.Stereo != str2bool(r.FormValue("Stereo")) {
		renewal.Stereo = str2bool(r.FormValue("Stereo"))
	}
	if current.FlexibleStatus != str2bool(r.FormValue("FlexibleStatus")) {
		renewal.FlexibleStatus = str2bool(r.FormValue("FlexibleStatus"))
	}
	if current.Screenx != str2bool(r.FormValue("Screenx")) {
		renewal.Screenx = str2bool(r.FormValue("Screenx"))
	}
	t, err := ditime.ToFullTime(19, r.FormValue("Deadline"))
	if err != nil {
		log.Println(err)
	}
	if current.Deadline != t {
		renewal.Deadline = t
	}
	if current.Director != r.FormValue("Director") {
		renewal.Director = r.FormValue("Director")
	}
	if current.Super != r.FormValue("Super") {
		renewal.Super = r.FormValue("Super")
	}
	renewal.OnsetSuper = r.FormValue("OnsetSuper")
	renewal.CgSuper = r.FormValue("CgSuper")
	renewal.Pd = r.FormValue("Pd")
	renewal.Pm = r.FormValue("Pm")
	renewal.PmEmail = r.FormValue("PmEmail")
	renewal.Pa = r.FormValue("Pa")
	renewal.Edit = r.FormValue("Edit")
	renewal.EditContact = r.FormValue("EditContact")
	renewal.Di = r.FormValue("Di")
	renewal.DiContact = r.FormValue("DiContact")
	renewal.Sound = r.FormValue("Sound")
	renewal.SoundContact = r.FormValue("SoundContact")
	renewal.Message = r.FormValue("Message")
	renewal.Wiki = r.FormValue("Wiki")
	renewal.Daily = r.FormValue("Daily")
	renewal.EditDir = r.FormValue("EditDir")
	fps, err := strconv.ParseFloat(r.FormValue("Fps"), 64)
	if err == nil {
		renewal.Fps = fps
	}
	aspectratio, err := strconv.ParseFloat(r.FormValue("AspectRatio"), 64)
	if err == nil {
		renewal.AspectRatio = aspectratio
	}
	startframe, err := strconv.Atoi(r.FormValue("StartFrame"))
	if err == nil {
		renewal.StartFrame = startframe
	}
	versionnum, err := strconv.Atoi(r.FormValue("VersionNum"))
	if err == nil {
		renewal.VersionNum = versionnum
	}
	seqnum, err := strconv.Atoi(r.FormValue("SeqNum"))
	if err == nil {
		renewal.SeqNum = seqnum
	}
	renewal.Issue = r.FormValue("Issue")
	platewidth, err := strconv.Atoi(r.FormValue("PlateWidth"))
	if err == nil {
		renewal.PlateWidth = platewidth
	}
	plateheight, err := strconv.Atoi(r.FormValue("PlateHeight"))
	if err == nil {
		renewal.PlateHeight = plateheight
	}
	renewal.ResizeType = r.FormValue("ResizeType")
	renewal.PlateExt = r.FormValue("PlateExt")
	renewal.ExrCompression = r.FormValue("ExrCompression")
	renewal.Camera = r.FormValue("Camera")
	renewal.PlateInColorspace = r.FormValue("PlateInColorspace")
	renewal.PlateOutColorspace = r.FormValue("PlateOutColorspace")
	renewal.ProxyOutColorspace = r.FormValue("ProxyOutColorspace")
	renewal.PostProductionProxyCodec = r.FormValue("PostProductionProxyCodec")
	outputmovWidth, err := strconv.Atoi(r.FormValue("OutputMov.Width"))
	if err == nil {
		renewal.OutputMov.Width = outputmovWidth
	}
	outputmovHeight, err := strconv.Atoi(r.FormValue("OutputMov.Height"))
	if err == nil {
		renewal.OutputMov.Height = outputmovHeight
	}
	renewal.OutputMov.Codec = r.FormValue("OutputMov.Codec")
	outputmovFps, err := strconv.ParseFloat(r.FormValue("OutputMov.Fps"), 64)
	if err == nil {
		renewal.OutputMov.Fps = outputmovFps
	}
	renewal.OutputMov.InColorspace = r.FormValue("OutputMov.InColorspace")
	renewal.OutputMov.OutColorspace = r.FormValue("OutputMov.OutColorspace")
	editmovWidth, err := strconv.Atoi(r.FormValue("EditMov.Width"))
	if err == nil {
		renewal.EditMov.Width = editmovWidth
	}
	editmovHeight, err := strconv.Atoi(r.FormValue("EditMov.Height"))
	if err == nil {
		renewal.EditMov.Height = editmovHeight
	}
	renewal.EditMov.Codec = r.FormValue("EditMov.Codec")
	editmovFps, err := strconv.ParseFloat(r.FormValue("EditMov.Fps"), 64)
	if err == nil {
		renewal.EditMov.Fps = editmovFps
	}
	renewal.EditMov.InColorspace = r.FormValue("EditMov.InColorspace")
	renewal.EditMov.OutColorspace = r.FormValue("EditMov.OutColorspace")
	// 마일스톤 추가하기.
	status, err := strconv.Atoi(r.FormValue("Status"))
	if err == nil {
		renewal.Status = ProjectStatus(status)
	}
	renewal.Lut = r.FormValue("Lut")
	renewal.LutInColorspace = r.FormValue("LutInColorspace")
	renewal.LutOutColorspace = r.FormValue("LutOutColorspace")
	renewal.Description = r.FormValue("Description")
	renewal.NukeGizmo = r.FormValue("NukeGizmo")
	renewal.FxElement = r.FormValue("FxElement")
	renewal.MayaCropMaskSize = r.FormValue("MayaCropMaskSize")
	cropaspectratio, err := strconv.ParseFloat(r.FormValue("CropAspectRatio"), 64)
	if err == nil {
		renewal.CropAspectRatio = cropaspectratio
	}
	houdiniImportScale, err := strconv.ParseFloat(r.FormValue("HoudiniImportScale"), 64)
	if err == nil {
		renewal.HoudiniImportScale = houdiniImportScale
	}
	screenxOverlay, err := strconv.ParseFloat(r.FormValue("ScreenxOverlay"), 64)
	if err == nil {
		renewal.ScreenxOverlay = screenxOverlay
	}
	renewal.AWSS3 = r.FormValue("AWSS3")
	renewal.AWSProfile = r.FormValue("AWSProfile")
	renewal.AWSLocalpath = r.FormValue("AWSLocalpath")
	renewal.SlackWebhookURL = r.FormValue("SlackWebhookURL")
	// 새로 변경된 정보를 DB에 저장한다.
	err = setProject(session, renewal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/projectinfo", http.StatusSeeOther)
}

// handleEditProject 함수는 프로젝트 편집페이지이다.
func handleEditProject(w http.ResponseWriter, r *http.Request) {
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
	id := q.Get("id") // 프로젝트id에 사용할 것
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		Project            `json:"project"`
		User               `json:"user"`
		Devmode            bool `json:"devmode"`
		SearchOption       `json:"searchoption"`
		DefaultColorspaces []string `json:"defaultcolorspace"`
		OCIOColorspaces    []string `json:"ociocolorspaces"`
	}
	rcp := recipe{}
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Devmode = *flagDevmode
	p, err := getProject(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Project = p
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	rcp.OCIOColorspaces, err = loadOCIOConfig(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.DefaultColorspaces = []string{"default", "linear", "sRGB", "rec709", "Cineon", "AlexaV3LogC", "REDLog", "Gamma2.2", "ACEScg", "ACES2065-1"}
	err = TEMPLATES.ExecuteTemplate(w, "editProject", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return

}

// handleRmProject 함수는 project을 삭제하는 페이지이다.
func handleRmProject(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel != AdminAccessLevel { // Admin
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
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
		User        User
		Projectlist []string
		Devmode     bool
		SearchOption
		AdminSetting Setting
	}
	rcp := recipe{}
	rcp.AdminSetting = CachedAdminSetting
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "rmproject", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleRmProjectSubmit 함수는 project를 삭제한다.
func handleRmProjectSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel != AdminAccessLevel {
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
	rmReviews := str2bool(r.FormValue("rmreviews"))
	type recipe struct {
		User    User
		Devmode bool
		SearchOption
		Error   string
		Project string
	}
	rcp := recipe{}
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Devmode = *flagDevmode
	rcp.Project = project
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	// 리뷰데이터 삭제
	if rmReviews {
		// 1. 해당 프로젝트의 리뷰 데이터를 가지고 온다.
		reviews, err := searchReview(session, "project:"+rcp.Project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. 리뷰 데이터의 물리적인 삭제
		for _, r := range reviews {
			// 동영상 데이터가 있다면 삭제한다.
			mp4Path := fmt.Sprintf("%s/%s.mp4", CachedAdminSetting.ReviewDataPath, r.ID.Hex())
			if _, err := os.Stat(mp4Path); !os.IsNotExist(err) {
				err = os.Remove(mp4Path) // Review 데이터가 존재하면 삭제한다.
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			// 드로잉 데이터가 있다면 삭제한다.
			for _, sketch := range r.Sketches {
				if _, err := os.Stat(sketch.SketchPath); err == nil {
					err = os.Remove(sketch.SketchPath)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}
		}
		// 3. 실제 Review DB 삭제
		err = RmProjectReview(session, rcp.Project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 4. 프로젝트 삭제
	err = rmProject(session, rcp.Project)
	if err != nil {
		rcp.Error = err.Error()
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "rmproject_success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleNoOnProject 함수는 OnProject가 없을 때 접근하는 페이지이다.
func handleNoOnProject(w http.ResponseWriter, r *http.Request) {
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
	type recipe struct {
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
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u

	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "noonproject", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
