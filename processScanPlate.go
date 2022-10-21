package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"image"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ProcessScanPlateRender() {
	// 버퍼 채널을 만든다.
	jobs := make(chan ScanPlate, *flagProcessBufferSize) // 작업을 대기할 버퍼를 만든다.
	// worker 프로세스를 지정한 개수만큼 실행시킨다.
	for w := 1; w <= *flagMaxProcessNum; w++ {
		go workerScanPlate(jobs)
	}
	// queueingItem을 실행시킨다.
	go queueingScanPlateItem(jobs)

	// ProcessMain()이 종료되지 않는 역할을 한다.
	select {}
}

// workerScanPlate 합수는 ScanPlate 데이터를 jobs로 보낸다.
func workerScanPlate(jobs <-chan ScanPlate) {
	for job := range jobs {
		// job은 리뷰타입이다.
		switch job.Ext {
		case ".jpg", ".png":
			processingScanPlateImageItem(job)
		case ".exr", ".dpx":
			processingScanPlateImageItem(job)
		case ".mov":
			processingScanPlateImageItem(job)
		default:
			processingScanPlateImageItem(job)
		}
	}
}

// queueingScanPlateItem 은 연산할 ScanPlate 아이템을 jobs 채널에 전송한다.
func queueingScanPlateItem(jobs chan<- ScanPlate) {
	for {
		if *flagDebug {
			fmt.Printf("wait %d sec before scanplate process\n", *flagProcessDuration)
		}
		time.Sleep(time.Second * time.Duration(*flagProcessDuration))
		// ProcessStatus가 wait인 item을 가져온다.
		scanPlate, err := GetWaitProcessStatusScanPlate() // 이 함수로 반환되는 아이템은 리뷰 아이템은 상태가 queued가 된 리뷰 아이템이다.
		if err != nil {
			// 가지고 올 문서가 없다면 기다렸다가 continue.
			if err != mongo.ErrNoDocuments {
				if *flagDebug {
					log.Println(err)
				}
			}
			continue
		}
		if *flagDebug {
			fmt.Println(scanPlate)
		}
		jobs <- scanPlate
	}
}

