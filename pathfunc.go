package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"image"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

// userTemppath 함수는 id를 받아서 각 유저별 Temp 경로를 생성 반환한다.
func userTemppath(id string) (string, error) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	path := filepath.Dir(tempDir) + "/" + id
	err = os.MkdirAll(path, 0766)
	if err != nil {
		return path, err
	}
	return path, nil
}

// RemoveXLSX 함수는 폴더 내부 .xlsx 파일을 지운다.
func RemoveXLSX(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		if filepath.Ext(name) != ".xlsx" {
			continue
		}
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveJSON 함수는 폴더 내부 .xlsx 파일을 지운다.
func RemoveJSON(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		if filepath.Ext(name) != ".json" {
			continue
		}
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// GetXLSX 함수는 폴더 내부에 .xlsx 파일을 찾는다.
func GetXLSX(dir string) ([]string, error) {
	var result []string
	d, err := os.Open(dir)
	if err != nil {
		return result, err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return result, err
	}
	for _, name := range names {
		if filepath.Ext(name) != ".xlsx" {
			continue
		}
		result = append(result, filepath.Join(dir, name))
	}
	return result, nil
}

// GetJSON 함수는 폴더 내부에 .json 파일을 찾는다.
func GetJSON(dir string) ([]string, error) {
	var result []string
	d, err := os.Open(dir)
	if err != nil {
		return result, err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return result, err
	}
	for _, name := range names {
		if filepath.Ext(name) != ".json" {
			continue
		}
		result = append(result, filepath.Join(dir, name))
	}
	return result, nil
}

// searchSeq 함수는 탐색할 경로를 입력받고 dpx, exr, png ... 정보를 수집 반환한다.
func searchSeq(searchpath string) ([]ScanPlate, error) {
	// 경로가 존재하는지 체크한다.
	_, err := os.Stat(searchpath)
	if err != nil {
		return nil, err
	}
	paths := make(map[string]ScanPlate)
	err = filepath.Walk(searchpath, func(path string, info os.FileInfo, err error) error {
		// 숨김폴더는 스킵한다.
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return nil //filepath.SkipDir
		}
		// 숨김파일도 스킵한다.
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".dpx", ".exr", ".tif", ".tiff", ".tga", ".png", ".jpg", ".jpeg":
			key, num, err := Seqnum2Sharp(path)
			if err != nil {
				if *flagDebug {
					fmt.Fprintf(os.Stderr, "%s\n", err)
				}
				return nil
			}
			if _, has := paths[key]; has {
				// 이미 수집된 경로가 존재할 때 처리되는 코드
				item := paths[key]
				item.Length++
				item.FrameOut = num
				item.RenderOut = num
				paths[key] = item
			} else {
				var width int
				var height int
				var timecode string
				// width, height 구하기
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
					width, height, err = imageSize(path)
					if err != nil {
						log.Printf("error parsing image %q: %v\n", path, err)
					}
				} else if ext == ".exr" || ext == ".dpx" || ext == ".tif" {
					width, height, err = imageSizeFromIinfo(path)
					if err != nil {
						log.Printf("error parsing image %q: %v\n", path, err)
					}
				}

				// timecode 구하기
				if ext == ".exr" || ext == ".dpx" {
					timecode, err = timecodeFromIinfo(path)
					if err != nil {
						log.Printf("error parsing timecode %q: %v\n", path, err)
					}
				}

				// 한번도 처리된적 없는 이미지가 존재하면 처리되는 코드
				item := ScanPlate{
					Searchpath:  searchpath,
					Dir:         filepath.Dir(path),
					Base:        filepath.Base(key),
					Ext:         ext,
					Length:      1,
					FrameIn:     num,
					FrameOut:    num,
					RenderIn:    num,
					ConvertExt:  ext,
					Width:       width,
					Height:      height,
					TimecodeIn:  timecode, // 시작지점에서 한번 연산한다.
					TimecodeOut: timecode, // 1프레임밖에 없을 수 있다. 디폴트로 저장한다.
				}
				paths[key] = item
			}
		default:
			return nil
		}
		return nil
	})
	if err != nil {
		log.Printf("error walking the path %q: %v\n", searchpath, err)
	}
	var items []ScanPlate
	for _, value := range paths {
		// timecode out 을 이곳에서 처리한다.
		if (value.Ext == ".exr" || value.Ext == ".dpx") && value.Length > 1 {
			endframePath := value.Dir + "/" + fmt.Sprintf(value.Base, value.FrameOut)
			timecodeout, err := timecodeFromIinfo(endframePath)
			if err != nil {
				log.Printf("error parsing timecode %q: %v\n", endframePath, err)
			}
			value.TimecodeOut = timecodeout
		}
		items = append(items, value)
	}
	if len(items) == 0 {
		return nil, errors.New("소스가 존재하지 않습니다")
	}
	return items, nil
}

