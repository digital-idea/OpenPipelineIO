package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
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
	case "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": // MS-Excel, Google & Libre Excel
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
		log.Printf("Not support: %s", mimeType)
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
	}
	rcp := recipe{}
	rcp.Sheet = "Sheet1"
	rcp.Projectlist, err = OnProjectlist(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.SessionID = ssid.ID
	rcp.Devmode = *flagDevmode
	rcp.SearchOption = handleRequestToSearchOption(r)
	rcp.User, err = getUser(session, ssid.ID)
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

		row.checkerror() // 각 값을 에러체크한다.
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// 롤넘버
		rnum, err := f.GetCellValue(rcp.Sheet, fmt.Sprintf("B%d", n+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if rnum != "" {
			_, err := SetRnum(session, project, name, rnum)
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
			_, err = AddComment(session, project, name, ssid.ID, time.Now().Format(time.RFC3339), comment, "")
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
				_, err = AddTag(session, project, name, tag)
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
	var items []Item
	if format == "cglist-shot" {
		items, err = SearchAllShot(session, project, sortkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if format == "cglist-asset" {
		items, err = SearchAllAsset(session, project, sortkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		items, err = SearchAll(session, project, sortkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
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
		"Rollnumber",
		"Thumbnail",
		"Type(2d/3d)",
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
	tasks, err := TasksettingNames(session)
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
		// 롤넘버
		pos, err = excelize.CoordinatesToCellName(2, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Rnum)
		f.SetCellStyle(sheet, pos, pos, style)
		// 썸네일
		pos, err = excelize.CoordinatesToCellName(3, n+2)
		if err != nil {
			log.Println(err)
		}
		imgPath := fmt.Sprintf("%s/%s/%s.jpg", *flagThumbPath, project, i.ID)
		f.AddPicture(sheet, pos, imgPath, `{"x_offset": 1, "y_offset": 1, "x_scale": 0.359, "y_scale": 0.359, "print_obj": true, "lock_aspect_ratio": true, "locked": true}`)
		// Type
		pos, err = excelize.CoordinatesToCellName(4, n+2)
		if err != nil {
			log.Println(err)
		}
		if i.Type == "asset" {
			f.SetCellValue(sheet, pos, strings.ToUpper(i.Assettype))
		} else {
			f.SetCellValue(sheet, pos, strings.ToUpper(i.Shottype))
		}
		f.SetCellStyle(sheet, pos, pos, style)
		// 상태
		pos, err = excelize.CoordinatesToCellName(5, n+2)
		if err != nil {
			log.Println(err)
		}
		statusStyle, err := f.NewStyle(
			fmt.Sprintf(`
			{"alignment":{"horizontal":"center","vertical":"center"},
			"fill":{"type":"pattern","color":["%s"],"pattern":1},
			"border":[
				{"type":"left","color":"888888","style":1},
				{"type":"top","color":"888888","style":1},
				{"type":"bottom","color":"888888","style":1},
				{"type":"right","color":"888888","style":1}]
				}`, itemStatus2color(i.Status)))
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, Status2capString(i.Status))
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// 작업내용
		pos, err = excelize.CoordinatesToCellName(6, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Note.Text)
		f.SetCellStyle(sheet, pos, pos, textStyle)
		// 수정사항
		pos, err = excelize.CoordinatesToCellName(7, n+2)
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
		pos, err = excelize.CoordinatesToCellName(8, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, strings.Join(i.Tag, ","))

		f.SetCellStyle(sheet, pos, pos, style)
		// JustTimecodeIn
		pos, err = excelize.CoordinatesToCellName(9, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.JustTimecodeIn)
		f.SetCellStyle(sheet, pos, pos, style)
		// JustTimecodeOut
		pos, err = excelize.CoordinatesToCellName(10, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.JustTimecodeOut)
		f.SetCellStyle(sheet, pos, pos, style)
		// ScanTimecodeIn
		pos, err = excelize.CoordinatesToCellName(11, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.ScanTimecodeIn)
		f.SetCellStyle(sheet, pos, pos, style)
		// ScanTimecodeOut
		pos, err = excelize.CoordinatesToCellName(12, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.ScanTimecodeOut)
		f.SetCellStyle(sheet, pos, pos, style)
		// Deadline2D
		pos, err = excelize.CoordinatesToCellName(13, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, ToNormalTime(i.Ddline2d))
		f.SetCellStyle(sheet, pos, pos, style)
		// Deadline3D
		pos, err = excelize.CoordinatesToCellName(14, n+2)
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
			pos, err = excelize.CoordinatesToCellName(15+taskOrder, n+2)
			if err != nil {
				log.Println(err)
			}
			statusStyle, err = f.NewStyle(
				fmt.Sprintf(`
				{"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
				"fill":{"type":"pattern","color":["%s"],"pattern":1},
				"border":[
					{"type":"left","color":"888888","style":1},
					{"type":"top","color":"888888","style":1},
					{"type":"bottom","color":"888888","style":1},
					{"type":"right","color":"888888","style":1}]
					}`, itemStatus2color(i.Tasks[t].Status)))
			if err != nil {
				log.Println(err)
			}
			text := Status2capString(i.Tasks[t].Status)
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
	w.Header().Add("Content-Disposition", fmt.Sprintf("Attachment; filename=%s-%s.xlsx", project, format))
	http.ServeFile(w, r, tempDir+"/"+filename)
}

// handleDownloadExcelTemplate 함수는 export template을 다운로드 한다.
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
