package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
)

// handleImportExcel 함수는 엑셀파일을 Import 하는 페이지 이다.
func handleImportExcel(w http.ResponseWriter, r *http.Request) {
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
	defer session.Close()
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "importexcel", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleImportExcel 함수는 엑셀파일을 Import 하는 페이지 이다.
func handlePresetExcel(w http.ResponseWriter, r *http.Request) {
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
		Files       []string
		SearchOption
	}
	rcp := recipe{}
	rcp.Devmode = *flagDevmode
	rcp.SessionID = ssid.ID
	tmp, err := userTemppath(ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	f, err := os.Open(tmp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, file := range files {
		rcp.Files = append(rcp.Files, file.Name())
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
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
	if len(rcp.Projectlist) == 0 {
		http.Redirect(w, r, "/noonproject", http.StatusSeeOther)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "presetexcel", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUploadExcel(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	// dropzone setting
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	mimeType := header.Header.Get("Content-Type")
	switch mimeType {
	case "text/csv":
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		tmp, err := userTemppath(ssid.ID)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
		}
		path := tmp + "/" + header.Filename // 업로드한 파일 리스트를 불러오기 위해 뒤에 붙는 Unixtime을 제거한다.
		err = ioutil.WriteFile(path, data, 0666)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		fmt.Println(path)
	case "application/vnd.ms-excel":
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		tmp, err := userTemppath(ssid.ID)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
		}
		path := tmp + "/" + header.Filename // 업로드한 파일 리스트를 불러오기 위해 뒤에 붙는 Unixtime을 제거한다.
		err = ioutil.WriteFile(path, data, 0666)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
	default:
		// 지원하지 않는 파일. 저장하지 않는다.

	}
}

// RowCSV 자료구조는 .csv 형식의 자료구조이다.
// 샷네임;작업종류(2D,3D);작업내용;수정사항;링크자료(제목:경로);3D마감;2D마감;FIN날짜;FIN버젼;테그;롤넘버;핸들IN;핸들OUT;JUST타임코드IN;JUST타임코드OUT
type RowCSV struct {
	Name            string
	Shottype        string
	Note            string
	Comment         string
	Link            string
	Ddline3D        string
	Ddline2D        string
	Findate         string
	Finver          string
	Tags            string
	Rnum            string
	HandleIn        string
	HandleOut       string
	JustTimecodeIn  string
	JustTimecodeOut string
}

// handlePresetExcelSubmit 함수는 excel 파일을 체크하고 분석 보고서로 Redirection 한다.
func handlePresetExcelSubmit(w http.ResponseWriter, r *http.Request) {
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
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	project := r.FormValue("project")
	filename := r.FormValue("filename")
	separator := r.FormValue("separator")
	overwrite := str2bool(r.FormValue("overwrite"))
	tmppath, err := userTemppath(ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// csv 파일을 읽는다.
	csvFile, err := os.Open(tmppath + "/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reader := csv.NewReader(csvFile)
	switch separator {
	case ";":
		reader.Comma = ';'
	case ":":
		reader.Comma = ':'
	default:
		reader.Comma = ','
	}
	reader.Comment = '#'
	var rows []RowCSV
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			http.Error(w, ".csv 파일을 읽다가 에러가 발생했습니다", http.StatusBadRequest)
			return
		}
		if len(line) != 15 { // Template 요소의 갯수가 다르다.
			http.Error(w, ".csv 열의 갯수가 약속된 형식과 다릅니다", http.StatusBadRequest)
			return
		}
		row := RowCSV{}
		row.Name = line[0]             // item name
		row.Shottype = line[1]         // shottype 2d,3d
		row.Note = line[2]             // 작업내용
		row.Comment = line[3]          // 수정사항
		row.Link = line[4]             // 링크자료(제목:경로)
		row.Ddline3D = line[5]         // 3D마감
		row.Ddline2D = line[6]         // 2D마감
		row.Findate = line[7]          // FIN날짜
		row.Finver = line[8]           // FIN버전
		row.Tags = line[9]             // 태그
		row.Rnum = line[10]            // 롤넘버
		row.HandleIn = line[11]        // 핸들IN
		row.HandleOut = line[12]       // 핸들OUT
		row.JustTimecodeIn = line[13]  // JUST타임코드IN
		row.JustTimecodeOut = line[14] // JUST타임코드OUT
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		http.Error(w, ".csv 파일 내부가 비어있습니다", http.StatusBadRequest)
		return
	}
	type recipe struct {
		Project   string
		Filename  string
		Separator string
		Overwrite bool
		Rows      []RowCSV
		User
		SessionID string
		Devmode   bool
		SearchOption
	}
	rcp := recipe{}
	rcp.SessionID = ssid.ID
	rcp.Devmode = *flagDevmode

	rcp.SearchOption = handleRequestToSearchOption(r)
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Rows = rows
	rcp.Project = project
	rcp.Filename = filename
	rcp.Separator = separator
	rcp.Overwrite = overwrite
	fmt.Println(rows)
	// project에 샷이 존재하는지 체크한다.
	// 각 col 값이 정상적인지 체크한다.
	// 결과 보고서를 만든다.
	err = TEMPLATES.ExecuteTemplate(w, "reportcsv", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