func processingScanPlateImageItem(scan ScanPlate) {
	scanID := scan.ID.Hex()
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 연산중으로 설정한다.
	err = SetScanPlateProcessStatus(client, scanID, "processing")
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	initStatusID, err := GetInitStatusIDV2(client)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
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
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	err = thumbnailImagePathTmpl.Execute(&thumbnailImagePath, item)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	item.Thumpath = thumbnailImagePath.String()

	// 썸네일 MOV 경로를 설정합니다.
	var thumbnailMovPath bytes.Buffer
	thumbnailMovPathTmpl, err := template.New("thumbnailMovPath").Parse(CachedAdminSetting.ThumbnailMovPath)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	err = thumbnailMovPathTmpl.Execute(&thumbnailMovPath, item)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	item.Thummov = thumbnailMovPath.String()

	// 플레이트 경로를 설정합니다.
	item.Platepath, err = PlatePath(item)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	// Task 셋팅
	tasks, err := AllTaskSettingsV2(client)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
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
			Status:       ASSIGN, // legacy, 샷의 경우 합성팀을 무조건 거쳐야 한다. Assign상태로 만든다.
			StatusV2:     initStatusID,
			Pipelinestep: task.Pipelinestep, // 파이프라인 스텝을 설정한다.
		}
		item.Tasks[task.Name] = t
	}

	err = addItemV2(client, item)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 썸네일처리
	uid, err := strconv.Atoi(CachedAdminSetting.ThumbnailImagePathUID)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	gid, err := strconv.Atoi(CachedAdminSetting.ThumbnailImagePathGID)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	permission, err := strconv.ParseInt(CachedAdminSetting.ThumbnailImagePathPermission, 8, 64)
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 썸네일 경로를 생성한다.
	path, _ := path.Split(thumbnailImagePath.String())
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 폴더를 생성한다.
		err = os.MkdirAll(path, os.FileMode(permission))
		if err != nil {
			err = SetScanPlateErrStatus(client, scanID, err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		// 위 폴더가 잘 생성되어 존재한다면 폴더의 권한을 설정한다.
		if _, err := os.Stat(path); os.IsExist(err) {
			err = os.Chown(path, uid, gid)
			if err != nil {
				err = SetScanPlateErrStatus(client, scanID, err.Error())
				if err != nil {
					log.Println(err)
				}
				return
			}
		}
	}

	// 썸네일 이미지가 이미 존재하는 경우 이미지 파일을 지운다.
	if _, err := os.Stat(thumbnailImagePath.String()); os.IsExist(err) {
		err = os.Remove(thumbnailImagePath.String())
		if err != nil {
			err = SetScanPlateErrStatus(client, scanID, err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
	}

	if scan.Ext == ".jpg" || scan.Ext == ".png" {
		// .jpg, .png 라면 바로 썸네일을 처리한다.
		thumbnailData, err := os.Open(item.Scanname)
		if err != nil {
			err = SetScanPlateErrStatus(client, scanID, err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		// 사용자가 업로드한 데이터를 이미지 자료구조로 만들고 리사이즈 한다.
		img, _, err := image.Decode(thumbnailData) // 전송된 바이트 파일을 이미지 자료구조로 변환한다.
		if err != nil {
			err = SetScanPlateErrStatus(client, scanID, err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		resizedImage := imaging.Fill(img, CachedAdminSetting.ThumbnailImageWidth, CachedAdminSetting.ThumbnailImageHeight, imaging.Center, imaging.Lanczos)
		err = imaging.Save(resizedImage, thumbnailImagePath.String())
		if err != nil {
			err = SetScanPlateErrStatus(client, scanID, err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
	} else if scan.Ext == ".exr" {
		// 중간 프레임을 구한다.
		middleFrame := scan.FrameIn + (scan.FrameOut-scan.FrameIn)/2
		filename := fmt.Sprintf(scan.Base, middleFrame)
		src := fmt.Sprintf("%s/%s", scan.Dir, filename)
		dst := thumbnailImagePath.String()

		// oiiotool을 이용해서 썸네일 경로에 바로 쓰기 한다.
		args := []string{
			src,
			"--colorconfig",
			CachedAdminSetting.OCIOConfig,
			"--colorconvert",
			scan.InColorspace,
			scan.OutColorspace,
			"--fit",
			fmt.Sprintf("%dx%d", CachedAdminSetting.ThumbnailImageWidth, CachedAdminSetting.ThumbnailImageHeight),
			"-o",
			dst,
		}
		if *flagDebug {
			fmt.Println(CachedAdminSetting.OpenImageIO, strings.Join(args, " "))
		}
		err = exec.Command(CachedAdminSetting.OpenImageIO, args...).Run()
		if err != nil {
			err = SetScanPlateErrStatus(client, scanID, err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
	} else if scan.Ext == ".mov" {
		fmt.Println("mov")
	}

	// 디스크 물리적인 연산을 수행합니다.

	// 플레이트 경로를 생성합니다.
	if scan.GenPlatePath {
		err = GenPlatePath(item.Platepath)
		if err != nil {
			err = SetScanPlateErrStatus(client, scanID, err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
	}

	// 플레이트 복사가 필요하다면 플레이트를 복사합니다.
	if scan.CopyPlate {
		for i := scan.FrameIn; i <= scan.FrameOut; i++ {
			filename := fmt.Sprintf(scan.Base, i)
			src := fmt.Sprintf("%s/%s", scan.Dir, filename)
			dst := fmt.Sprintf("%s/%s", item.Platepath, filename)
			err = CopyPlate(src, dst)
			if err != nil {
				err = SetScanPlateErrStatus(client, scanID, err.Error())
				if err != nil {
					log.Println(err)
				}
				return
			}
		}
	}

	// Proxy Jpg 옵션이 켜 있다면 Proxy jpg를 생성한다.
	if scan.ProxyJpg {
		// Proxy Jpg 가 생성될 경로를 만든다.
		err = GenPlatePath(item.Platepath + "_jpg")
		if err != nil {
			err = SetScanPlateErrStatus(client, scanID, err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
		for i := scan.FrameIn; i <= scan.FrameOut; i++ {
			filename := fmt.Sprintf(scan.Base, i)
			src := fmt.Sprintf("%s/%s", scan.Dir, filename)
			dst := fmt.Sprintf("%s_jpg/%s", item.Platepath, strings.Replace(filename, ".exr", ".jpg", -1))

			args := []string{
				src,
				"--colorconfig",
				CachedAdminSetting.OCIOConfig,
				"--colorconvert",
				scan.InColorspace,
				scan.OutColorspace,
				"-o",
				dst,
			}
			if *flagDebug {
				fmt.Println(CachedAdminSetting.OpenImageIO, strings.Join(args, " "))
			}
			err = exec.Command(CachedAdminSetting.OpenImageIO, args...).Run()
			if err != nil {
				err = SetScanPlateErrStatus(client, scanID, err.Error())
				if err != nil {
					log.Println(err)
				}
				return
			}
		}
	}

	// 연산 상태를 done 으로 바꾼다.
	err = SetScanPlateProcessStatus(client, scanID, "done")
	if err != nil {
		err = SetScanPlateErrStatus(client, scanID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
}
