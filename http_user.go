package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/disintegration/imaging"
	"gopkg.in/mgo.v2"
)

// handleUser 함수는 유저정보를 출력하는 페이지이다.
func handleUser(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	id := q.Get("id")
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		User
		SessionID string
		Devmode   bool
	}
	rcp := recipe{}
	rcp.Devmode = *flagDevmode
	rcp.SessionID = ssid.ID
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User, err = getUser(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "user", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// handleEditUser 함수는 유저정보를 수정하는 페이지이다.
func handleEditUser(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	id := q.Get("id")
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		MailDNS string
		User
		Divisions   []Division
		Departments []Department
		Teams       []Team
		Roles       []Role
		Positions   []Position
	}
	rcp := recipe{}
	rcp.MailDNS = *flagMailDNS
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Divisions, err = allDivisions(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Departments, err = allDepartments(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Teams, err = allTeams(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Roles, err = allRoles(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Positions, err = allPositions(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User, err = getUser(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "edituser", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleEditUserSubmit 함수는 회원정보를 수정받는 페이지이다.
func handleEditUserSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 유저 정보를 가지고 온다.
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUser(session, ssid.ID)
	u.FirstNameKor = r.FormValue("FirstNameKor")
	u.LastNameKor = r.FormValue("LastNameKor")
	u.FirstNameEng = strings.Title(strings.ToLower(r.FormValue("FirstNameEng")))
	u.LastNameEng = strings.Title(strings.ToLower(r.FormValue("LastNameEng")))
	u.FirstNameChn = r.FormValue("FirstNameChn")
	u.LastNameChn = r.FormValue("LastNameChn")
	u.Email = r.FormValue("Email")
	u.EmailExternal = r.FormValue("EmailExternal")
	u.Phone = r.FormValue("Phone")
	u.Hotline = r.FormValue("Hotline")
	u.Location = r.FormValue("Location")
	u.Tags = Str2Tags(r.FormValue("Tags"))

	u.Timezone = r.FormValue("Timezone")
	u.LastIP = host
	u.LastPort = port

	// Oraganization 정보를 분석해서 사용자에 Organization 정보를 등록한다.
	u.OrganizationsForm = r.FormValue("OrganizationsForm")
	u.Organizations, err = OrganizationsFormToOrganizations(session, u.OrganizationsForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 조직 정보로 태그를 자동으로 생성한다.
	u.SetTags()
	file, fileHandler, fileErr := r.FormFile("Photo")
	if fileErr == nil {
		if !(fileHandler.Header.Get("Content-Type") == "image/jpeg" || fileHandler.Header.Get("Content-Type") == "image/png") {
			err := errors.New("업로드 파일이 jpeg 또는 png 파일이 아닙니다")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		thumbnailPath := fmt.Sprintf("%s/user/%s.jpg", *flagThumbPath, u.ID)
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
		}
		// 이미지변환
		src, err := imaging.Open(tempPath)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Resize the cropped image to width = 200px preserving the aspect ratio.
		dst := imaging.Fill(src, 200, 280, imaging.Center, imaging.Lanczos)
		err = imaging.Save(dst, thumbnailPath)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.Thumbnail = true
	}
	err = setUser(session, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/user?id="+ssid.ID, http.StatusSeeOther)
}

// handleSignup 함수는 회원가입 페이지이다.
func handleSignup(w http.ResponseWriter, r *http.Request) {
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		Company     string
		CaptchaID   string
		MailDNS     string
		Divisions   []Division
		Departments []Department
		Teams       []Team
		Roles       []Role
		Positions   []Position
	}
	rcp := recipe{}
	rcp.Company = strings.Title(*flagCompany)
	rcp.MailDNS = *flagMailDNS
	rcp.CaptchaID = captcha.New()
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Divisions, err = allDivisions(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Departments, err = allDepartments(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Teams, err = allTeams(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Roles, err = allRoles(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Positions, err = allPositions(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, strings.Trim(r.URL.Path, "/"), rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleInvalidAccess 함수는 사용자의 레벨이 부족할 때 접속하는 페이지이다.
func handleInvalidAccess(w http.ResponseWriter, r *http.Request) {
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.ExecuteTemplate(w, "invalidaccess", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// OrganizationsFormToOrganizations 함수는 form 문자를 받아서 []Organization 을 생성한다.
func OrganizationsFormToOrganizations(session *mgo.Session, s string) ([]Organization, error) {
	var results []Organization
	orgs := strings.Split(s, ":")
	for _, org := range orgs {
		parts := strings.Split(org, ",")
		if len(parts) != 5 {
			return results, errors.New("조직 정보가 5개가 아닙니다")
		}
		org := Organization{}
		if parts[0] != "unknown" {
			division, err := getDivision(session, parts[0])
			if err != nil {
				return results, err
			}
			org.Division = division
		}
		if parts[1] != "unknown" {
			department, err := getDepartment(session, parts[1])
			if err != nil {
				return results, err
			}
			org.Department = department
		}
		if parts[2] != "unknown" {
			team, err := getTeam(session, parts[2])
			if err != nil {
				return results, err
			}
			org.Team = team
		}
		if parts[3] != "unknown" {
			role, err := getRole(session, parts[3])
			if err != nil {
				return results, err
			}
			org.Role = role
		}
		if parts[4] != "unknown" {
			position, err := getPosition(session, parts[4])
			if err != nil {
				return results, err
			}
			org.Position = position
		}
		results = append(results, org)
	}
	return results, nil
}

// handleSignupSubmit 함수는 회원가입 페이지이다.
func handleSignupSubmit(w http.ResponseWriter, r *http.Request) {
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if !captcha.VerifyString(r.FormValue("CaptchaID"), r.FormValue("CaptchaNum")) {
		err := errors.New("CaptchaID 값과 CaptchaNum 값이 다릅니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("ID") == "" {
		err := errors.New("ID 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("Password") == "" {
		err := errors.New("Password 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("Password") != r.FormValue("ConfirmPassword") {
		err := errors.New("입력받은 2개의 패스워드가 서로 다릅니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := strings.ToLower(r.FormValue("ID"))
	u := *NewUser(id)
	pw, err := Encrypt(r.FormValue("Password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Password = pw
	u.FirstNameKor = r.FormValue("FirstNameKor")
	u.LastNameKor = r.FormValue("LastNameKor")
	u.FirstNameEng = strings.Title(strings.ToLower(r.FormValue("FirstNameEng")))
	u.LastNameEng = strings.Title(strings.ToLower(r.FormValue("LastNameEng")))
	u.FirstNameChn = r.FormValue("FirstNameChn")
	u.LastNameChn = r.FormValue("LastNameChn")
	u.Email = r.FormValue("Email")
	u.EmailExternal = r.FormValue("EmailExternal")
	u.Phone = r.FormValue("Phone")
	u.Hotline = r.FormValue("Hotline")
	u.Location = r.FormValue("Location")
	u.Tags = Str2Tags(r.FormValue("Tags"))
	u.Timezone = r.FormValue("Timezone")
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.LastIP = host
	u.LastPort = port
	u.Token = base64.StdEncoding.EncodeToString([]byte(pw))
	// 이 데이터가 DB로 들어가야 한다.
	session, err := mgo.DialWithTimeout(*flagDBIP, 2*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Oraganization 정보를 분석해서 사용자에 Organization 정보를 등록한다.
	u.OrganizationsForm = r.FormValue("OrganizationsForm")
	u.Organizations, err = OrganizationsFormToOrganizations(session, u.OrganizationsForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 조직 정보로 태그를 자동으로 생성한다.
	u.SetTags()
	err = addUser(session, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = addToken(session, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SetSessionID(w, u.ID, u.AccessLevel, "")
	// 가입완료 페이지로 이동
	err = t.ExecuteTemplate(w, "signup_success", u)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSignin 함수는 로그인 페이지이다.
func handleSignin(w http.ResponseWriter, r *http.Request) {
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		Company string
	}
	rcp := recipe{}
	rcp.Company = strings.Title(*flagCompany)
	err = t.ExecuteTemplate(w, "signin", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSigninSubmit 함수는 회원가입 페이지이다.
func handleSigninSubmit(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("ID") == "" {
		err := errors.New("ID 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("Password") == "" {
		err := errors.New("Password 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.FormValue("ID")
	pw := r.FormValue("Password")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	err = vaildUser(session, id, pw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 로그인에 성공하면 접속한 아이피와 포트를 DB에 기록한다.
	u, err := getUser(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.LastIP = host
	u.LastPort = port
	err = setUser(session, u)
	// session을 저장후 로그인 성공페이지로 이동한다.
	err = SetSessionID(w, u.ID, u.AccessLevel, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/signin_success", http.StatusSeeOther)
}

func handleSigninSuccess(w http.ResponseWriter, r *http.Request) {
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
	err = t.ExecuteTemplate(w, "signin_success", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSignout 함수는 로그아웃 페이지이다.
func handleSignout(w http.ResponseWriter, r *http.Request) {
	_, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	RmSessionID(w)
	err = t.ExecuteTemplate(w, "signout", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUpdatePassword 함수는 사용자의 패스워드를 수정하는 페이지이다.
func handleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	_, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q := r.URL.Query()
	id := q.Get("id")
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		User
		Devmode bool
	}
	rcp := recipe{}
	rcp.Devmode = *flagDevmode
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User, err = getUser(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "updatepassword", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUpdatePasswordSubmit 함수는 회원가입 페이지이다.
func handleUpdatePasswordSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.ID == "" {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if r.FormValue("OldPassword") == "" {
		err := errors.New("Password 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("NewPassword") == "" {
		err := errors.New("Password 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("NewPassword") != r.FormValue("ConfirmNewPassword") {
		err := errors.New("Password 값이 빈 문자열 입니다")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pw := r.FormValue("OldPassword")
	newPw := r.FormValue("NewPassword")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	err = updatePasswordUser(session, ssid.ID, pw, newPw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 기존 토큰을 제거한다.
	err = rmToken(session, ssid.ID)
	if err != nil {
		log.Println(err)
	}
	// 새로운 사용자 정보를 불러와서 토큰을 생성한다.
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = addToken(session, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 기존 쿠키를 제거하고 새로 다시 로그인을 합니다.
	RmSessionID(w)
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

// handleUsers 함수는 유저리스트를 검색하는 페이지이다. (기본 정렬은 사번순이다.)
func handleUsers(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		searchword := r.FormValue("searchword")
		http.Redirect(w, r, "/users?search="+searchword, http.StatusSeeOther)
		return
	}
	q := r.URL.Query()
	var searchword string
	searchword = q.Get("search")
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
		User                // 로그인한 사용자 정보
		Users      []User   // 검색된 사용자를 담는 리스트
		Tags       []string // 부서,파트 태그
		Searchword string   // searchform에 들어가는 문자
		Usernum    int      // 검색된 인원수
		Devmode    bool     // 개발모드
	}
	rcp := recipe{}
	rcp.Devmode = *flagDevmode
	rcp.Searchword = searchword
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	searchwords := strings.Split(strings.ReplaceAll(searchword, " ", ","), ",")
	rcp.Users, err = searchUsers(session, searchwords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Usernum = len(rcp.Users)
	rcp.Tags, err = UserTags(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "users", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// handleReplaceTag 함수는 유저에 설정된 태그를 변경하는 페이지이다.
func handleReplaceTag(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type recipe struct {
		User // 로그인한 사용자 정보
	}
	rcp := recipe{}
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "replacetag", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleReplaceTagSubmit 함수는 유저에 설정된 부서 태그를 변경하는 페이지이다.
func handleReplaceTagSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	src := r.FormValue("src")
	dst := r.FormValue("dst")

	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Tags replace
	err = ReplaceTags(session, src, dst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// users 리다리렉트한다.
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
