package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/dchest/captcha"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"gopkg.in/mgo.v2"
)

// MaxFileSize 사이즈는 웹에서 전송할 수 있는 최대 사이즈를 2기가로 제한한다.(인트라넷)
const MaxFileSize = 2000 * 1000 * 1024

// LoadTemplates 함수는 템플릿을 로딩합니다.
func LoadTemplates() (*template.Template, error) {
	t := template.New("").Funcs(funcMap)
	t, err := vfstemplate.ParseGlob(assets, t, "/template/*.html")
	return t, err
}

//템플릿 함수를 로딩합니다.
var funcMap = template.FuncMap{
	"title":               strings.Title,
	"itemStatus2color":    itemStatus2color,
	"projectStatus2color": projectStatus2color,
	"statusnum2string":    statusnum2string,
	"name2seq":            name2seq,
	"note2body":           note2body,
	"pmnote2body":         pmnote2body,
	"GetPath":             GetPath,
	"ReverseStringSlice":  ReverseStringSlice,
	"ToShortTime":         ToShortTime,
	"Tags2str":            Tags2str,
	"CheckDate":           CheckDate,
	"CheckUpdate":         CheckUpdate,
	"CheckDdline":         CheckDdline,
	"ToHumantime":         ToHumantime,
	"Framecal":            Framecal,
	"Add":                 Add,
	"Minus":               Minus,
	"Review":              Review,
	"Scanname2RollMedia":  Scanname2RollMedia,
	"Hashtag2tag":         Hashtag2tag,
	"Username2Elements":   Username2Elements,
	"RemovePath":          RemovePath,
	"ShortPhoneNum":       ShortPhoneNum,
}

// 도움말 페이지 입니다.
func handleHelp(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type recipy struct {
		Wfs  string
		User User
	}
	rcp := recipy{}
	rcp.User = u
	rcp.Wfs = *flagWFS
	err = t.ExecuteTemplate(w, "help", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 전송되는 컨텐츠의 캐쉬 수명을 설정하는 핸들러입니다.
func maxAgeHandler(seconds int, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))
		h.ServeHTTP(w, r)
	})
}

// webserver함수는 웹서버의 URL을 선언하는 함수입니다.
func webserver(port string) {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assets)))
	http.Handle("/thumbnail/", maxAgeHandler(3600, http.StripPrefix("/thumbnail/", http.FileServer(http.Dir(*flagThumbPath)))))
	http.HandleFunc("/", handleIndex)
	// Item
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/tag/", handleTags)
	http.HandleFunc("/assettags/", handleAssettags)
	http.HandleFunc("/ddline/", handleDdline)
	http.HandleFunc("/edit", handleEdit)
	http.HandleFunc("/edit_item_submit", handleEditItemSubmit)
	http.HandleFunc("/help", handleHelp)
	http.HandleFunc("/setellite", handleSetellite)
	http.HandleFunc("/uploadsetellite", handleUploadSetellite)

	// Project
	http.HandleFunc("/projectinfo", handleProjectinfo)
	http.HandleFunc("/addproject", handleAddProject)
	http.HandleFunc("/addproject_submit", handleAddProjectSubmit)
	http.HandleFunc("/editproject", handleEditProject)
	http.HandleFunc("/editproject_submit", handleEditProjectSubmit)

	// User
	http.HandleFunc("/signup", handleSignup)
	http.HandleFunc("/signup_submit", handleSignupSubmit)
	http.HandleFunc("/signin", handleSignin)
	http.HandleFunc("/signin_submit", handleSigninSubmit)
	http.HandleFunc("/signin_success", handleSigninSuccess)
	http.HandleFunc("/signout", handleSignout)
	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/updatepassword", handleUpdatePassword)
	http.HandleFunc("/updatepassword_submit", handleUpdatePasswordSubmit)
	http.HandleFunc("/edituser", handleEditUser)
	http.HandleFunc("/edituser_submit", handleEditUserSubmit)
	http.HandleFunc("/replacepart", handleReplacePart)
	http.HandleFunc("/replacepart_submit", handleReplacePartSubmit)

	// rest API / API2
	http.HandleFunc("/api/project", handleAPIProject)
	http.HandleFunc("/api/projects", handleAPIProjects)
	http.HandleFunc("/api/addproject", handleAPIAddproject)
	http.HandleFunc("/api/projecttags", handleAPIProjectTags)
	http.HandleFunc("/api/item", handleAPIItem)
	http.HandleFunc("/api2/items", handleAPI2Items)
	http.HandleFunc("/api/searchname", handleAPISearchname)
	http.HandleFunc("/api/seqs", handleAPISeqs)
	http.HandleFunc("/api/shots", handleAPIShots)
	http.HandleFunc("/api/shot", handleAPIShot)
	http.HandleFunc("/api/setellite", handleAPISetelliteItems)
	http.HandleFunc("/api/setellitesearch", handleAPISetelliteSearch)
	http.HandleFunc("/api/setmov", handleAPISetmov)
	http.HandleFunc("/api/setplatesize", handleAPISetPlateSize)
	http.HandleFunc("/api/setdistortionsize", handleAPISetDistortionSize)
	http.HandleFunc("/api/setrendersize", handleAPISetRenderSize)
	http.HandleFunc("/api/setcamerapubpath", handleAPISetCameraPubPath)
	http.HandleFunc("/api/setcamerapubtask", handleAPISetCameraPubTask)
	http.HandleFunc("/api/setcameraprojection", handleAPISetCameraProjection)
	http.HandleFunc("/api/setthummov", handleAPISetThummov)
	http.HandleFunc("/api/setstatus", handleAPISetStatus)
	http.HandleFunc("/api/setjustin", handleAPISetJustIn)
	http.HandleFunc("/api/setjustout", handleAPISetJustOut)
	http.HandleFunc("/api/setstartdate", handleAPISetStartdate)
	http.HandleFunc("/api/setpredate", handleAPISetPredate)
	http.HandleFunc("/api/user", handleAPIUser)
	http.HandleFunc("/api/users", handleAPISearchUser)

	// Deprecated: 사용하지 않는 API, 과거호환성을 위해서 남겨둠
	http.HandleFunc("/api/items", handleAPIItems)

	// Web Cmd
	http.HandleFunc("/cmd", handleCmd) // 리펙토링이 필요해보임.

	// Captcha
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	if port == ":443" || port == ":8443" {
		err := http.ListenAndServeTLS(port, *flagCertFullchanin, *flagCertPrivkey, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
