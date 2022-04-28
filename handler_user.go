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
	"os/exec"
	"path/filepath"
	"strconv"
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
	q := r.URL.Query()
	id := q.Get("id")
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		User
		QueryUser User
		SessionID string
		Devmode   bool
		SearchOption
		Setting Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.Devmode = *flagDevmode
	rcp.SessionID = ssid.ID
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
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.QueryUser, err = getUser(session, id)
	if err != nil {
		http.Redirect(w, r, "/nouser?id="+id, http.StatusSeeOther)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "user", rcp)
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
	q := r.URL.Query()
	id := q.Get("id")
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		MailDNS string
		User
		SessionUser User
		Divisions   []Division
		Departments []Department
		Teams       []Team
		Roles       []Role
		Positions   []Position
		Setting     Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.MailDNS = *flagMailDNS
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
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
	rcp.SessionUser, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "edituser", rcp)
	if err != nil {
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
	_, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 유저 정보를 가지고 온다.
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	id := r.FormValue("ID")
	if id != ssid.ID {
		if ssid.AccessLevel != AdminAccessLevel {
			http.Error(w, "사용자를 수정하기 위해서는 Accesslevel 10이 필요합니다", http.StatusUnauthorized)
			return
		}
	}
	u, err := getUser(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.EmployeeNumber = strings.TrimSpace(r.FormValue("EmployeeNumber"))
	u.RocketChatID = strings.TrimSpace(r.FormValue("RocketChatID"))
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
	u.Tags = Str2List(r.FormValue("Tags"))
	u.AccessProjects = Str2List(r.FormValue("AccessProjects"))
	u.Timezone = r.FormValue("Timezone")
	if r.FormValue("AccessLevel") != "" {
		level, err := strconv.Atoi(r.FormValue("AccessLevel"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 사용자 레벨을 업데이트한다.
		u.AccessLevel = AccessLevel(level)
		// 사용자 토큰을 업데이트한다.
		t, err := getToken(session, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.AccessLevel = AccessLevel(level)
		err = setToken(session, t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 퇴사를하게 되면 레벨:0, 토큰레벨: 0 으로 수정한다.
	if str2bool(r.FormValue("IsLeave")) {
		u.AccessLevel = AccessLevel(0)
		u.IsLeave = true
		// 사용자 토큰을 업데이트한다.
		t, err := getToken(session, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.AccessLevel = AccessLevel(0)
		err = setToken(session, t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Oraganization 정보를 분석해서 사용자에 Organization 정보를 등록한다.
	u.OrganizationsForm = r.FormValue("OrganizationsForm")
	if u.OrganizationsForm != "" {
		u.Organizations, err = OrganizationsFormToOrganizations(session, u.OrganizationsForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 조직 정보로 태그를 자동으로 생성한다.
	u.SetTags()
	file, fileHandler, fileErr := r.FormFile("Photo")
	if fileErr == nil {
		if !(fileHandler.Header.Get("Content-Type") == "image/jpeg" || fileHandler.Header.Get("Content-Type") == "image/png") {
			http.Error(w, "업로드 파일이 jpeg 또는 png 파일이 아닙니다", http.StatusInternalServerError)
			return
		}
		// 파일이 없다면 fileErr 값은 "http: no such file" 값이 된다.
		// 썸네일 파일이 존재한다면 아래 프로세스를 거친다.
		mediatype, fileParams, err := mime.ParseMediaType(fileHandler.Header.Get("Content-Disposition"))
		if err != nil {
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 사용자가 업로드한 파일을 tempFile에 복사한다.
		io.Copy(tempFile, io.LimitReader(file, MaxFileSize))
		tempFile.Close()
		defer os.Remove(tempPath)
		if CachedAdminSetting.ThumbnailRootPath == "" {
			http.Error(w, "Admin 셋팅의 ThumbnailRootPath가 설정되어있지 않습니다", http.StatusInternalServerError)
			return
		}
		thumbnailPath := fmt.Sprintf("%s/user/%s.jpg", CachedAdminSetting.ThumbnailRootPath, u.ID)
		thumbnailDir := filepath.Dir(thumbnailPath)
		// 썸네일을 생성할 user 경로가 존재하지 않는다면 생성한다.
		_, err = os.Stat(thumbnailDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(thumbnailDir, 0775)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// 이미지변환
		src, err := imaging.Open(tempPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Resize the cropped image to width = 200px preserving the aspect ratio.
		dst := imaging.Fill(src, 200, 280, imaging.Center, imaging.Lanczos)
		err = imaging.Save(dst, thumbnailPath)
		if err != nil {
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

	// 사용자 수정이후 처리할 스크립트가 admin setting에 선언되어 있다면, 실행합니다.
	setting, err := GetAdminSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if setting.RunScriptAfterEditUserProfile != "" {
		for _, line := range strings.Split(setting.RunScriptAfterEditUserProfile, "\r\n") {
			cmds := strings.Split(line, " ")
			// 해당 코드가 연산이 오래걸리면 페이지로 리다이렉션시 오래걸린다.
			if len(cmds) < 2 {
				err := exec.Command(line).Run()
				if err != nil {
					log.Println(err)
				}
			} else {
				err := exec.Command(cmds[0], cmds[1:]...).Run()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	// 리다이렉션
	if id != ssid.ID {
		// 관리자가 일반유저를 수정하는 경우 리다이렉션 URL
		http.Redirect(w, r, "/users?search="+id, http.StatusSeeOther)
	} else {
		// 자기 자신이 수정하는 경우 리다이렉션 URL
		http.Redirect(w, r, "/user?id="+ssid.ID, http.StatusSeeOther)
	}
}

// handleSignup 함수는 회원가입 페이지이다.
func handleSignup(w http.ResponseWriter, r *http.Request) {
	RmSessionID(w) // SignIn을 할 때 역시 기존의 세션을 지운다. 여러사용자 2중 로그인 방지
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
	defer session.Close()
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
	err = TEMPLATES.ExecuteTemplate(w, strings.Trim(r.URL.Path, "/"), rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleInvalidAccess 함수는 사용자의 레벨이 부족할 때 접속하는 페이지이다.
func handleInvalidAccess(w http.ResponseWriter, r *http.Request) {
	RmSessionID(w)
	err := TEMPLATES.ExecuteTemplate(w, "invalidaccess", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleInvalidPass 함수는 사용자의 패스워드가 많이 틀려서 접속되는 페이지이다.
func handleInvalidPass(w http.ResponseWriter, r *http.Request) {
	RmSessionID(w)
	err := TEMPLATES.ExecuteTemplate(w, "invalidpass", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleNoUser 함수는 등록되지 않은 사용자를 볼 때 출력되는 페이지이다.
func handleNoUser(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("id")
	type recipe struct {
		ID string `json:"id"`
	}
	rcp := recipe{}
	rcp.ID = id
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "nouser", rcp)
	if err != nil {
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
		if len(parts) != 6 {
			return results, errors.New("조직 정보가 충분하지 않습니다")
		}
		org := Organization{}
		if parts[0] == "true" {
			org.Primary = true
		} else {
			org.Primary = false
		}
		if parts[1] != "unknown" {
			division, err := getDivision(session, parts[1])
			if err != nil {
				return results, err
			}
			org.Division = division
		}
		if parts[2] != "unknown" {
			department, err := getDepartment(session, parts[2])
			if err != nil {
				return results, err
			}
			org.Department = department
		}
		if parts[3] != "unknown" {
			team, err := getTeam(session, parts[3])
			if err != nil {
				return results, err
			}
			org.Team = team
		}
		if parts[4] != "unknown" {
			role, err := getRole(session, parts[4])
			if err != nil {
				return results, err
			}
			org.Role = role
		}
		if parts[5] != "unknown" {
			position, err := getPosition(session, parts[5])
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
	w.Header().Set("Content-Type", "text/html")
	if !captcha.VerifyString(r.FormValue("CaptchaID"), r.FormValue("CaptchaNum")) {
		http.Redirect(w, r, "/error-captcha", http.StatusSeeOther)
		return
	}
	id := r.FormValue("ID")
	if id == "" {
		http.Error(w, "ID 값이 빈 문자열 입니다", http.StatusBadRequest)
		return
	}
	if !regexpUserID.MatchString(id) {
		http.Error(w, "ID값은 영문,숫자로만 이루어져야 합니다", http.StatusBadRequest)
		return
	}
	pw := r.FormValue("Password")
	if pw == "" {
		http.Error(w, "Password 값이 빈 문자열 입니다", http.StatusBadRequest)
		return
	}
	if pw != r.FormValue("ConfirmPassword") {
		http.Error(w, "입력받은 2개의 패스워드가 서로 다릅니다", http.StatusInternalServerError)
		return
	}
	u := *NewUser(id)
	pw, err := Encrypt(pw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Password = pw
	u.EmployeeNumber = strings.TrimSpace(r.FormValue("EmployeeNumber"))
	u.FirstNameKor = strings.TrimSpace(r.FormValue("FirstNameKor"))
	u.LastNameKor = strings.TrimSpace(r.FormValue("LastNameKor"))
	u.FirstNameEng = strings.Title(strings.TrimSpace(strings.ToLower(r.FormValue("FirstNameEng"))))
	u.LastNameEng = strings.Title(strings.TrimSpace(strings.ToLower(r.FormValue("LastNameEng"))))
	u.FirstNameChn = strings.TrimSpace(r.FormValue("FirstNameChn"))
	u.LastNameChn = strings.TrimSpace(r.FormValue("LastNameChn"))
	u.Email = strings.TrimSpace(r.FormValue("Email"))
	u.EmailExternal = strings.TrimSpace(r.FormValue("EmailExternal"))
	u.Phone = r.FormValue("Phone")
	u.Hotline = r.FormValue("Hotline")
	u.Location = r.FormValue("Location")
	u.Tags = Str2List(r.FormValue("Tags"))
	u.Timezone = r.FormValue("Timezone")
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
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
	if u.OrganizationsForm != "" {
		u.Organizations, err = OrganizationsFormToOrganizations(session, u.OrganizationsForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	// JWT 토큰으로 쿠키를 저장한다.
	SetSessionID(w, u.ID, u.AccessLevel, "")
	// 가입이후 처리할 스크립트가 admin setting에 선언되어 있다면, 실행합니다.
	setting, err := GetAdminSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if setting.RunScriptAfterSignup != "" {
		for _, line := range strings.Split(setting.RunScriptAfterSignup, "\r\n") {
			cmds := strings.Split(line, " ")
			if len(cmds) < 2 {
				err := exec.Command(line).Run()
				if err != nil {
					log.Println(err)
				}
			} else {
				err := exec.Command(cmds[0], cmds[1:]...).Run()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	// 가입완료 페이지로 이동
	err = TEMPLATES.ExecuteTemplate(w, "signup_success", u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSignin 함수는 로그인 페이지이다.
func handleSignin(w http.ResponseWriter, r *http.Request) {
	RmSessionID(w) // SignIn을 할 때 역시 기존의 세션을 지운다. 여러사용자 2중 로그인 방지
	type recipe struct {
		Company string
		Message string
		ID      string
	}
	rcp := recipe{}
	q := r.URL.Query()
	errorCode := q.Get("status")
	rcp.ID = q.Get("id")
	passwordAttempt := q.Get("passwordattempt")
	switch errorCode {
	case "wrongpw":
		if passwordAttempt != "" {
			rcp.Message = passwordAttempt + "회 패스워드를 틀렸습니다. 다시 로그인 해주세요."
		} else {
			rcp.Message = "패스워드를 틀렸습니다. 다시 로그인 해주세요."
		}
	}
	rcp.Company = strings.Title(*flagCompany)
	if rcp.ID == "" {
		// ID가 없다면 저장된 쿠키 ID를 가지고 온다.
		c, err := r.Cookie("CookieUserID")
		if err != nil {
			rcp.ID = ""
		} else {
			rcp.ID = c.Value
		}
	}
	err := TEMPLATES.ExecuteTemplate(w, "signin", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSigninSubmit 함수는 사용자 로그인값을 처리하는 핸들러이다.
func handleSigninSubmit(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("ID")
	if !regexpUserID.MatchString(id) {
		http.Error(w, "ID값은 영문,숫자로만 이루어져야 합니다", http.StatusBadRequest)
		return
	}
	if id == "" {
		http.Error(w, "ID 값이 빈 문자열 입니다", http.StatusBadRequest)
		return
	}
	pw := r.FormValue("Password")
	if pw == "" {
		http.Error(w, "Password 값이 빈 문자열 입니다", http.StatusBadRequest)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	// 사용자가 과거에 패스워드를 5회이상 틀렸다면 로그인을 허용하지 않는다.
	u, err := getUser(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if u.PasswordAttempt > 4 {
		http.Redirect(w, r, "/invalidpass", http.StatusSeeOther)
		return
	}
	err = vaildUser(session, id, pw)
	if err != nil {
		// 패스워드 시도횟수를 추가한다.
		addPasswordAttempt(session, id)
		// 패스워드 시도횟수를 가지고 오기 위해서 사용자 정보를 가지고 온다.
		u, err := getUser(session, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 다시 로그인 페이지로 리다이렉트한다.
		http.Redirect(w, r, fmt.Sprintf("/signin?status=wrongpw&passwordattempt=%d&id=%s", u.PasswordAttempt, id), http.StatusSeeOther)
		return
	}

	// 로그인에 성공하면 접속한 아이피와 포트를 DB에 기록한다.
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.LastIP = host
	u.LastPort = port
	u.PasswordAttempt = 0 // 로그인에 성공하면 기존 시도한 패스워드 횟수를 초기화 한다.
	err = setUser(session, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	// 로그인에 성공하면 CookieUserID 값을 쿠키에 저장한다.
	c := http.Cookie{
		Name:   "CookieUserID",
		Value:  ssid.ID,
		MaxAge: int(*flagCookieAge),
	}
	http.SetCookie(w, &c)

	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	type recipe struct {
		User
	}
	rcp := recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "signin_success", rcp)
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
	RmSessionID(w)
	err = TEMPLATES.ExecuteTemplate(w, "signout", nil)
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
	q := r.URL.Query()
	id := q.Get("id")
	w.Header().Set("Content-Type", "text/html")
	type recipe struct {
		User
		Devmode bool
		SearchOption
		Setting Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.Devmode = *flagDevmode
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
	rcp.User, err = getUser(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "updatepassword", rcp)
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
		http.Error(w, "Password 값이 빈 문자열 입니다", http.StatusInternalServerError)
		return
	}
	if r.FormValue("NewPassword") != r.FormValue("ConfirmNewPassword") {
		http.Error(w, "Password 값이 빈 문자열 입니다", http.StatusInternalServerError)
		return
	}
	pw := r.FormValue("OldPassword")
	newPw := r.FormValue("NewPassword")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
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

// handleUsers 함수는 유저리스트를 검색하는 페이지이다. (기본정렬: 사번순)
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
	searchword := q.Get("search")
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
		SearchOption
		Setting Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	err = TEMPLATES.ExecuteTemplate(w, "users", rcp)
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
	type recipe struct {
		User // 로그인한 사용자 정보
	}
	rcp := recipe{}
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "replacetag", rcp)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	// Tags replace
	err = ReplaceTags(session, src, dst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// users 리다이렉트한다.
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