// searchImage 함수는 탐색할 경로를 입력받고 jpg, png ... 정보를 수집 반환한다.
func searchImage(searchpath string) ([]ScanPlate, error) {
	// 경로가 존재하는지 체크한다.
	_, err := os.Stat(searchpath)
	if err != nil {
		return nil, err
	}
	paths := make(map[string]ScanPlate)
	err = filepath.Walk(searchpath, func(path string, info os.FileInfo, err error) error {
		// 숨김폴더는 스킵한다.
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return nil //filepath.SkipDir
		}
		// 숨김파일도 스킵한다.
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".png", ".jpg", ".jpeg":
			var width int
			var height int
			// width, height 구하기
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				width, height, err = imageSize(path)
				if err != nil {
					log.Printf("error parsing image %q: %v\n", path, err)
				}
			}
			// 한번도 처리된적 없는 이미지가 존재하면 처리되는 코드
			item := ScanPlate{
				Searchpath: searchpath,
				Dir:        filepath.Dir(path),
				Base:       filepath.Base(path),
				Ext:        ext,
				Length:     1,
				FrameIn:    1,
				FrameOut:   1,
				RenderIn:   1,
				ConvertExt: ext,
				Width:      width,
				Height:     height,
			}
			paths[path] = item
		default:
			return nil
		}
		return nil
	})
	if err != nil {
		log.Printf("error walking the path %q: %v\n", searchpath, err)
	}
	var items []ScanPlate
	for _, value := range paths {
		items = append(items, value)
	}
	if len(items) == 0 {
		return nil, errors.New("소스가 존재하지 않습니다")
	}
	return items, nil
}

// Seqnum2Sharp 함수는 경로와 파일명을 받아서 시퀀스부분을 #문자열로 바꾸고 시퀀스의 숫자를 int로 바꾼다.
// "test.0002.jpg" -> "test.####.jpg", 2, nil
func Seqnum2Sharp(filename string) (string, int, error) {
	re, err := regexp.Compile(`([0-9]+)(\.[a-zA-Z]+$)`)
	// 이 정보를 통해서 파일명을 구하는 방식으로 바꾼다.
	if err != nil {
		return filename, -1, errors.New("정규 표현식이 잘못되었습니다")
	}
	results := re.FindStringSubmatch(filename)
	if results == nil {
		return filename, -1, errors.New("경로가 시퀀스 형식이 아닙니다")
	}
	seq := results[1]
	ext := results[2]
	header := filename[:strings.LastIndex(filename, seq+ext)]
	seqNum, err := strconv.Atoi(seq)
	if err != nil {
		return filename, -1, err
	}
	return header + "%0" + strconv.Itoa(len(seq)) + "d" + ext, seqNum, nil
}

func imageSize(path string) (int, int, error) {
	reader, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer reader.Close()
	im, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}
	return im.Width, im.Height, nil
}

