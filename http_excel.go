package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
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
		row.Shottype, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("B%d", n+1)) // shottype 2d,3d
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Note, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("C%d", n+1)) // 작업내용
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Comment, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("D%d", n+1)) // 수정사항
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Link, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("E%d", n+1)) // 링크자료(제목:경로)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Ddline3D, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("F%d", n+1)) // 3D마감
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Ddline2D, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("G%d", n+1)) // 2D마감
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Findate, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("H%d", n+1)) // FIN날짜
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Finver, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("I%d", n+1)) // FIN버전
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Tags, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("J%d", n+1)) // 태그
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.Rnum, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("K%d", n+1)) // 롤넘버
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.HandleIn, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("L%d", n+1)) // 핸들IN
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.HandleOut, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("M%d", n+1)) // 핸들OUT
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.JustTimecodeIn, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("N%d", n+1)) // JUST타임코드IN
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.JustTimecodeOut, err = f.GetCellValue(rcp.Sheet, fmt.Sprintf("O%d", n+1)) // JUST타임코드OUT
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row.checkerror() // 각 값을 에러체크한다.
		rcp.Errornum += row.Errornum
		rows = append(rows, row)
	}
	rcp.SessionID = ssid.ID
	rcp.Devmode = *flagDevmode
	rcp.SearchOption = handleRequestToSearchOption(r)
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	for n, line := range excelRows {
		if n == 0 { // 첫번째줄
			if len(line) != 15 {
				http.Error(w, "약속된 Cell 갯수가 다릅니다", http.StatusBadRequest)
				return
			}
			continue
		}
		name, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("A%d", n+1))     // Name
		shottype, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("B%d", n+1)) // shottype 2d,3d
		if shottype != "" {
			err := SetShotType(session, project, name, shottype)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		note, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("C%d", n+1)) // 작업내용
		if note != "" {
			_, err = SetNote(session, project, name, ssid.ID, note, overwrite)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		comment, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("D%d", n+1)) // 수정사항
		if comment != "" {
			err = AddComment(session, project, name, ssid.ID, time.Now().Format(time.RFC3339), comment)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		sources, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("E%d", n+1)) // 링크자료(제목:경로)
		if sources != "" {
			for _, s := range strings.Split(sources, "\n") {
				source := strings.Split(s, ":")
				title := strings.TrimSpace(source[0])
				path := strings.TrimSpace(source[1])
				err = AddSource(session, project, name, ssid.ID, title, path)
				if err != nil {
					rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
					continue
				}
			}
		}
		ddline3d, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("F%d", n+1)) // 3D마감
		if ddline3d != "" {
			date, err := ditime.ToFullTime(19, ddline3d)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = SetDeadline3D(session, project, name, date)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		ddline2d, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("G%d", n+1)) // 2D마감
		if ddline2d != "" {
			date, err := ditime.ToFullTime(19, ddline2d)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
			err = SetDeadline2D(session, project, name, date)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		findate, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("H%d", n+1)) // FIN날짜
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
		}
		finver, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("I%d", n+1)) // FIN버전
		if finver != "" {
			err = SetFinver(session, project, name, finver)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}

		tags, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("J%d", n+1)) // 태그
		if tags != "" {
			err = SetTags(session, project, name, strings.Split(tags, ","))
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		rnum, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("K%d", n+1)) // 롤넘버
		if rnum != "" {
			err = SetRnum(session, project, name, rnum)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		handleIn, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("L%d", n+1)) // 핸들IN
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
		}
		handleOut, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("M%d", n+1)) // 핸들OUT
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
		}
		justTimecodeIn, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("N%d", n+1)) // JUST타임코드IN
		if justTimecodeIn != "" {
			err = SetJustTimecodeIn(session, project, name, justTimecodeIn)
			if err != nil {
				rcp.ErrorItems = append(rcp.ErrorItems, ErrorItem{Name: name, Error: err.Error()})
				continue
			}
		}
		justTimecodeOut, _ := f.GetCellValue(rcp.Sheet, fmt.Sprintf("O%d", n+1)) // JUST타임코드OUT
		if justTimecodeOut != "" {
			err = SetJustTimecodeIn(session, project, name, justTimecodeOut)
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
	rcp.Projectlist, err = OnProjectlist(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User, err = getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	var items []Item
	if format == "cglist-shot" {
		items, err = SearchAllShot(session, project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if format == "cglist-asset" {
		items, err = SearchAllAsset(session, project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		items, err = SearchAll(session, project)
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
	f.SetCellValue(sheet, "A1", "Name")
	f.SetCellValue(sheet, "B1", "Thumbnail")
	f.SetCellValue(sheet, "C1", "롤넘버")
	f.SetCellValue(sheet, "D1", "Type")
	f.SetCellValue(sheet, "E1", "상태")
	f.SetCellValue(sheet, "F1", "작업내용")
	f.SetCellValue(sheet, "G1", "태그")
	f.SetCellValue(sheet, "H1", "수정사항")
	f.SetCellValue(sheet, "I1", "JustTimecodeIn")
	f.SetCellValue(sheet, "J1", "JustTimecodeOut")
	f.SetCellValue(sheet, "K1", "2D마감")
	f.SetCellValue(sheet, "L1", "comp")
	f.SetCellValue(sheet, "M1", "matte")
	f.SetCellValue(sheet, "N1", "mg")
	f.SetCellValue(sheet, "O1", "3D마감")
	f.SetCellValue(sheet, "P1", "model")
	f.SetCellValue(sheet, "Q1", "mm")
	f.SetCellValue(sheet, "R1", "layout")
	f.SetCellValue(sheet, "S1", "ani")
	f.SetCellValue(sheet, "T1", "fur")
	f.SetCellValue(sheet, "U1", "sim")
	f.SetCellValue(sheet, "V1", "crowd")
	f.SetCellValue(sheet, "W1", "fx")
	f.SetCellValue(sheet, "X1", "light")
	f.SetCellValue(sheet, "Y1", "previz")
	f.SetColWidth(sheet, "A", "Y", 20)
	f.SetCellStyle(sheet, "A1", "Y1", style)
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
		// 썸네일
		pos, err = excelize.CoordinatesToCellName(2, n+2)
		if err != nil {
			log.Println(err)
		}
		imgPath := fmt.Sprintf("%s/%s/%s.jpg", *flagThumbPath, project, i.ID)
		f.AddPicture(sheet, pos, imgPath, `{"x_offset": 1, "y_offset": 1, "x_scale": 0.359, "y_scale": 0.359, "print_obj": true, "lock_aspect_ratio": true, "locked": true}`)
		// 롤넘버
		pos, err = excelize.CoordinatesToCellName(3, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, i.Rnum)
		f.SetCellStyle(sheet, pos, pos, style)
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
		// 태그
		pos, err = excelize.CoordinatesToCellName(7, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, strings.Join(i.Tag, ","))

		f.SetCellStyle(sheet, pos, pos, style)
		// 수정사항
		pos, err = excelize.CoordinatesToCellName(8, n+2)
		if err != nil {
			log.Println(err)
		}
		comments := []string{}
		for _, c := range ReverseCommentSlice(i.Comments) {
			comments = append(comments, c.Text)
		}
		f.SetCellValue(sheet, pos, strings.Join(comments, "\n"))
		f.SetCellStyle(sheet, pos, pos, textStyle)
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
		// Deadline2D
		pos, err = excelize.CoordinatesToCellName(11, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, ToNormalTime(i.Ddline2d))
		f.SetCellStyle(sheet, pos, pos, style)
		// Comp
		pos, err = excelize.CoordinatesToCellName(12, n+2)
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
				}`, itemStatus2color(i.Comp.Status)))
		if err != nil {
			log.Println(err)
		}
		text := Status2capString(i.Comp.Status)
		text += "\n" + i.Comp.User
		text += "\n" + ToNormalTime(i.Comp.Predate)
		text += "\n" + ToNormalTime(i.Comp.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Matte
		pos, err = excelize.CoordinatesToCellName(13, n+2)
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
				}`, itemStatus2color(i.Matte.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Matte.Status)
		text += "\n" + i.Matte.User
		text += "\n" + ToNormalTime(i.Matte.Predate)
		text += "\n" + ToNormalTime(i.Matte.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Mg
		pos, err = excelize.CoordinatesToCellName(14, n+2)
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
				}`, itemStatus2color(i.Mg.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Mg.Status)
		text += "\n" + i.Mg.User
		text += "\n" + ToNormalTime(i.Mg.Predate)
		text += "\n" + ToNormalTime(i.Mg.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Deadline3D
		pos, err = excelize.CoordinatesToCellName(15, n+2)
		if err != nil {
			log.Println(err)
		}
		f.SetCellValue(sheet, pos, ToNormalTime(i.Ddline3d))
		f.SetCellStyle(sheet, pos, pos, style)
		// Model
		pos, err = excelize.CoordinatesToCellName(16, n+2)
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
				}`, itemStatus2color(i.Model.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Model.Status)
		text += "\n" + i.Model.User
		text += "\n" + ToNormalTime(i.Model.Predate)
		text += "\n" + ToNormalTime(i.Model.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Mm
		pos, err = excelize.CoordinatesToCellName(17, n+2)
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
				}`, itemStatus2color(i.Mm.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Mm.Status)
		text += "\n" + i.Mm.User
		text += "\n" + ToNormalTime(i.Mm.Predate)
		text += "\n" + ToNormalTime(i.Mm.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Layout
		pos, err = excelize.CoordinatesToCellName(18, n+2)
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
				}`, itemStatus2color(i.Layout.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Layout.Status)
		text += "\n" + i.Layout.User
		text += "\n" + ToNormalTime(i.Layout.Predate)
		text += "\n" + ToNormalTime(i.Layout.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Ani
		pos, err = excelize.CoordinatesToCellName(19, n+2)
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
				}`, itemStatus2color(i.Ani.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Ani.Status)
		text += "\n" + i.Ani.User
		text += "\n" + ToNormalTime(i.Ani.Predate)
		text += "\n" + ToNormalTime(i.Ani.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Fur
		pos, err = excelize.CoordinatesToCellName(20, n+2)
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
				}`, itemStatus2color(i.Fur.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Fur.Status)
		text += "\n" + i.Fur.User
		text += "\n" + ToNormalTime(i.Fur.Predate)
		text += "\n" + ToNormalTime(i.Fur.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Sim
		pos, err = excelize.CoordinatesToCellName(21, n+2)
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
				}`, itemStatus2color(i.Sim.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Sim.Status)
		text += "\n" + i.Sim.User
		text += "\n" + ToNormalTime(i.Sim.Predate)
		text += "\n" + ToNormalTime(i.Sim.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Crowd
		pos, err = excelize.CoordinatesToCellName(22, n+2)
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
				}`, itemStatus2color(i.Crowd.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Crowd.Status)
		text += "\n" + i.Crowd.User
		text += "\n" + ToNormalTime(i.Crowd.Predate)
		text += "\n" + ToNormalTime(i.Crowd.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Fx
		pos, err = excelize.CoordinatesToCellName(23, n+2)
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
				}`, itemStatus2color(i.Fx.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Fx.Status)
		text += "\n" + i.Fx.User
		text += "\n" + ToNormalTime(i.Fx.Predate)
		text += "\n" + ToNormalTime(i.Fx.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Light
		pos, err = excelize.CoordinatesToCellName(24, n+2)
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
				}`, itemStatus2color(i.Light.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Light.Status)
		text += "\n" + i.Light.User
		text += "\n" + ToNormalTime(i.Light.Predate)
		text += "\n" + ToNormalTime(i.Light.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
		// Priviz
		pos, err = excelize.CoordinatesToCellName(25, n+2)
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
				}`, itemStatus2color(i.Previz.Status)))
		if err != nil {
			log.Println(err)
		}
		text = Status2capString(i.Previz.Status)
		text += "\n" + i.Previz.User
		text += "\n" + ToNormalTime(i.Previz.Predate)
		text += "\n" + ToNormalTime(i.Previz.Date)
		f.SetCellValue(sheet, pos, text)
		f.SetCellStyle(sheet, pos, pos, statusStyle)
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
