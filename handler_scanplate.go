package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func handleScanPlate(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type recipe struct {
		User
		SearchOption
		Setting          Setting
		Pipelinesteps    []Pipelinestep
		Projectlist      []string
		TasksettingNames []string
		Stages           []Stage
		Status           []Status
		Colorspaces      []string `json:"colorspaces"`
	}
	rcp := recipe{}
	rcp.Projectlist, err = ProjectlistV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.TasksettingNames, err = TaskSettingNamesV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Status, err = AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Stages, err = AllStagesV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User, err = getUserV2(client, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Pipelinesteps, err = allPipelinesteps(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Setting = CachedAdminSetting
	rcp.Colorspaces, err = loadOCIOConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "scanplate", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUploadScanPlate(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}

	buffer := CachedAdminSetting.MultipartFormBufferSize
	err = r.ParseMultipartForm(int64(buffer)) // grab the multipart form
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	path := CachedAdminSetting.ScanPlateUploadPath + "/temp"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0775)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	type recipe struct {
		Path     string `json:"path"`
		Ext      string `json:"ext"`
		Filename string `json:"filename"`
	}
	for _, files := range r.MultipartForm.File {
		for _, f := range files {
			if f.Size == 0 {
				http.Error(w, "파일사이즈가 0 바이트입니다", http.StatusInternalServerError)
				return
			}
			file, err := f.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				continue
			}
			defer file.Close()
			mimeType := f.Header.Get("Content-Type")
			switch mimeType {
			case "image/jpeg", "image/png", "video/quicktime", "image/x-exr", "application/octet-stream":
				ext := strings.ToLower(filepath.Ext(f.Filename))
				if ext != ".jpg" && ext != ".png" && ext != ".mov" && ext != ".dpx" && ext != ".exr" { // .jpg .dpx .exr 외에는 허용하지 않는다.
					http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
					return
				}
				data, err := ioutil.ReadAll(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = ioutil.WriteFile(path+"/"+f.Filename, data, 0775)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				rcp := recipe{}
				rcp.Path = path
				rcp.Filename = f.Filename
				rcp.Ext = ext
				json, err := json.Marshal(rcp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Write(json)
			default:
				http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
				return
			}
		}
	}
}

func deleteScanPlateTemp(w http.ResponseWriter, r *http.Request) {
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// token 체크
	_, accessLevel, err := TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accessLevel != AdminAccessLevel {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	path := CachedAdminSetting.ScanPlateUploadPath + "/temp"

	os.RemoveAll(path)
	os.MkdirAll(path, 0775)

	type recipe struct {
		Path string `json:"path"`
	}
	rcp := recipe{}
	rcp.Path = path
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPISearchFootages(w http.ResponseWriter, r *http.Request) {
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// token 체크
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "need path value", http.StatusBadRequest)
		return
	}
	typ := r.FormValue("type")
	if typ == "" {
		typ = "org"
	}
	method := r.FormValue("method")
	incolorspace := r.FormValue("incolorspace")
	outcolorspace := r.FormValue("outcolorspace")
	regex := r.FormValue("regex")

	regexp, regerr := regexp.Compile(regex)
	if regerr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var items []ScanPlate
	if method == "" || method == "image" {
		items, err = searchImage(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		items, err = searchSeq(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 사용자에게 받은 type을 설정한다.
	for n, _ := range items {
		items[n].Type = typ
		items[n].InColorspace = incolorspace   // 사용자에게 받은 InColorspace를 설정한다.
		items[n].OutColorspace = outcolorspace // 사용자에게 받은 OutColorspace를 설정한다.

		regGroup := regexp.SubexpNames()
		regValue := regexp.FindStringSubmatch(items[n].Base)
		// 사용자에게 받은 Regex값을 이용해서 Name값을 추출한다.
		if regerr == nil {
			for order, groupName := range regGroup {
				switch groupName {
				case "name", "Name":
					items[n].Name = regValue[order]
				case "type", "Type":
					items[n].Type = regValue[order]
				}
			}
		}
	}

	data, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIOcioColorspace(w http.ResponseWriter, r *http.Request) {
	type recipe struct {
		Colorspaces []string `json:"colorspaces"`
	}
	rcp := recipe{}
	colorspace, err := loadOCIOConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Colorspaces = colorspace
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIScanPlate(w http.ResponseWriter, r *http.Request) {
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// token 체크
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	initStatusID, err := GetInitStatusIDV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scan := ScanPlate{}
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&scan)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// ScanPlate 자료구조로 Item 자료구조를 만든다.
	item := Item{}
	item.Name = scan.Name
	item.Type = scan.Type
	item.UseType = scan.Type
	item.Project = scan.Project
	item.ID = scan.Name + "_" + scan.Type
	item.Status = ASSIGN
	item.StatusV2 = initStatusID
	item.Scantime = time.Now().Format(time.RFC3339)
	item.Updatetime = time.Now().Format(time.RFC3339)
	item.Scanname = scan.Dir + "/" + scan.Base
	item.Dataname = scan.Dir + "/" + scan.Base
	if scan.Width != 0 && scan.Height != 0 {
		item.Platesize = fmt.Sprintf("%dx%d", scan.Width, scan.Height)
	}
	if scan.UndistortionWidth != 0 && scan.UndistortionHeight != 0 {
		item.Undistortionsize = fmt.Sprintf("%dx%d", scan.UndistortionWidth, scan.UndistortionHeight)
		item.Dsize = fmt.Sprintf("%dx%d", scan.UndistortionWidth, scan.UndistortionHeight)
	}
	if scan.RenderWidth != 0 && scan.RenderHeight != 0 {
		item.Rendersize = fmt.Sprintf("%dx%d", scan.RenderWidth, scan.RenderHeight)
	}

	item.Tasks = make(map[string]Task)
	item.SetSeq()
	item.SetCut()

	// 썸네일 이미지 경로를 설정합니다.
	var thumbnailImagePath bytes.Buffer
	thumbnailImagePathTmpl, err := template.New("thumbnailImagePath").Parse(CachedAdminSetting.ThumbnailImagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = thumbnailImagePathTmpl.Execute(&thumbnailImagePath, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	item.Thumpath = thumbnailImagePath.String()

	// 썸네일 MOV 경로를 설정합니다.
	var thumbnailMovPath bytes.Buffer
	thumbnailMovPathTmpl, err := template.New("thumbnailMovPath").Parse(CachedAdminSetting.ThumbnailMovPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = thumbnailMovPathTmpl.Execute(&thumbnailMovPath, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	item.Thummov = thumbnailMovPath.String()

	// 플레이트 경로를 설정합니다.
	var platePath bytes.Buffer
	platePathTmpl, err := template.New("platePath").Parse(CachedAdminSetting.PlatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = platePathTmpl.Execute(&platePath, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	item.Platepath = platePath.String()

	// Task 셋팅
	tasks, err := AllTaskSettingsV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for _, task := range tasks {
		if !task.InitGenerate {
			continue
		}
		if task.Type != "shot" {
			continue
		}
		t := Task{
			Title:        task.Name,
			Status:       ASSIGN, // 샷의 경우 합성팀을 무조건 거쳐야 한다. Assign상태로 만든다. // legacy
			StatusV2:     initStatusID,
			Pipelinestep: task.Pipelinestep, // 파이프라인 스텝을 설정한다.
		}
		item.Tasks[task.Name] = t
	}

	err = item.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = addItemV2(client, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// 썸네일처리
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

	// 썸네일 이미지가 이미 존재하는 경우 이미지 파일을 지운다.
	if _, err := os.Stat(thumbnailImagePath.String()); os.IsExist(err) {
		err = os.Remove(thumbnailImagePath.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 썸네일 경로를 생성한다.
	path, _ := path.Split(thumbnailImagePath.String())
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

	thumbnailData, err := os.Open(item.Scanname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 사용자가 업로드한 데이터를 이미지 자료구조로 만들고 리사이즈 한다.
	img, _, err := image.Decode(thumbnailData) // 전송된 바이트 파일을 이미지 자료구조로 변환한다.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resizedImage := imaging.Fill(img, CachedAdminSetting.ThumbnailImageWidth, CachedAdminSetting.ThumbnailImageHeight, imaging.Center, imaging.Lanczos)
	err = imaging.Save(resizedImage, thumbnailImagePath.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Json 값 출력
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
