package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/disintegration/imaging"
	"golang.org/x/sys/unix"
	"gopkg.in/mgo.v2"
)

// handleAPIUploadThumbnail 함수는 thumbnail 이미지를 업로드 하는 RestAPI 이다.
func handleAPIUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		Project        string `json:"project"`
		Name           string `json:"name"`
		Type           string `json:"type"`
		Path           string `json:"path"`
		UploadFilename string `json:"uploadfilename"`
	}
	rcp := Recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// 어드민 셋팅을 불러온다.
	umask, err := strconv.Atoi(CachedAdminSetting.Umask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uid, err := strconv.Atoi(CachedAdminSetting.ThumbnailImagePathUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	gid, err := strconv.Atoi(CachedAdminSetting.ThumbnailImagePathGID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	permission, err := strconv.ParseInt(CachedAdminSetting.ThumbnailImagePathPermission, 8, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 폼을 분석한다.
	err = r.ParseMultipartForm(int64(CachedAdminSetting.MultipartFormBufferSize))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "need project", http.StatusBadRequest)
		return
	}
	rcp.Project = project

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "need name", http.StatusBadRequest)
		return
	}
	rcp.Name = name
	typ := r.FormValue("type")
	if typ == "" {
		itemType, err := Type(session, rcp.Project, rcp.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.Type = itemType
	} else {
		rcp.Type = typ
	}

	if len(r.MultipartForm.File) == 0 { // 파일이 없다면 에러처리한다.
		http.Error(w, "need thumbnail image path", http.StatusBadRequest)
		return
	}

	// 썸네일이 존재한다면 썸네일을 처리한다.
	for _, files := range r.MultipartForm.File {
		if len(files) != 1 { // 파일이 복수일 때
			http.Error(w, "multiple files cannot be set", http.StatusBadRequest)
			return
		}
		for _, f := range files {
			rcp.UploadFilename = f.Filename
			if f.Size == 0 {
				http.Error(w, "file size is 0 bytes", http.StatusBadRequest)
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
				// adminsetting에 설정된 썸네일 템플릿에 실제 값을 넣는다.
				var thumbImgPath bytes.Buffer
				thumbImgPathTmpl, err := template.New("thumbImgPath").Parse(CachedAdminSetting.ThumbnailImagePath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = thumbImgPathTmpl.Execute(&thumbImgPath, rcp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// 썸네일 이미지가 이미 존재하는 경우 이미지 파일을 지운다.
				if _, err := os.Stat(thumbImgPath.String()); os.IsExist(err) {
					err = os.Remove(thumbImgPath.String())
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
				// 썸네일 경로를 생성한다.
				path, _ := path.Split(thumbImgPath.String())
				if _, err := os.Stat(path); os.IsNotExist(err) {
					// 폴더를 생성한다.
					err = os.MkdirAll(path, os.FileMode(permission))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					// 위 폴더가 잘 생성되어 존재한다면 폴더의 권한을 설정한다.
					if _, err := os.Stat(path); os.IsExist(err) {
						err = os.Chown(path, uid, gid)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
					}
				}
				// 사용자가 업로드한 데이터를 이미지 자료구조로 만들고 리사이즈 한다.
				img, _, err := image.Decode(bytes.NewReader(data)) // 전송된 바이트 파일을 이미지 자료구조로 변환한다.
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				resizedImage := imaging.Fill(img, CachedAdminSetting.ThumbnailImageWidth, CachedAdminSetting.ThumbnailImageHeight, imaging.Center, imaging.Lanczos)
				err = imaging.Save(resizedImage, thumbImgPath.String())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				rcp.Path = thumbImgPath.String()
			default:
				http.Error(w, "not support format", http.StatusBadRequest)
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
