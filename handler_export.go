package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/digital-idea/dilog"
	"github.com/digital-idea/ditime"
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
		SessionID   string
		Devmode     bool
		Projectlist []string
		Setting     Setting
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
	// 기존 Temp 경로 내부 .xlsx 데이터를 삭제한다.
	tmp, err := userTemppath(ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = RemoveXLSX(tmp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = TEMPLATES.ExecuteTemplate(w, "importexcel", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleImportJSON 함수는 JSON 파일을 Import 하는 페이지 이다.
func handleImportJSON(w http.ResponseWriter, r *http.Request) {
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
		Setting     Setting
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
	// 기존 Temp 경로 내부 .json 데이터를 삭제한다.
	tmp, err := userTemppath(ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = RemoveJSON(tmp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "importjson", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleUploadExcel 핸들러는 Excel 파일을 받아 서버에 저장한다.
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	mimeType := header.Header.Get("Content-Type")
	switch mimeType {
	case "text/csv":
		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmp, err := userTemppath(ssid.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		path := tmp + "/" + header.Filename // 업로드한 파일 리스트를 불러오기 위해 뒤에 붙는 Unixtime을 제거한다.
		err = ioutil.WriteFile(path, data, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": // MS-Excel, Google & Libre Excel
		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmp, err := userTemppath(ssid.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		path := tmp + "/" + header.Filename // 업로드한 파일 리스트를 불러오기 위해 뒤에 붙는 Unixtime을 제거한다.
		err = ioutil.WriteFile(path, data, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("Not support: %s", mimeType), http.StatusInternalServerError) // 지원하지 않는 파일. 저장하지 않는다.
		return
	}
}

// handleUploadJSON 핸들러는 JSON 파일을 받아 서버에 저장한다.
func handleUploadJSON(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	mimeType := header.Header.Get("Content-Type")
	switch mimeType {
	case "application/json":
		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmp, err := userTemppath(ssid.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		path := tmp + "/" + header.Filename // 업로드한 파일 리스트를 불러오기 위해 뒤에 붙는 Unixtime을 제거한다.
		err = ioutil.WriteFile(path, data, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("Not support: %s", mimeType), http.StatusInternalServerError) // 지원하지 않는 파일. 저장하지 않는다.
		return
	}
}

// handleReportExcel 함수는 excel 파일을 체크하고 분석 보고서로 Redirection 한다.
func handleReportExcel(w http.ResponseWriter, r *http.Request) {
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
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	// 파일네임을 구한다.
	tmppath, err := userTemppath(ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// .xlsx 파일을 읽는다.
	xlsxs, err := GetXLSX(tmppath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(xlsxs) != 1 {
		http.Redirect(w, r, "/importexcel", http.StatusSeeOther)
		return
	}
	f, err := excelize.OpenFile(xlsxs[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type recipe struct {
		Project   string
		Filename  string
		Sheet     string
		Overwrite bool
		Rows      []Excelrow
		User
		SessionID string
		Devmode   bool
		SearchOption
		Errornum    int
		Projectlist []string
		Setting     Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.Sheet = "Sheet1"

	rcp.SessionID = ssid.ID
	rcp.Devmode = *flagDevmode
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

	var rows []Excelrow
	excelRows, err := f.GetRows(rcp.Sheet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(excelRows) == 0 {
		http.Error(w, rcp.Sheet+"값이 비어있습니다.", http.StatusBadRequest)
		return
	}
	for n, line := range excelRows {
		if n == 0 { // 첫번째줄
			if len(line) != 15 {
				http.Error(w, "약속된 Cell 갯수가 다릅니다", http.StatusBadRequest)
				return
			}
			continue
		}
		row := Excelrow{}
		row.Name, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("A%d", n+1)) // Name
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if row.Name == "" { // Name이 비어있다면 넘긴다.
			continue
		}

		row.Rnum, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("B%d", n+1)) // Rollnumber
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Shottype, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("C%d", n+1)) // Shottype(2d,3d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Note, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("D%d", n+1)) // 작업내용
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Comment, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("E%d", n+1)) // 수정사항
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Tags, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("F%d", n+1)) // Tags
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Link, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("G%d", n+1)) // Source(제목:경로)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.JustTimecodeIn, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("H%d", n+1)) // JustTimecodeIn
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.JustTimecodeOut, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("I%d", n+1)) // JustTimecodeOut
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Ddline2D, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("J%d", n+1)) // 2D마감
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Ddline3D, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("K%d", n+1)) // 3D마감
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Findate, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("L%d", n+1)) // FIN날짜
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Finver, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("M%d", n+1)) // FIN버전
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		row.HandleIn, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("N%d", n+1)) // 핸들IN
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.HandleOut, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("O%d", n+1)) // 핸들OUT
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		row.checkerror(session, project)
		rcp.Errornum += row.Errornum
		rows = append(rows, row)
	}

	rcp.Rows = rows
	err = TEMPLATES.ExecuteTemplate(w, "reportexcel", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// handleReportJSON 함수는 json 파일을 체크하고 분석 보고서로 Redirection 한다.
func handleReportJSON(w http.ResponseWriter, r *http.Request) {
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
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	// 파일네임을 구한다.
	tmppath, err := userTemppath(ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// .json 파일을 읽는다.
	jsons, err := GetJSON(tmppath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(jsons) != 1 {
		http.Redirect(w, r, "/importjson", http.StatusSeeOther)
		return
	}
	jsonFile, err := ioutil.ReadFile(jsons[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 파일이 정상인지 체크한다.
	type recipe struct {
		Project   string
		Filename  string
		Overwrite bool
		Rows      []Item
		User
		SessionID string
		Devmode   bool
		SearchOption
		Projectlist []string
		Setting     Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.Project = project
	rcp.SessionID = ssid.ID
	rcp.Devmode = *flagDevmode
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

	var rows []Item
	err = json.Unmarshal(jsonFile, &rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(rows) == 0 {
		http.Error(w, "json 값이 비어있습니다.", http.StatusBadRequest)
		return
	}
	rcp.Rows = rows
	err = TEMPLATES.ExecuteTemplate(w, "reportjson", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleExcelSubmit 함수는 excel 파일을 전송한다.
func handleExcelSubmit(w http.ResponseWriter, r *http.Request) {
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
	// 사용자의 이름을 구한다.
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized) // 사용자가 존재하지 않으면 당연히 Comment를 작성하면 안된다.
		return
	}
	authorName := u.LastNameKor + u.FirstNameKor
	// 파일네임을 구한다.
	tmppath, err := userTemppath(ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 로그 기록을 위해서 host 값을 구한다.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// .xlsx 파일을 읽는다.
	xlsxs, err := GetXLSX(tmppath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(xlsxs) != 1 {
		http.Redirect(w, r, "/importexcel", http.StatusSeeOther)
		return
	}
	f, err := excelize.OpenFile(xlsxs[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type ErrorItem struct {
		Name  string
		Error string
	}
	type recipe struct {
		Filename string
		Sheet    string
		User
		SessionID string
		Devmode   bool
		SearchOption
		ErrorItems []ErrorItem
		Setting    Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.SessionID = ssid.ID
	rcp.Devmode = *flagDevmode
	rcp.SearchOption = handleRequestToSearchOption(r)
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Sheet = "Sheet1"
	project := r.FormValue("project")
	overwrite := str2bool(r.FormValue("overwrite"))

	excelRows, err := f.GetRows(rcp.Sheet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(excelRows) == 0 {
		http.Error(w, rcp.Sheet+"값이 비어있습니다.", http.StatusBadRequest)
		return
	}
	// 로그 처리시 로그서버에는 로그를 기록하지만, 대량이 들어갈 때 slack에는 전달하지 않습니다.
	// slack에 대량의 로그가 쌓이는것을 원치않기 때문입니다.

	for n, line := range excelRows {
		if n == 0 { // 첫번째줄
			if len(line) != 15 {
				http.Error(w, "약속된 Cell 갯수가 다릅니다", http.StatusBadRequest)
				return
			}
			continue
		}
		name, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("A%d", n+1)) // Name
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if name == "" { // 샷이름이 없다면 넘긴다.
			continue
		}
		typ, err := Type(session, project, name)
		if err != nil {
			continue // 샷 타입을 가지고 올 수 없다면 넘긴다.
		}
		// 롤넘버
		rnum, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("B%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if rnum != "" {
			err := SetRnum(session, project, name+"_"+typ, rnum)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, "Set Rnum: "+rnum, project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// Shottype 2d,3d
		shottype, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("C%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if shottype != "" {
			_, err := SetShotType(session, project, name, shottype)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Shottype: %s", shottype), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// 작업내용
		note, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("D%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if note != "" {
			itemName, _, err := SetNote(session, project, name+"_"+typ, ssid.ID, note, overwrite)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: itemName, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, "Set Note: "+note, project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// 수정사항
		comment, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("E%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if comment != "" {
			_, err = AddComment(session, project, name, ssid.ID, authorName, time.Now().Format(time.RFC3339), comment, "", "")
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Add Comment: %s, Media: %s", comment, ""), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// Tags
		tags, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("F%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if tags != "" {
			for _, tag := range strings.Split(tags, ",") {
				removeSpaceTag := strings.Replace(tag, " ", "", -1) // Tag에 존재하는 띄어쓰기를 제거한다.
				if !regexpTag.MatchString(removeSpaceTag) {
					rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: "tag에는 특수문자를 사용할 수 없습니다"})
					continue
				}
				_, err = AddTag(session, project, name+"_"+typ, removeSpaceTag)
				if err != nil {
					rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
					continue
				}
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Add Tag: %s", tags), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// Source(제목:경로)
		sources, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("G%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if sources != "" {
			for _, s := range strings.Split(sources, "\n") {
				source := strings.Split(s, ":")
				title := strings.TrimSpace(source[0])
				path := strings.TrimSpace(source[1])
				_, err = AddSource(session, project, name, ssid.ID, title, path)
				if err != nil {
					rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
					continue
				}
				// log
				err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Add Source: %s, %s", title, path), rcp.Project, name, "csi3", ssid.ID, 180)
				if err != nil {
					rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
					continue
				}
			}

		}
		// JustTimecodeIn
		justTimecodeIn, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("H%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if justTimecodeIn != "" {
			err = SetJustTimecodeIn(session, project, name, justTimecodeIn)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("JustTimecodeIn: %s", justTimecodeIn), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// JustTimecoeOut
		justTimecodeOut, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("I%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if justTimecodeOut != "" {
			err = SetJustTimecodeOut(session, project, name, justTimecodeOut)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("JustTimecodeOut: %s", justTimecodeOut), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// 2D마감
		ddline2d, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("J%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ddline2d != "" {
			date, err := ditime.ToFullTime(19, ddline2d)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			_, err = SetDeadline2D(session, project, name, date)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Set Deadline2D: %s", date), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// 3D마감
		ddline3d, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("K%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if ddline3d != "" {
			date, err := ditime.ToFullTime(19, ddline3d)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			_, err = SetDeadline3D(session, project, name, date)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Set Deadline3D: %s", date), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// FIN날짜
		findate, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("L%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if findate != "" {
			date, err := ditime.ToFullTime(19, findate)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = SetFindate(session, project, name, date)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Set FinDate: %s", date), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// FIN버전
		finver, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("M%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if finver != "" {
			err = SetFinver(session, project, name, finver)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Set FinVersion: %s", finver), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// HandleIn
		handleIn, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("N%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if handleIn != "" {
			num, err := strconv.Atoi(handleIn)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = SetFrame(session, project, name, "handlein", num)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Handle In: %d", num), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		// HandleOut
		handleOut, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("O%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if handleOut != "" {
			num, err := strconv.Atoi(handleOut)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = SetFrame(session, project, name, "handleout", num)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = dilog.Add(*flagDBIP, host, fmt.Sprintf("Handle Out: %d", num), project, name, "csi3", ssid.ID, 180)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
	}
	err = TEMPLATES.ExecuteTemplate(w, "resultexcel", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// handleJSONSubmit 함수는 json 파일을 전송한다.
func handleJSONSubmit(w http.ResponseWriter, r *http.Request) {
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
	// 파일네임을 구한다.
	tmppath, err := userTemppath(ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 로그 기록을 위해서 host 값을 구한다.
	_, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// .xlsx 파일을 읽는다.
	jsonFiles, err := GetJSON(tmppath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(jsonFiles) != 1 {
		http.Redirect(w, r, "/importexcel", http.StatusSeeOther)
		return
	}
	jsonFile, err := ioutil.ReadFile(jsonFiles[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type recipe struct {
		Filename string
		User
		SessionID string
		Devmode   bool
		SearchOption
		Setting Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.SessionID = ssid.ID
	rcp.Devmode = *flagDevmode
	rcp.SearchOption = handleRequestToSearchOption(r)
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	project := r.FormValue("project")
	overwrite := str2bool(r.FormValue("overwrite"))
	var rows []Item
	err = json.Unmarshal(jsonFile, &rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(rows) == 0 {
		http.Error(w, "json 값이 비어있습니다.", http.StatusBadRequest)
		return
	}
	for _, i := range rows {
		if overwrite {
			err = setItem(session, project, i) // 기존데이터를 덮어쓰기 한다.
			if err != nil && err == mgo.ErrNotFound {
				// 새로운 데이터를 추가한다.
				err = addItem(session, project, i)
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Println(err)
			}
		} else {
			// 새로운 데이터를 추가한다.
			err = addItem(session, project, i)
			if err != nil {
				log.Println(err)
			}
		}
	}
	err = TEMPLATES.ExecuteTemplate(w, "resultjson", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleExportExcel 함수는 엑셀파일을 Export 하는 페이지 이다.
func handleExportExcel(w http.ResponseWriter, r *http.Request) {
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
		Projectlist []string
		SessionID   string
		Devmode     bool
		Setting     Setting
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
	// 만약 사용자에게 AccessProjects가 설정되어있다면 해당프로젝트만 보여야한다.
	if len(rcp.User.AccessProjects) != 0 {
		var accessProjects []string
		for _, i := range rcp.Projectlist {
			for _, j := range rcp.User.AccessProjects {
				if i != j {
					continue
				}
				accessProjects = append(accessProjects, i)
			}
		}
		rcp.Projectlist = accessProjects
	}

	err = TEMPLATES.ExecuteTemplate(w, "exportexcel", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleExportJSON 함수는 엑셀파일을 Export 하는 페이지 이다.
func handleExportJSON(w http.ResponseWriter, r *http.Request) {
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
		Projectlist []string
		SessionID   string
		Devmode     bool
		Setting     Setting
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
	// 만약 사용자에게 AccessProjects가 설정되어있다면 해당프로젝트만 보여야한다.
	if len(rcp.User.AccessProjects) != 0 {
		var accessProjects []string
		for _, i := range rcp.Projectlist {
			for _, j := range rcp.User.AccessProjects {
				if i != j {
					continue
				}
				accessProjects = append(accessProjects, i)
			}
		}
		rcp.Projectlist = accessProjects
	}

	err = TEMPLATES.ExecuteTemplate(w, "exportjson", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleExportExcelSubmit 함수는 export excel을 처리한다.
func handleExportExcelSubmit(w http.ResponseWriter, r *http.Request) {
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
	format := r.FormValue("format")
	sortkey := r.FormValue("sortkey")
	task := r.FormValue("task")
	statusv2 := str2bool(r.FormValue("statusv2"))

	var searchItems []Item
	switch format {
	case "shot":
		searchItems, err = SearchAllShot(session, project, sortkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "asset":
		searchItems, err = SearchAllAsset(session, project, sortkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "", "all":
		searchItems, err = SearchAll(session, project, sortkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// task가 선택되어 있다면 item을 돌면서 item을 거른다.
	var items []Item
	if task == "all" {
		items = searchItems
	} else {
		for _, i := range searchItems {
			if _, found := i.Tasks[task]; !found {
				continue
			}
			items = append(items, i)
		}
	}
	// status에 필요한 컬러를 불러온다.
	bgcolor := make(map[string]string)
	textcolor := make(map[string]string)
	status, err := AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, s := range status {
		bgcolor[s.ID] = s.BGColor
		textcolor[s.ID] = s.TextColor
	}
	bgcolor[""] = "#FFFFFF"   // Default BG color
	textcolor[""] = "#000000" // Default Text color
	bgcolor["0"] = "#3D3B3B"  // None, legacy
	bgcolor["1"] = "#606161"  // Hold, legacy
	bgcolor["2"] = "#E4D2B7"  // Done, legacy
	bgcolor["3"] = "#EEA4F1"  // Out, legacy
	bgcolor["4"] = "#FFF76B"  // Assign, legacy
	bgcolor["5"] = "#BEEF37"  // Ready, legacy
	bgcolor["6"] = "#77BB40"  // Wip, legacy
	bgcolor["7"] = "#54D6FD"  // Confirm, legacy
	bgcolor["8"] = "#FC9F55"  // Omit, legacy
	bgcolor["9"] = "#FFFFFF"  // Client, legacy

	f := excelize.NewFile()
	sheet := "Sheet1"
	index := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	// 스타일
	style, err := f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center","wrap_text":true}}`)
	if err != nil {
		log.Println(err)
	}
	textStyle, err := f.NewStyle(`{"alignment":{"horizontal":"left","vertical":"top", "wrap_text":true}}`)
	if err != nil {
		log.Println(err)
	}
	// 제목생성
	titles := []string{
		"Name",
		"Type",
		"Rollnumber",
		"Thumbnail",
		"ShotType(2d/3d)",
		"UseType(재스캔사용버전)",
		"상태",
		"작업내용",
		"수정사항",
		"Tags",
		"JustTimecodeIn",
		"JustTimecodeOut",
		"ScanTimecodeIn",
		"ScanTimecodeOut",
		"2D마감",
		"3D마감",
	}
	tasks, err := TasksettingNamesByExcelOrder(session)
	if err != nil {
		log.Println(err)
	}
	titles = append(titles, tasks...)
	for n, i := range titles {
		pos, err := excelize.CoordinatesToCellName(1+n, 1)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i)
		colName, err := excelize.ColumnNumberToName(n + 1)
		if err != nil {
			log.Println(err)
		}
		f.SetColWidth(sheet, colName, colName, 20)
		f.SetCellStyle(sheet, pos, pos, style)
	}

	// 엑셀파일 생성
	for n, i := range items {
		f.SetRowHeight(sheet, n+2, 60)
		// 이름
		pos, err := excelize.CoordinatesToCellName(1, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Name)
		f.SetCellStyle(sheet, pos, pos, style)
		// Type
		pos, err = excelize.CoordinatesToCellName(2, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Type)
		f.SetCellStyle(sheet, pos, pos, style)
		// 롤넘버
		pos, err = excelize.CoordinatesToCellName(3, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Rnum)
		f.SetCellStyle(sheet, pos, pos, style)
		// 썸네일
		pos, err = excelize.CoordinatesToCellName(4, n+2)
		if err != nil {
			log.Println(err)
		}
		imgPath := fmt.Sprintf("%s/%s/%s.jpg", CachedAdminSetting.ThumbnailRootPath, project, i.ID)
		f.AddPicture(sheet, pos, imgPath, `{"x_offset": 1, "y_offset": 1, "x_scale": 0.359, "y_scale": 0.359, "print_obj": true, "lock_aspect_ratio": true, "locked": true}`)
		// ShotType
		pos, err = excelize.CoordinatesToCellName(5, n+2)
		if err != nil {
			log.Println(err)
		}
		if i.Type == "asset" {
			f.SetCellValue(sheet, pos, strings.ToUpper(i.Assettype))
		} else {
			f.SetCellValue(sheet, pos, strings.ToUpper(i.Shottype))
		}
		f.SetCellStyle(sheet, pos, pos, style)
		// UseType(재스캔사용버전)
		pos, err = excelize.CoordinatesToCellName(6, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.UseType)
		f.SetCellStyle(sheet, pos, pos, style)
		// 상태
		pos, err = excelize.CoordinatesToCellName(7, n+2)
		if err != nil {
			log.Println(err)
		}
		// 기존에 Static한 상태의 컬러를 사용한다. // legacy
		statusStyle, err := f.NewStyle(
			fmt.Sprintf(`{
				"alignment":{"horizontal":"center","vertical":"center"},
				"fill":{"type":"pattern","color":["%s"],"pattern":1},
				"border":[
					{"type":"left","color":"888888","style":1},
					{"type":"top","color":"888888","style":1},
					{"type":"bottom","color":"888888","style":1},
					{"type":"right","color":"888888","style":1}]
				}`, bgcolor[i.Status]))
		if err != nil {
			log.Println(err)
		}
		// 다이나나믹 Status가 설정되어 있다면 해당 컬러를 사용한다.
		if statusv2 {
			statusStyle, err = f.NewStyle(
				fmt.Sprintf(`{
					"alignment":{"horizontal":"center","vertical":"center"},
					"font":{"color":"%s"},
					"fill":{"type":"pattern","color":["%s"],"pattern":1},
					"border":[
						{"type":"left","color":"888888","style":1},
						{"type":"top","color":"888888","style":1},
						{"type":"bottom","color":"888888","style":1},
						{"type":"right","color":"888888","style":1}]
					}`, textcolor[i.StatusV2], bgcolor[i.StatusV2]))
			if err != nil {
				log.Println(err)
			}
		}
		if statusv2 {
			f.SetCellValue(sheet, pos, i.StatusV2)
		} else {
			f.SetCellValue(sheet, pos, Status2capString(i.Status)) // legacy
		}

		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// 작업내용
		pos, err = excelize.CoordinatesToCellName(8, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Note.Text)
		f.SetCellStyle(sheet, pos, pos, textStyle)
		// 수정사항
		pos, err = excelize.CoordinatesToCellName(9, n+2)
		if err != nil {
			log.Println(err)
		}
		comments := []string{}
		for _, c := range ReverseCommentSlice(i.Comments) {
			comments = append(comments, c.Text)
		}
		f.SetCellValue(sheet, pos, strings.Join(comments, "\n"))
		f.SetCellStyle(sheet, pos, pos, textStyle)
		// Tags
		pos, err = excelize.CoordinatesToCellName(10, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, strings.Join(i.Tag, ","))

		f.SetCellStyle(sheet, pos, pos, style)
		// JustTimecodeIn
		pos, err = excelize.CoordinatesToCellName(11, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.JustTimecodeIn)
		f.SetCellStyle(sheet, pos, pos, style)
		// JustTimecodeOut
		pos, err = excelize.CoordinatesToCellName(12, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.JustTimecodeOut)
		f.SetCellStyle(sheet, pos, pos, style)
		// ScanTimecodeIn
		pos, err = excelize.CoordinatesToCellName(13, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.ScanTimecodeIn)
		f.SetCellStyle(sheet, pos, pos, style)
		// ScanTimecodeOut
		pos, err = excelize.CoordinatesToCellName(14, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.ScanTimecodeOut)
		f.SetCellStyle(sheet, pos, pos, style)
		// Deadline2D
		pos, err = excelize.CoordinatesToCellName(15, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, ToNormalTime(i.Ddline2d))
		f.SetCellStyle(sheet, pos, pos, style)
		// Deadline3D
		pos, err = excelize.CoordinatesToCellName(16, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, ToNormalTime(i.Ddline3d))
		f.SetCellStyle(sheet, pos, pos, style)
		// Tasks
		for taskOrder, t := range tasks {
			if _, found := i.Tasks[t]; !found {
				continue
			}
			pos, err = excelize.CoordinatesToCellName(17+taskOrder, n+2)
			if err != nil {
				log.Println(err)
			}
			// 기존 Static Status 일 때. // legacy
			statusStyle, err = f.NewStyle(
				fmt.Sprintf(`{
					"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
					"fill":{"type":"pattern","color":["%s"],"pattern":1},
					"border":[
						{"type":"left","color":"888888","style":1},
						{"type":"top","color":"888888","style":1},
						{"type":"bottom","color":"888888","style":1},
						{"type":"right","color":"888888","style":1}]
					}`, bgcolor[i.Tasks[t].Status]))
			if err != nil {
				log.Println(err)
			}
			// 만약 다이나믹 Status 일 때는 해당 컬러를 사용한다.
			if statusv2 {
				statusStyle, err = f.NewStyle(
					fmt.Sprintf(`{
						"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
						"font":{"color":"%s"},
						"fill":{"type":"pattern","color":["%s"],"pattern":1},
						"border":[
							{"type":"left","color":"888888","style":1},
							{"type":"top","color":"888888","style":1},
							{"type":"bottom","color":"888888","style":1},
							{"type":"right","color":"888888","style":1}]
						}`, textcolor[i.Tasks[t].StatusV2], bgcolor[i.Tasks[t].StatusV2]))
				if err != nil {
					log.Println(err)
				}
			}
			var text string
			if statusv2 {
				text = i.Tasks[t].StatusV2
			} else {
				text = Status2capString(i.Tasks[t].Status) // legacy
			}
			text += "\n" + i.Tasks[t].User
			text += "\n" + ToNormalTime(i.Tasks[t].Predate)
			text += "\n" + ToNormalTime(i.Tasks[t].Date)
			f.SetCellValue(sheet, pos, text)
			f.SetCellStyle(sheet, pos, pos, statusStyle)
		}
	}
	tempDir, err := ioutil.TempDir("", "excel")
	if err != nil {
		log.Println(err)
	}
	defer os.RemoveAll(tempDir)
	filename := format + ".xlsx"
	err = f.SaveAs(tempDir + "/" + filename)
	if err != nil {
		log.Println(err)
	}
	// 저장된 Excel 파일을 다운로드 시킨다.
	w.Header().Add("Content-Disposition", fmt.Sprintf("Attachment; filename=%s-%s-%s.xlsx", project, format, task))
	http.ServeFile(w, r, tempDir+"/"+filename)
}

// handleExportJSONSubmit 함수는 export json을 처리한다.
func handleExportJSONSubmit(w http.ResponseWriter, r *http.Request) {
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
	format := r.FormValue("format")
	sortkey := r.FormValue("sortkey")
	task := r.FormValue("task")

	var searchItems []Item
	switch format {
	case "shot":
		searchItems, err = SearchAllShot(session, project, sortkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "asset":
		searchItems, err = SearchAllAsset(session, project, sortkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "", "all":
		searchItems, err = SearchAll(session, project, sortkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// task가 선택되어 있다면 item을 돌면서 item을 거른다.
	var items []Item
	if task == "all" {
		items = searchItems
	} else {
		for _, i := range searchItems {
			if _, found := i.Tasks[task]; !found {
				continue
			}
			items = append(items, i)
		}
	}
	data, err := json.MarshalIndent(items, "", "    ") // json 파일이 보기 좋게 정렬되어 있어야 한다.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tempDir, err := ioutil.TempDir("", "json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)
	filename := format + ".json"
	err = ioutil.WriteFile(tempDir+"/"+filename, data, 0664)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 저장된 Json 파일을 다운로드 시킨다.
	w.Header().Add("Content-Disposition", fmt.Sprintf("Attachment; filename=%s-%s-%s.json", project, format, task))
	http.ServeFile(w, r, tempDir+"/"+filename)
}

// handleDownloadExcelFile 함수는 전송된 값을 이용해서 export excel을 처리한다.
func handleDownloadExcelFile(w http.ResponseWriter, r *http.Request) {
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
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	q := r.URL.Query()
	op := SearchOption{}
	project := q.Get("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	op.Project = project
	op.Task = q.Get("task")
	searchword := q.Get("searchword")
	if searchword == "" {
		http.Error(w, "검색어를 설정해주세요", http.StatusBadRequest)
		return
	}
	op.Searchword = searchword
	op.Sortkey = q.Get("sortkey")
	op.SearchbarTemplate = q.Get("searchbartemplate") // legacy
	op.Assign = str2bool(q.Get("assign"))             // legacy
	op.Ready = str2bool(q.Get("ready"))               // legacy
	op.Wip = str2bool(q.Get("wip"))                   // legacy
	op.Confirm = str2bool(q.Get("confirm"))           // legacy
	op.Done = str2bool(q.Get("done"))                 // legacy
	op.Omit = str2bool(q.Get("omit"))                 // legacy
	op.Hold = str2bool(q.Get("hold"))                 // legacy
	op.Out = str2bool(q.Get("out"))                   // legacy
	op.None = str2bool(q.Get("none"))                 // legacy
	op.TrueStatus = strings.Split(q.Get("truestatus"), ",")
	op.Shot = true
	op.Assets = true
	op.Type2d = true
	op.Type3d = true
	items, err := Search(session, op)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	statusv2 := false
	if op.SearchbarTemplate == "searchbarV2" {
		statusv2 = true
	}

	// status에 필요한 컬러를 불러온다.
	bgcolor := make(map[string]string)
	textcolor := make(map[string]string)
	status, err := AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, s := range status {
		bgcolor[s.ID] = s.BGColor
		textcolor[s.ID] = s.TextColor
	}
	bgcolor[""] = "#FFFFFF"   // Default BG color
	textcolor[""] = "#000000" // Default Text color
	bgcolor["0"] = "#3D3B3B"  // None, legacy
	bgcolor["1"] = "#606161"  // Hold, legacy
	bgcolor["2"] = "#E4D2B7"  // Done, legacy
	bgcolor["3"] = "#EEA4F1"  // Out, legacy
	bgcolor["4"] = "#FFF76B"  // Assign, legacy
	bgcolor["5"] = "#BEEF37"  // Ready, legacy
	bgcolor["6"] = "#77BB40"  // Wip, legacy
	bgcolor["7"] = "#54D6FD"  // Confirm, legacy
	bgcolor["8"] = "#FC9F55"  // Omit, legacy
	bgcolor["9"] = "#FFFFFF"  // Client, legacy

	f := excelize.NewFile()
	sheet := "Sheet1"
	index := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	// 스타일
	style, err := f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center","wrap_text":true}}`)
	if err != nil {
		log.Println(err)
	}
	textStyle, err := f.NewStyle(`{"alignment":{"horizontal":"left","vertical":"top", "wrap_text":true}}`)
	if err != nil {
		log.Println(err)
	}
	// 제목생성
	titles := []string{
		"Name",
		"Type",
		"Rollnumber",
		"Thumbnail",
		"ShotType(2d/3d)",
		"UseType(재스캔사용버전)",
		"상태",
		"작업내용",
		"수정사항",
		"Tags",
		"JustTimecodeIn",
		"JustTimecodeOut",
		"ScanTimecodeIn",
		"ScanTimecodeOut",
		"2D마감",
		"3D마감",
	}
	tasks, err := TasksettingNamesByExcelOrder(session)
	if err != nil {
		log.Println(err)
	}
	titles = append(titles, tasks...)
	for n, i := range titles {
		pos, err := excelize.CoordinatesToCellName(1+n, 1)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i)
		colName, err := excelize.ColumnNumberToName(n + 1)
		if err != nil {
			log.Println(err)
		}
		f.SetColWidth(sheet, colName, colName, 20)
		f.SetCellStyle(sheet, pos, pos, style)
	}

	// 엑셀파일 생성
	for n, i := range items {
		f.SetRowHeight(sheet, n+2, 60)
		// 이름
		pos, err := excelize.CoordinatesToCellName(1, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Name)
		f.SetCellStyle(sheet, pos, pos, style)
		// Type
		pos, err = excelize.CoordinatesToCellName(2, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Type)
		f.SetCellStyle(sheet, pos, pos, style)
		// 롤넘버
		pos, err = excelize.CoordinatesToCellName(3, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Rnum)
		f.SetCellStyle(sheet, pos, pos, style)
		// 썸네일
		pos, err = excelize.CoordinatesToCellName(4, n+2)
		if err != nil {
			log.Println(err)
		}
		imgPath := fmt.Sprintf("%s/%s/%s.jpg", CachedAdminSetting.ThumbnailRootPath, project, i.ID)
		f.AddPicture(sheet, pos, imgPath, `{"x_offset": 1, "y_offset": 1, "x_scale": 0.359, "y_scale": 0.359, "print_obj": true, "lock_aspect_ratio": true, "locked": true}`)
		// ShotType
		pos, err = excelize.CoordinatesToCellName(5, n+2)
		if err != nil {
			log.Println(err)
		}
		if i.Type == "asset" {
			f.SetCellValue(sheet, pos, strings.ToUpper(i.Assettype))
		} else {
			f.SetCellValue(sheet, pos, strings.ToUpper(i.Shottype))
		}
		f.SetCellStyle(sheet, pos, pos, style)
		// UseType(재스캔사용버전)
		pos, err = excelize.CoordinatesToCellName(6, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.UseType)
		f.SetCellStyle(sheet, pos, pos, style)
		// 상태
		pos, err = excelize.CoordinatesToCellName(7, n+2)
		if err != nil {
			log.Println(err)
		}
		// 기존에 Static한 상태의 컬러를 사용한다. // legacy
		statusStyle, err := f.NewStyle(
			fmt.Sprintf(`{
					"alignment":{"horizontal":"center","vertical":"center"},
					"fill":{"type":"pattern","color":["%s"],"pattern":1},
					"border":[
						{"type":"left","color":"888888","style":1},
						{"type":"top","color":"888888","style":1},
						{"type":"bottom","color":"888888","style":1},
						{"type":"right","color":"888888","style":1}]
					}`, bgcolor[i.Status]))
		if err != nil {
			log.Println(err)
		}
		// 다이나나믹 Status가 설정되어 있다면 해당 컬러를 사용한다.
		if statusv2 {
			statusStyle, err = f.NewStyle(
				fmt.Sprintf(`{
						"alignment":{"horizontal":"center","vertical":"center"},
						"font":{"color":"%s"},
						"fill":{"type":"pattern","color":["%s"],"pattern":1},
						"border":[
							{"type":"left","color":"888888","style":1},
							{"type":"top","color":"888888","style":1},
							{"type":"bottom","color":"888888","style":1},
							{"type":"right","color":"888888","style":1}]
						}`, textcolor[i.StatusV2], bgcolor[i.StatusV2]))
			if err != nil {
				log.Println(err)
			}
		}
		if statusv2 {
			f.SetCellValue(sheet, pos, i.StatusV2)
		} else {
			f.SetCellValue(sheet, pos, Status2capString(i.Status)) // legacy
		}

		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// 작업내용
		pos, err = excelize.CoordinatesToCellName(8, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Note.Text)
		f.SetCellStyle(sheet, pos, pos, textStyle)
		// 수정사항
		pos, err = excelize.CoordinatesToCellName(9, n+2)
		if err != nil {
			log.Println(err)
		}
		comments := []string{}
		for _, c := range ReverseCommentSlice(i.Comments) {
			comments = append(comments, c.Text)
		}
		f.SetCellValue(sheet, pos, strings.Join(comments, "\n"))
		f.SetCellStyle(sheet, pos, pos, textStyle)
		// Tags
		pos, err = excelize.CoordinatesToCellName(10, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, strings.Join(i.Tag, ","))

		f.SetCellStyle(sheet, pos, pos, style)
		// JustTimecodeIn
		pos, err = excelize.CoordinatesToCellName(11, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.JustTimecodeIn)
		f.SetCellStyle(sheet, pos, pos, style)
		// JustTimecodeOut
		pos, err = excelize.CoordinatesToCellName(12, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.JustTimecodeOut)
		f.SetCellStyle(sheet, pos, pos, style)
		// ScanTimecodeIn
		pos, err = excelize.CoordinatesToCellName(13, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.ScanTimecodeIn)
		f.SetCellStyle(sheet, pos, pos, style)
		// ScanTimecodeOut
		pos, err = excelize.CoordinatesToCellName(14, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.ScanTimecodeOut)
		f.SetCellStyle(sheet, pos, pos, style)
		// Deadline2D
		pos, err = excelize.CoordinatesToCellName(15, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, ToNormalTime(i.Ddline2d))
		f.SetCellStyle(sheet, pos, pos, style)
		// Deadline3D
		pos, err = excelize.CoordinatesToCellName(16, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, ToNormalTime(i.Ddline3d))
		f.SetCellStyle(sheet, pos, pos, style)
		// Tasks
		for taskOrder, t := range tasks {
			if _, found := i.Tasks[t]; !found {
				continue
			}
			pos, err = excelize.CoordinatesToCellName(17+taskOrder, n+2)
			if err != nil {
				log.Println(err)
			}
			// 기존 Static Status 일 때. // legacy
			statusStyle, err = f.NewStyle(
				fmt.Sprintf(`{
						"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
						"fill":{"type":"pattern","color":["%s"],"pattern":1},
						"border":[
							{"type":"left","color":"888888","style":1},
							{"type":"top","color":"888888","style":1},
							{"type":"bottom","color":"888888","style":1},
							{"type":"right","color":"888888","style":1}]
						}`, bgcolor[i.Tasks[t].Status]))
			if err != nil {
				log.Println(err)
			}
			// 만약 다이나믹 Status 일 때는 해당 컬러를 사용한다.
			if statusv2 {
				statusStyle, err = f.NewStyle(
					fmt.Sprintf(`{
							"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
							"font":{"color":"%s"},
							"fill":{"type":"pattern","color":["%s"],"pattern":1},
							"border":[
								{"type":"left","color":"888888","style":1},
								{"type":"top","color":"888888","style":1},
								{"type":"bottom","color":"888888","style":1},
								{"type":"right","color":"888888","style":1}]
							}`, textcolor[i.Tasks[t].StatusV2], bgcolor[i.Tasks[t].StatusV2]))
				if err != nil {
					log.Println(err)
				}
			}
			var text string
			if statusv2 {
				text = i.Tasks[t].StatusV2
			} else {
				text = Status2capString(i.Tasks[t].Status) // legacy
			}
			text += "\n" + i.Tasks[t].User
			text += "\n" + ToNormalTime(i.Tasks[t].Predate)
			text += "\n" + ToNormalTime(i.Tasks[t].Date)
			f.SetCellValue(sheet, pos, text)
			f.SetCellStyle(sheet, pos, pos, statusStyle)
		}
	}
	tempDir, err := ioutil.TempDir("", "excel")
	if err != nil {
		log.Println(err)
	}
	defer os.RemoveAll(tempDir)
	filename := "currentPage.xlsx"
	err = f.SaveAs(tempDir + "/" + filename)
	if err != nil {
		log.Println(err)
	}
	// 저장된 Excel 파일을 다운로드 시킨다.
	w.Header().Add("Content-Disposition", fmt.Sprintf("Attachment; filename=%s-%s%s.xlsx", project, "currentPage", op.Task))
	http.ServeFile(w, r, tempDir+"/"+filename)
}

// handleDownloadJSONFile 함수는 전송된 값을 이용해서 export json을 처리한다.
func handleDownloadJSONFile(w http.ResponseWriter, r *http.Request) {
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
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	q := r.URL.Query()
	op := SearchOption{}
	project := q.Get("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	op.Project = project
	op.Task = q.Get("task")
	searchword := q.Get("searchword")
	if searchword == "" {
		http.Error(w, "검색어를 설정해주세요", http.StatusBadRequest)
		return
	}
	op.Searchword = searchword
	op.Sortkey = q.Get("sortkey")
	op.SearchbarTemplate = q.Get("searchbartemplate") // legacy
	op.Assign = str2bool(q.Get("assign"))             // legacy
	op.Ready = str2bool(q.Get("ready"))               // legacy
	op.Wip = str2bool(q.Get("wip"))                   // legacy
	op.Confirm = str2bool(q.Get("confirm"))           // legacy
	op.Done = str2bool(q.Get("done"))                 // legacy
	op.Omit = str2bool(q.Get("omit"))                 // legacy
	op.Hold = str2bool(q.Get("hold"))                 // legacy
	op.Out = str2bool(q.Get("out"))                   // legacy
	op.None = str2bool(q.Get("none"))                 // legacy
	op.TrueStatus = strings.Split(q.Get("truestatus"), ",")
	op.Shot = true
	op.Assets = true
	op.Type2d = true
	op.Type3d = true
	items, err := Search(session, op)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.MarshalIndent(items, "", "    ") // json 파일이 보기 좋게 정렬되어 있어야 한다.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tempDir, err := ioutil.TempDir("", "json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)
	filename := "currentPage.json"
	err = ioutil.WriteFile(tempDir+"/"+filename, data, 0664)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 저장된 Json 파일을 다운로드 시킨다.
	w.Header().Add("Content-Disposition", fmt.Sprintf("Attachment; filename=%s-%s%s.json", project, "currentPage", op.Task))
	http.ServeFile(w, r, tempDir+"/"+filename)
}

// handleDownloadCsvFile 함수는 전송된 값을 이용해서 export csv을 처리한다.
func handleDownloadCsvFile(w http.ResponseWriter, r *http.Request) {
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
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	q := r.URL.Query()
	op := SearchOption{}
	project := q.Get("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	op.Project = project
	op.Task = q.Get("task")
	searchword := q.Get("searchword")
	if searchword == "" {
		http.Error(w, "검색어를 설정해주세요", http.StatusBadRequest)
		return
	}
	op.Searchword = searchword
	op.Sortkey = q.Get("sortkey")
	op.SearchbarTemplate = q.Get("searchbartemplate") // legacy
	op.Assign = str2bool(q.Get("assign"))             // legacy
	op.Ready = str2bool(q.Get("ready"))               // legacy
	op.Wip = str2bool(q.Get("wip"))                   // legacy
	op.Confirm = str2bool(q.Get("confirm"))           // legacy
	op.Done = str2bool(q.Get("done"))                 // legacy
	op.Omit = str2bool(q.Get("omit"))                 // legacy
	op.Hold = str2bool(q.Get("hold"))                 // legacy
	op.Out = str2bool(q.Get("out"))                   // legacy
	op.None = str2bool(q.Get("none"))                 // legacy
	op.TrueStatus = strings.Split(q.Get("truestatus"), ",")
	op.Shot = true
	op.Assets = true
	op.Type2d = true
	op.Type3d = true
	items, err := Search(session, op)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	records := [][]string{}
	titles := strings.Split(q.Get("titles"), ",") // 타이틀을 설정한다.
	records = append(records, titles)
	// 데이터를 찾고 넣는다.
	for _, i := range items {
		datas := []string{}
		for _, title := range titles {
			switch title {
			// 혀용하지 않는 것(리스트형,맵등)
			case "Tag", "Assettags", "Note", "Sources", "References", "Comments", "Tasks", "OnsetCam", "ProductionCam":
				continue
			case "CrowdAsset":
				// bool형
				datas = append(datas, strconv.FormatBool(i.GetFieldBool(title)))
			case "ScanFrame", "ScanIn", "ScanOut", "HandleIn", "HandleOut", "JustIn", "JustOut", "PlateIn", "PlateOut", "ObjectidIn", "ObjectidOut":
				// int형
				datas = append(datas, fmt.Sprintf("%d", i.GetFieldInteger(title)))
			case "OverscanRatio":
				// float64 형태
				datas = append(datas, fmt.Sprintf("%f", i.GetFieldFloat64(title)))
			default:
				// 기본적으로 문자열로 처리
				datas = append(datas, i.GetFieldString(title))
			}
		}
		records = append(records, datas)
	}

	// csv를 만들기 위해서 tempdir을 생성한다.
	tempDir, err := ioutil.TempDir("", "csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)
	filename := "currentPage.csv"

	f, err := os.Create(tempDir + "/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	d := csv.NewWriter(f)
	err = d.WriteAll(records) // calls Flush internally
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Disposition", fmt.Sprintf("Attachment; filename=%s-%s%s.csv", project, "currentPage", op.Task))
	http.ServeFile(w, r, tempDir+"/"+filename)
}

// handleDownloadExcelTemplate 함수는 빈 export template을 다운로드 한다.
func handleDownloadExcelTemplate(w http.ResponseWriter, r *http.Request) {
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
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	f := excelize.NewFile()
	sheet := "Sheet1"
	index := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	// 스타일
	style, err := f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center","wrap_text":true}}`)
	if err != nil {
		log.Println(err)
	}
	// Template 제목생성
	f.SetCellValue(sheet, "A1", "Name")
	f.SetCellValue(sheet, "B1", "Rollnumber")
	f.SetCellValue(sheet, "C1", "Type(2d/3d)")
	f.SetCellValue(sheet, "D1", "작업내용")
	f.SetCellValue(sheet, "E1", "수정사항")
	f.SetCellValue(sheet, "F1", "Tags")
	f.SetCellValue(sheet, "G1", "Source(Key:Value)")
	f.SetCellValue(sheet, "H1", "JustTimecodeIn")
	f.SetCellValue(sheet, "I1", "JustTimecodeOut")
	f.SetCellValue(sheet, "J1", "2D마감")
	f.SetCellValue(sheet, "K1", "3D마감")
	f.SetCellValue(sheet, "L1", "Final Date")
	f.SetCellValue(sheet, "M1", "Final Version")
	f.SetCellValue(sheet, "N1", "Handle In")
	f.SetCellValue(sheet, "O1", "Handle Out")
	f.SetColWidth(sheet, "A", "O", 20)
	f.SetCellStyle(sheet, "A1", "O1", style)
	tempDir, err := ioutil.TempDir("", "excel")
	if err != nil {
		log.Println(err)
	}
	defer os.RemoveAll(tempDir)
	filename := "template.xlsx"
	err = f.SaveAs(tempDir + "/" + filename)
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("Content-Disposition", fmt.Sprintf("Attachment; filename=%s", filename))
	http.ServeFile(w, r, tempDir+"/"+filename)
}

// handleExportDumpProject 함수는 전송된 값을 이용해서 export json을 처리한다.
func handleExportDumpProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, accessLevel, err := TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if accessLevel < IoAccessLevel {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Temp 경로를 만든다.
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)
	dbName := "project"
	// Temp경로에 DB를 dump한다.
	args := []string{
		"-d",
		dbName,
		"-h",
		*flagDBIP,
		"--quiet",
		"--out",
		tempDir,
	}
	err = exec.Command(CachedAdminSetting.MongodumpPath, args...).Run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 압축 Tempfile을 생성한다.
	filename := "csi-dbdump.zip"
	zipDumpfile, err := ioutil.TempFile("", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer os.Remove(zipDumpfile.Name())
	// Temp 경로를 압축한다.
	err = ZipWriter(tempDir+"/"+dbName, zipDumpfile.Name())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 저장된 Json 파일을 다운로드 시킨다.
	w.Header().Add("Content-Disposition", fmt.Sprintf("Attachment; filename=%s", filename))
	http.ServeFile(w, r, zipDumpfile.Name())
}