func imageSizeFromIinfo(path string) (int, int, error) {
	var width int
	var height int
	_, err := os.Stat(CachedAdminSetting.Iinfo)
	if err != nil {
		return 0, 0, err
	}
	cmd := exec.Command(CachedAdminSetting.Iinfo, path)
	stdout, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	re, err := regexp.Compile(`(\d+)\s+x\s+(\d+)`)
	if err != nil {
		return 0, 0, errors.New("the regular expression is invalid")
	}
	results := re.FindStringSubmatch(string(stdout))
	if results == nil {
		return 0, 0, errors.New("there were no results matching the regular expression condition")
	}
	width, err = strconv.Atoi(results[1]) // results[0]은 "1920 x 1080" 의 묶음이다.
	if err != nil {
		return 0, 0, err
	}
	height, err = strconv.Atoi(results[2])
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func timecodeFromExrheader(path string) (string, error) {
	_, err := os.Stat(CachedAdminSetting.Exrheader)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(CachedAdminSetting.Exrheader, path)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	re, err := regexp.Compile(`time\s+(\d{2}:\d{2}:\d{2}:\d{2})`)
	if err != nil {
		return "", errors.New("the regular expression is invalid")
	}
	results := re.FindStringSubmatch(string(stdout))
	if results == nil {
		return "", errors.New("there were no results matching the regular expression condition")
	}
	timecode := results[1]
	return timecode, nil
}

func timecodeFromIinfo(path string) (string, error) {
	_, err := os.Stat(CachedAdminSetting.Iinfo)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(CachedAdminSetting.Iinfo, "-v", path)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	re, err := regexp.Compile(`smpte:TimeCode:\s+(\d{2}:\d{2}:\d{2}:\d{2})`)
	if err != nil {
		return "", errors.New("the regular expression is invalid")
	}
	results := re.FindStringSubmatch(string(stdout))
	if results == nil {
		return "", errors.New("there were no results matching the regular expression condition")
	}
	timecode := results[1]
	return timecode, nil
}

func PlatePath(item Item) (string, error) {
	var path bytes.Buffer
	pathTmpl, err := template.New("temp").Parse(CachedAdminSetting.PlatePath)
	if err != nil {
		return "", err
	}
	err = pathTmpl.Execute(&path, item)
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

func ShotRootPath(item Item) (string, error) {
	// /Users/woong/show/{{.Project}}/seq
	var path bytes.Buffer
	pathTmpl, err := template.New("temp").Parse(CachedAdminSetting.ShotRootPath)
	if err != nil {
		return "", err
	}
	err = pathTmpl.Execute(&path, item)
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

func ShotPath(item Item) (string, error) {
	// /Users/woong/show/{{.Project}}/seq
	var path bytes.Buffer
	pathTmpl, err := template.New("temp").Parse(CachedAdminSetting.ShotPath)
	if err != nil {
		return "", err
	}
	err = pathTmpl.Execute(&path, item)
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

func SeqPath(item Item) (string, error) {
	// /Users/woong/show/{{.Project}}/seq
	var path bytes.Buffer
	pathTmpl, err := template.New("temp").Parse(CachedAdminSetting.SeqPath)
	if err != nil {
		return "", err
	}
	err = pathTmpl.Execute(&path, item)
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

func AssetRootPath(item Item) (string, error) {
	var path bytes.Buffer
	pathTmpl, err := template.New("temp").Parse(CachedAdminSetting.AssetRootPath)
	if err != nil {
		return "", err
	}
	err = pathTmpl.Execute(&path, item)
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

func AssetTypePath(item Item) (string, error) {
	var path bytes.Buffer
	pathTmpl, err := template.New("temp").Parse(CachedAdminSetting.AssetTypePath)
	if err != nil {
		return "", err
	}
	err = pathTmpl.Execute(&path, item)
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

func AssetPath(item Item) (string, error) {
	var path bytes.Buffer
	pathTmpl, err := template.New("temp").Parse(CachedAdminSetting.AssetPath)
	if err != nil {
		return "", err
	}
	err = pathTmpl.Execute(&path, item)
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

func GenPlatePath(path string) error {
	// 존재하면 폴더를 만들필요가 없다. 바로 리턴한다.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	if CachedAdminSetting.Umask == "" {
		unix.Umask(0)
	} else {
		umask, err := strconv.Atoi(CachedAdminSetting.Umask)
		if err != nil {
			return err
		}
		unix.Umask(umask)
	}
	// 퍼미션을 가지고 온다.
	per, err := strconv.ParseInt(CachedAdminSetting.PlatePathPermission, 8, 64)
	if err != nil {
		return err
	}
	// 폴더를 생성한다.
	err = os.MkdirAll(path, os.FileMode(per))
	if err != nil {
		return err
	}
	// uid, gid 를 설정한다.
	uid, err := strconv.Atoi(CachedAdminSetting.PlatePathUID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(CachedAdminSetting.PlatePathGID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

func GenShotRootPath(path string) error {
	// 존재하면 폴더를 만들필요가 없다. 바로 리턴한다.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	if CachedAdminSetting.Umask == "" {
		unix.Umask(0)
	} else {
		umask, err := strconv.Atoi(CachedAdminSetting.Umask)
		if err != nil {
			return err
		}
		unix.Umask(umask)
	}
	// 퍼미션을 가지고 온다.
	per, err := strconv.ParseInt(CachedAdminSetting.ShotRootPathPermission, 8, 64)
	if err != nil {
		return err
	}
	// 폴더를 생성한다.
	err = os.MkdirAll(path, os.FileMode(per))
	if err != nil {
		return err
	}
	// uid, gid 를 설정한다.
	uid, err := strconv.Atoi(CachedAdminSetting.ShotRootPathUID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(CachedAdminSetting.ShotRootPathGID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

func GenShotPath(path string) error {
	// 존재하면 폴더를 만들필요가 없다. 바로 리턴한다.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	if CachedAdminSetting.Umask == "" {
		unix.Umask(0)
	} else {
		umask, err := strconv.Atoi(CachedAdminSetting.Umask)
		if err != nil {
			return err
		}
		unix.Umask(umask)
	}
	// 퍼미션을 가지고 온다.
	per, err := strconv.ParseInt(CachedAdminSetting.ShotPathPermission, 8, 64)
	if err != nil {
		return err
	}
	// 폴더를 생성한다.
	err = os.MkdirAll(path, os.FileMode(per))
	if err != nil {
		return err
	}
	// uid, gid 를 설정한다.
	uid, err := strconv.Atoi(CachedAdminSetting.ShotPathUID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(CachedAdminSetting.ShotPathGID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

func GenSeqPath(path string) error {
	// 존재하면 폴더를 만들필요가 없다. 바로 리턴한다.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	if CachedAdminSetting.Umask == "" {
		unix.Umask(0)
	} else {
		umask, err := strconv.Atoi(CachedAdminSetting.Umask)
		if err != nil {
			return err
		}
		unix.Umask(umask)
	}
	// 퍼미션을 가지고 온다.
	per, err := strconv.ParseInt(CachedAdminSetting.SeqPathPermission, 8, 64)
	if err != nil {
		return err
	}
	// 폴더를 생성한다.
	err = os.MkdirAll(path, os.FileMode(per))
	if err != nil {
		return err
	}
	// uid, gid 를 설정한다.
	uid, err := strconv.Atoi(CachedAdminSetting.SeqPathUID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(CachedAdminSetting.SeqPathGID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

func GenAssetRootPath(path string) error {
	// 존재하면 폴더를 만들필요가 없다. 바로 리턴한다.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	if CachedAdminSetting.Umask == "" {
		unix.Umask(0)
	} else {
		umask, err := strconv.Atoi(CachedAdminSetting.Umask)
		if err != nil {
			return err
		}
		unix.Umask(umask)
	}
	// 퍼미션을 가지고 온다.
	per, err := strconv.ParseInt(CachedAdminSetting.AssetRootPathPermission, 8, 64)
	if err != nil {
		return err
	}
	// 폴더를 생성한다.
	err = os.MkdirAll(path, os.FileMode(per))
	if err != nil {
		return err
	}
	// uid, gid 를 설정한다.
	uid, err := strconv.Atoi(CachedAdminSetting.AssetRootPathUID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(CachedAdminSetting.AssetRootPathGID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

func GenAssetTypePath(path string) error {
	// 존재하면 폴더를 만들필요가 없다. 바로 리턴한다.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	if CachedAdminSetting.Umask == "" {
		unix.Umask(0)
	} else {
		umask, err := strconv.Atoi(CachedAdminSetting.Umask)
		if err != nil {
			return err
		}
		unix.Umask(umask)
	}
	// 퍼미션을 가지고 온다.
	per, err := strconv.ParseInt(CachedAdminSetting.AssetTypePathPermission, 8, 64)
	if err != nil {
		return err
	}
	// 폴더를 생성한다.
	err = os.MkdirAll(path, os.FileMode(per))
	if err != nil {
		return err
	}
	// uid, gid 를 설정한다.
	uid, err := strconv.Atoi(CachedAdminSetting.AssetTypePathUID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(CachedAdminSetting.AssetTypePathGID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

func GenAssetPath(path string) error {
	// 존재하면 폴더를 만들필요가 없다. 바로 리턴한다.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	if CachedAdminSetting.Umask == "" {
		unix.Umask(0)
	} else {
		umask, err := strconv.Atoi(CachedAdminSetting.Umask)
		if err != nil {
			return err
		}
		unix.Umask(umask)
	}
	// 퍼미션을 가지고 온다.
	per, err := strconv.ParseInt(CachedAdminSetting.AssetPathPermission, 8, 64)
	if err != nil {
		return err
	}
	// 폴더를 생성한다.
	err = os.MkdirAll(path, os.FileMode(per))
	if err != nil {
		return err
	}
	// uid, gid 를 설정한다.
	uid, err := strconv.Atoi(CachedAdminSetting.AssetPathUID)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(CachedAdminSetting.AssetPathGID)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	return nil
}
