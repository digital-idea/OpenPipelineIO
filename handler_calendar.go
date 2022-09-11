package main

import (
	"net/http"

	"github.com/dchest/captcha"
	"gopkg.in/mgo.v2"
)

func handleCalendarTest(w http.ResponseWriter, r *http.Request) {
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
	rcp.Company = *flagCompany
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
	err = TEMPLATES.ExecuteTemplate(w, "calendartest", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
