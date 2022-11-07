package main

import (
	"strings"

	"github.com/digital-idea/ditime"
	"gopkg.in/mgo.v2"
)

// Excelrow 자료구조는 .xlsx 형식의 자료구조이다.
// 샷네임;작업종류(2D,3D);작업내용;수정사항;링크자료(제목:경로);3D마감;2D마감;FIN날짜;FIN버젼;테그;롤넘버;핸들IN;핸들OUT;JUST타임코드IN;JUST타임코드OUT
type Excelrow struct {
	Name                 string
	NameError            string
	Shottype             string
	ShottypeError        string
	Note                 string
	Comment              string
	Link                 string
	LinkError            string
	Ddline3D             string
	Ddline3DError        string
	Ddline2D             string
	Ddline2DError        string
	Findate              string
	FindateError         string
	Finver               string
	FinverError          string
	Tags                 string
	Rnum                 string
	RnumError            string
	HandleIn             string
	HandleInError        string
	HandleOut            string
	HandleOutError       string
	JustTimecodeIn       string
	JustTimecodeInError  string
	JustTimecodeOut      string
	JustTimecodeOutError string
	Errornum             int
}

func (r *Excelrow) checkerror(session *mgo.Session, project string) {
	_, err := Type(session, project, r.Name)
	if err != nil {
		r.NameError = "등록된 Shot, Asset 이름이 아닙니다"
		r.Errornum++
	}
	if !(regexpShotname.MatchString(r.Name) || regexpAssetname.MatchString(r.Name)) { // 필수값
		r.NameError = "Shot, Asset 이름 형태가 아닙니다"
		r.Errornum++
	}
	if !(r.Shottype == "" || r.Shottype == "2d" || r.Shottype == "3d" || r.Shottype == "2D" || r.Shottype == "3D" || r.Shottype == "2.5d" || r.Shottype == "2.5D") {
		r.ShottypeError = "Not support shotype: 2d, 3d, 2.5d"
		r.Errornum++
	}
	if r.Link != "" {
		for _, l := range strings.Split(r.Link, "\n") {
			var key string
			var value string
			input := strings.Split(l, ":")
			if len(input) != 2 {
				r.LinkError = "key:value 형태로 작성되어있지 않습니다"
				r.Errornum++
			}
			key = input[0]
			value = input[1]
			if key == "" {
				r.LinkError = "key값이 빈 문자열 입니다"
				r.Errornum++
			}
			if value == "" {
				r.LinkError = "value값이 빈 문자열 입니다"
				r.Errornum++
			}
		}
	}
	if r.Rnum != "" && !regexpRnum.MatchString(r.Rnum) {
		r.RnumError = "롤넘버 형식이 A0001 형태가 아닙니다"
		r.Errornum++
	}
	if r.HandleIn != "" && !regexpHandle.MatchString(r.HandleIn) {
		r.HandleInError = "핸들 형식이 아닙니다"
		r.Errornum++
	}
	if r.HandleOut != "" && !regexpHandle.MatchString(r.HandleOut) {
		r.HandleOutError = "핸들 형식이 아닙니다"
		r.Errornum++
	}
	if r.JustTimecodeIn != "" && !regexpTimecode.MatchString(r.JustTimecodeIn) {
		r.JustTimecodeInError = "Timecode 형식이 아닙니다"
		r.Errornum++
	}
	if r.JustTimecodeOut != "" && !regexpTimecode.MatchString(r.JustTimecodeOut) {
		r.JustTimecodeOutError = "Timecode 형식이 아닙니다"
		r.Errornum++
	}
	if r.Finver != "" && !regexpVersion.MatchString(r.Finver) {
		r.FinverError = "값이 3자리 이하 숫자로 이루어져있지 않습니다"
		r.Errornum++
	}
	if r.Ddline3D != "" {
		ddline3d, err := ditime.ToFullTime(19, r.Ddline3D)
		if err != nil {
			r.Ddline3DError = err.Error()
			r.Errornum++
		}
		r.Ddline3D = ddline3d
	}
	if r.Ddline2D != "" {
		ddline2d, err := ditime.ToFullTime(19, r.Ddline2D)
		if err != nil {
			r.Ddline2DError = err.Error()
			r.Errornum++
		}
		r.Ddline2D = ddline2d
	}
	if r.Findate != "" {
		findate, err := ditime.ToFullTime(19, r.Findate)
		if err != nil {
			r.FindateError = err.Error()
			r.Errornum++
		}
		r.Findate = findate
	}
}
