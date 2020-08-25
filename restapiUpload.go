package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"golang.org/x/sys/unix"
	"gopkg.in/mgo.v2"
)

// handleAPIUploadThumbnail 함수는 thumbnail 이미지를 업로드 하는 RestAPI 이다.
func handleAPIUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		UserID  string `json:"userid"`
	}
	rcp := Recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.UserID, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// 어드민 셋팅을 불러온다.
	adminSetting, err := GetAdminSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	umask, err := strconv.Atoi(adminSetting.Umask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uid, err := strconv.Atoi(adminSetting.ThumbnailImagePathUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	gid, err := strconv.Atoi(adminSetting.ThumbnailImagePathGID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	permission, err := strconv.ParseInt(adminSetting.ThumbnailImagePathPermission, 8, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 폼을 분석한다.
	r.ParseForm()
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Project = project

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name을 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Name = name
	typ := r.FormValue("type")
	if typ == "" {
		http.Error(w, "type을 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Type = typ

	buffer := adminSetting.MultipartFormBufferSize
	err = r.ParseMultipartForm(int64(buffer))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, files := range r.MultipartForm.File {
		for _, f := range files {
			if f.Size == 0 {
				http.Error(w, "파일사이즈가 0 바이트입니다", http.StatusBadRequest)
				return
			}
			file, err := f.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				continue
			}
			defer file.Close()
			unix.Umask(umask)
			switch f.Header.Get("Content-Type") {
			case "image/jpeg", "image/png":
				data, err := ioutil.ReadAll(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				var thumbImgPath bytes.Buffer
				thumbImgPathTmpl, err := template.New("thumbImgPath").Parse(adminSetting.ThumbnailImagePath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = thumbImgPathTmpl.Execute(&thumbImgPath, rcp)
				if _, err := os.Stat(thumbImgPath.String()); os.IsExist(err) {
					// 썸네일 이미지가 이미 존재하는 경우, 지운다.
					err = os.Remove(thumbImgPath.String())
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
				// 썸네일 경로를 생서한다.
				path, _ := path.Split(thumbImgPath.String())
				if _, err := os.Stat(path); os.IsNotExist(err) {
					// 폴더를 생성한다.
					err = os.MkdirAll(path, os.FileMode(permission))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					// 폴더의 권한을 설정한다.
					err = os.Chown(path, uid, gid)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
				err = ioutil.WriteFile(thumbImgPath.String(), data, os.FileMode(permission))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			default:
				http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
				return
			}
		}
	}
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
