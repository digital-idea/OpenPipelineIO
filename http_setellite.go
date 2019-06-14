package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/mgo.v2"
)

// handleSetellite 함수는 현장데이터를 출력하는 페이지이다.
func handleSetellite(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		templates.ExecuteTemplate(w, "dberr", nil)
		return
	}
	defer session.Close()
	if r.Method == http.MethodPost {
		project := r.FormValue("project")
		searchWord := r.FormValue("searchWord")
		url := fmt.Sprintf("/setellite?project=%s&searchword=%s", project, searchWord)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
	w.Header().Set("Content-Type", "text/html; charset=\"utf-8\"")
	type recipe struct {
		Projectlist []string
		Project     string
		Searchword  string
		Items       []Setellite
		Error       string
	}
	rcp := recipe{}
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		rcp.Error = err.Error()
		templates.ExecuteTemplate(w, "setellite", rcp)
		return
	}
	r.ParseForm()
	for key, value := range r.Form {
		switch key {
		case "project":
			if len(value) != 1 {
				rcp.Error = "프로젝트값을 1개만 설정해주세요."
				templates.ExecuteTemplate(w, "setellite", rcp)
				return
			}
			rcp.Project = value[0]
		case "searchword":
			if len(value) != 1 {
				rcp.Error = "Searchword 값을 1개만 설정해주세요."
				templates.ExecuteTemplate(w, "setellite", rcp)
				return
			}
			rcp.Searchword = value[0]
		}
	}
	if rcp.Searchword == "" {
		rcp.Error = "검색어를 입력해주세요."
		templates.ExecuteTemplate(w, "setellite", rcp)
		return
	}
	rcp.Items, err = SetelliteSearch(session, rcp.Project, rcp.Searchword)
	if err != nil {
		rcp.Error = err.Error()
	}
	if len(rcp.Items) == 0 {
		rcp.Error = "검색결과 0 건"
	}
	templates.ExecuteTemplate(w, "setellite", rcp)
}

func handleUploadSetellite(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		templates.ExecuteTemplate(w, "dberr", nil)
		return
	}
	defer session.Close()
	type recipe struct {
		Projectlist []string
		Message     string
		Errors      []string // CSV를 처리하면서 각 라인별로 에러가 있다면 에러내용을 저장한다.
	}
	rcp := recipe{}
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		rcp.Message = err.Error()
		err = templates.ExecuteTemplate(w, "uploadSetellite", rcp)
		return
	}
	if r.Method == http.MethodGet {
		rcp.Message = "Setellite CSV파일을 업로드해주세요."
		w.Header().Set("Content-Type", "text/html")
		err = templates.ExecuteTemplate(w, "uploadSetellite", rcp)
		if err != nil {
			rcp.Message = err.Error()
			templates.ExecuteTemplate(w, "uploadSetellite", rcp)
			return
		}
	} else if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "text/html")
		project := r.FormValue("project")
		file, fileHandle, err := r.FormFile("csv")
		if err != nil {
			rcp.Message = err.Error()
			templates.ExecuteTemplate(w, "uploadSetellite", rcp)
			return
		}
		mediatype, fileParams, err := mime.ParseMediaType(fileHandle.Header.Get("Content-Disposition"))
		if err != nil {
			rcp.Message = err.Error()
			templates.ExecuteTemplate(w, "uploadSetellite", rcp)
			return
		}
		if mediatype != "form-data" {
			rcp.Message = "지원하지 않는 미디어 입니다."
			templates.ExecuteTemplate(w, "uploadSetellite", rcp)
			return
		}
		// csv 확장자 체크
		if ".csv" != strings.ToLower(filepath.Ext(fileParams["filename"])) {
			rcp.Message = "CSV 파일이 아닙니다."
			templates.ExecuteTemplate(w, "uploadSetellite", rcp)
			return
		}
		f, err := ioutil.TempFile(os.TempDir(), "")
		if err != nil {
			rcp.Message = err.Error()
			templates.ExecuteTemplate(w, "uploadSetellite", rcp)
			return
		}
		file.Close()
		defer os.Remove(f.Name())
		io.Copy(f, io.LimitReader(file, MaxFileSize))
		items := csv2SetelliteList(f.Name())
		for num, item := range items {
			err = addSetellite(session, project, item, true)
			if err != nil {
				rcp.Errors = append(rcp.Errors, fmt.Sprintf("%d줄 : %s", num+1, err.Error()))
			}
		}
		rcp.Message = "파일이 업로드 되었습니다. 업로드할 다른 CSV가 있다면 업로드해주세요."
		templates.ExecuteTemplate(w, "uploadSetellite", rcp)
		return
	} else {
		http.Error(w, "Get,Post Only", http.StatusMethodNotAllowed)
	}
}
