package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
)

// handleInputMode 함수는 수정을 편하게 입력하는 페이지 이다.
func handleInputMode(w http.ResponseWriter, r *http.Request) {
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
	type recipe struct {
		User
		SessionID   string
		Devmode     bool
		Projectlist []string
		Items       []Item
		Ddline3d    []string
		Ddline2d    []string
		Tags        []string
		Assettags   []string
		SearchOption
		Searchnum           Infobarnum
		Totalnum            Infobarnum
		Projectinfo         Project
		MailDNS             string
		Dilog               string
		Wfs                 string
		OS                  string
		TasksettingNames    []string
		TasksettingOrderMap map[string]float64
		Dday                string
		Status              []Status
		AllStatusIDs        []string
	}
	rcp := recipe{}
	_, rcp.OS, _ = GetInfoFromRequestHeader(r)
	rcp.Devmode = *flagDevmode
	rcp.MailDNS = *flagMailDNS
	rcp.Dilog = *flagDILOG
	rcp.Wfs = *flagWFS
	rcp.SessionID = ssid.ID
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	tasks, err := AllTaskSettings(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.TasksettingOrderMap = make(map[string]float64)
	for _, t := range tasks {
		rcp.TasksettingOrderMap[t.Name] = t.Order
	}
	rcp.TasksettingNames, err = TasksettingNames(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Status, err = AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, status := range rcp.Status {
		rcp.AllStatusIDs = append(rcp.AllStatusIDs, status.ID)
	}
	rcp.SearchOption = handleRequestToSearchOption(r)
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Projectlist, err = OnProjectlist(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 만약 사용자에게 AccessProjects가 설정되어있다면 해당리스트를 사용한다.
	if len(rcp.User.AccessProjects) != 0 {
		var accessProjects []string
		for _, i := range rcp.Projectlist {
			for _, j := range rcp.User.AccessProjects {
				if i != j {
					continue
				}
				accessProjects = append(accessProjects, j)
			}
		}
		rcp.Projectlist = accessProjects
	}
	// 마이크레이션중인 프로젝트는 제거한다.
	setting, err := GetAdminSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var excludeProjects []string
	if setting.ExcludeProject != "" {
		excludeProjects = strings.Split(strings.Replace(setting.ExcludeProject, " ", "", -1), ",")
		var accessProjects []string
		for _, p := range rcp.Projectlist {
			has := false
			for _, ep := range excludeProjects {
				if p == ep {
					has = true
				}
			}
			if !has { // 제외할 프로젝트가 아니면 프로젝트 리스트에 포함한다.
				accessProjects = append(accessProjects, p)
			}
		}
		rcp.Projectlist = accessProjects
	}
	if len(rcp.Projectlist) == 0 {
		http.Redirect(w, r, "/noonproject", http.StatusSeeOther)
		return
	}
	rcp.Ddline3d, err = DistinctDdline(session, rcp.SearchOption.Project, "ddline3d")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Ddline2d, err = DistinctDdline(session, rcp.SearchOption.Project, "ddline2d")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Tags, err = Distinct(session, rcp.SearchOption.Project, "tag")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Assettags, err = Distinct(session, rcp.SearchOption.Project, "assettags")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Totalnum, err = Totalnum(session, rcp.SearchOption.Project)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Totalnum.calculatePercent()
	if rcp.SearchOption.Project != "" {
		rcp.Projectinfo, err = getProject(session, rcp.SearchOption.Project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dday, err := ToDday(rcp.Projectinfo.Deadline)
		if err != nil {
			log.Println(err)
		}
		rcp.Dday = dday
	}
	rcp.Items, err = Searchv2(session, rcp.SearchOption)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Searchnum, err = Searchnum(rcp.SearchOption, rcp.Items)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 최종적으로 사용된 프로젝트명을 쿠키에 저장한다.
	cookie := http.Cookie{
		Name:   "Project",
		Value:  rcp.SearchOption.Project,
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Task",
		Value:  rcp.SearchOption.Task,
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Searchword",
		Value:  base64.StdEncoding.EncodeToString([]byte(rcp.SearchOption.Searchword)), //  쿠키는 UTF-8을 저장할 때 에러가 발생한다.
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Assign",
		Value:  bool2str(rcp.SearchOption.Assign),
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Ready",
		Value:  bool2str(rcp.SearchOption.Ready),
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Wip",
		Value:  bool2str(rcp.SearchOption.Wip),
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Confirm",
		Value:  bool2str(rcp.SearchOption.Confirm),
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Done",
		Value:  bool2str(rcp.SearchOption.Done),
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Omit",
		Value:  bool2str(rcp.SearchOption.Omit),
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Hold",
		Value:  bool2str(rcp.SearchOption.Hold),
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Out",
		Value:  bool2str(rcp.SearchOption.Out),
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "None",
		Value:  bool2str(rcp.SearchOption.None),
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Template",
		Value:  rcp.SearchOption.Template,
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)

	err = TEMPLATES.ExecuteTemplate(w, rcp.SearchOption.Template, rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
