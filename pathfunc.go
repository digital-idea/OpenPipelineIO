package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"html/template"
	"image"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/agorman/go-timecode/v2"
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

// searchMovs 함수는 탐색할 경로를 입력받고 mov 정보를 수집 반환한다.
func searchMovs(searchpath string) ([]ScanPlate, error) {
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
		case ".mov":
			var width int
			var height int
			var timecodeIn string
			var timecodeOut string
			var length int
			var fps string
			// width, height, timecode 구하기
			if ext == ".mov" {
				width, height, err = movSizeFromFFprobe(path)
				if err != nil {
					log.Printf("error parsing mov size %q: %v\n", path, err)
				}
				length, err = movDurationFromFFprobe(path)
				if err != nil {
					log.Printf("error parsing mov duration %q: %v\n", path, err)
				}
				timecodeIn, err = movTimecodeInFromFFprobe(path)
				if err != nil {
					log.Printf("error parsing mov timecode %q: %v\n", path, err)
				}
				tc, err := timecode.Parse(timecode.R30, timecodeIn)
				if err != nil {
					log.Printf("error parsing mov timecode %q: %v\n", path, err)
				}
				timecodeOut = tc.Add(uint64(length)).String()
				fps, err = movFpsFromFFprobe(path)
				if err != nil {
					log.Printf("error parsing mov %q: %v\n", path, err)
				}
			}

			// 한번도 처리된적 없는 이미지가 존재하면 처리되는 코드
			item := ScanPlate{
				Searchpath:  searchpath,
				Dir:         filepath.Dir(path),
				Base:        filepath.Base(path),
				Ext:         ext,
				Length:      length,
				FrameIn:     1,
				FrameOut:    length,
				RenderIn:    1,
				RenderOut:   length,
				ConvertExt:  ext,
				Width:       width,
				Height:      height,
				TimecodeIn:  timecodeIn,
				TimecodeOut: timecodeOut,
				Fps:         fps,
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
		return nil, errors.New("no sources")
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

func movSize(path string) (int, int, error) {
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

// inputdeviceDpx 함수는 dpx파일경로를 입력받아서 input device name metadata 값을 반환한다.
func inputdeviceDpx(path string) (string, error) {
	if strings.ToLower(filepath.Ext(path)) != ".dpx" {
		return "", errors.New("확장자가 dpx가 아닙니다.")
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buffer := bufio.NewReader(f)
	// dpx파일의 최초 4바이트를 읽어온다.
	c := make([]byte, 4) // 4자리의 바이트설정.
	_, err = buffer.Read(c)
	if err != nil {
		return "", err
	}
	if !(string(c) == "SDPX" || string(c) == "XPDS") {
		return "", errors.New("파일구조가 DPX 형태가 아닙니다.")
	}
	// input device name 을 구하려면 1556 바이트로 이동해야한다.
	_, err = buffer.Discard(1556 - len(c)) // magicNum을 읽었기 때문에 그 길이만큼 뺀다.
	if err != nil {
		return "", err
	}
	c = make([]byte, 32) // input device name은 32바이트 ASCII 로 구성되어있다.
	_, err = buffer.Read(c)
	if err != nil {
		return "", err
	}
	// input device name값은 big/little endian과 상관없다.
	str := []byte{}
	for _, r := range c {
		if r == 0 {
			break
		}
		str = append(str, r)
	}
	return string(str), nil
}

// imagesizeDpx 함수는 dpx파일경로를 입력받아서 x,y 픽셀수를 반환한다.
func imagesizeDpx(path string) (uint32, uint32, error) {
	// dpx파일의 x,y값을 구하기 위해서는..
	// 772바이트(x), 776바이트(y) 에서 사이즈를 읽으면 된다.
	if strings.ToLower(filepath.Ext(path)) != ".dpx" {
		return 0, 0, errors.New("확장자가 dpx가 아닙니다.")
	}
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	buffer := bufio.NewReader(f)
	// dpx파일의 최초 4바이트를 읽어온다.
	c := make([]byte, 4) // 4자리의 바이트설정.
	_, err = buffer.Read(c)
	if err != nil {
		return 0, 0, err
	}
	if !(string(c) == "SDPX" || string(c) == "XPDS") {
		return 0, 0, errors.New("파일구조가 DPX 형태가 아닙니다.")
	}
	magicNum := string(c)
	var w uint32
	var h uint32
	// X값 구하기.
	// DPX x값 영역으로 이동한다.
	// magicNum을 읽었기 때문에 그 길이만큼 뺀다.
	_, err = buffer.Discard(772 - len(c))
	if err != nil {
		return 0, 0, err
	}
	_, err = buffer.Read(c)
	if err != nil {
		return 0, 0, err
	}
	if magicNum == "SDPX" {
		w = binary.BigEndian.Uint32(c)
	} else {
		w = binary.LittleEndian.Uint32(c)
	}
	// Y값 구하기.
	_, err = buffer.Read(c) // 다음 4바이트를 읽는다.
	if err != nil {
		return 0, 0, err
	}
	if magicNum == "SDPX" {
		h = binary.BigEndian.Uint32(c)
	} else {
		h = binary.LittleEndian.Uint32(c)
	}
	return w, h, nil
}

// timecodeDpx는 dpx파일경로를 받아서 타임코드를 반환한다.
// 타임코드검색어가 파일에 들어있다면 실제 타임코드와 nil을 반환한다.
func timecodeDpx(path string) (string, error) {
	if strings.ToLower(filepath.Ext(path)) != ".dpx" {
		return "", errors.New("확장자가 dpx가 아닙니다.")
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buffer := bufio.NewReader(f)
	if err != nil {
		return "", err
	}
	// dpx파일의 최초 4바이트를 읽어온다.
	c := make([]byte, 4) // 4자리의 바이트설정.
	_, err = buffer.Read(c)
	if err != nil {
		return "", err
	}
	if !(string(c) == "SDPX" || string(c) == "XPDS") {
		return "", errors.New("파일구조가 DPX 형태가 아닙니다.")
	}
	magicNum := string(c)
	// Television Area timecode 영역인 1920 바이트로 이동한다.
	_, err = buffer.Discard(1920 - len(c))
	if err != nil {
		return "", err
	}
	_, err = buffer.Read(c)
	if err != nil {
		return "", err
	}
	var timecodeUint32 uint32
	if magicNum == "SDPX" {
		timecodeUint32 = binary.BigEndian.Uint32(c)
	} else {
		timecodeUint32 = binary.LittleEndian.Uint32(c)
	}
	timecode := fmt.Sprintf("%08x", timecodeUint32)
	return addColons(timecode), nil
}

// addColons 함수는 타임코드를 받아서 읽기편하도록 콜론을 넣어준다.
func addColons(timecode string) string {
	if len(timecode) == 8 {
		return fmt.Sprintf("%s:%s:%s:%s", timecode[0:2], timecode[2:4], timecode[4:6], timecode[6:8])
	}
	return timecode
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

func movSizeFromFFprobe(path string) (int, int, error) {
	var width int
	var height int
	_, err := os.Stat(CachedAdminSetting.FFprobe)
	if err != nil {
		return 0, 0, err
	}
	cmd := exec.Command(CachedAdminSetting.FFprobe, path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return 0, 0, err
	}
	re, err := regexp.Compile(`\s+(\d+)x(\d+),`)
	if err != nil {
		return 0, 0, errors.New("the regular expression is invalid")
	}
	results := re.FindStringSubmatch(stderr.String())
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

func movTimecodeInFromFFprobe(path string) (string, error) {
	_, err := os.Stat(CachedAdminSetting.FFprobe)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(CachedAdminSetting.FFprobe, path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	re, err := regexp.Compile(`<StartTimecode>(\d{2}:\d{2}:\d{2}:\d{2})</StartTimecode>`)
	if err != nil {
		return "", errors.New("the regular expression is invalid")
	}
	results := re.FindStringSubmatch(stderr.String())
	if results == nil {
		return "", errors.New("there were no results matching the regular expression condition")
	}
	timecode := results[1]
	return timecode, nil
}

func movDurationFromFFprobe(path string) (int, error) {
	_, err := os.Stat(CachedAdminSetting.FFprobe)
	if err != nil {
		return 0, err
	}
	cmd := exec.Command(CachedAdminSetting.FFprobe, path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return 0, err
	}
	re, err := regexp.Compile(`<Duration>(\d+)</Duration>`)
	if err != nil {
		return 0, errors.New("the regular expression is invalid")
	}
	results := re.FindStringSubmatch(stderr.String())
	if results == nil {
		return 0, errors.New("there were no results matching the regular expression condition")
	}
	duration, err := strconv.Atoi(results[1])
	if err != nil {
		return 0, err
	}
	return duration, nil
}

func movFpsFromFFprobe(path string) (string, error) {
	_, err := os.Stat(CachedAdminSetting.FFprobe)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(CachedAdminSetting.FFprobe, path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	re, err := regexp.Compile(`<FrameRate>([0-9.]+)([a-z]*)</FrameRate>`)
	if err != nil {
		return "", errors.New("the regular expression is invalid")
	}
	results := re.FindStringSubmatch(stderr.String())
	if results == nil {
		return "", errors.New("there were no results matching the regular expression condition")
	}
	fps := results[1]
	return fps, nil
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

func ThumbnailMovPath(item Item) (string, error) {
	var path bytes.Buffer
	pathTmpl, err := template.New("temp").Parse(CachedAdminSetting.ThumbnailMovPath)
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

func CopyPlate(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func GenSlateString(scan ScanPlate) string {
	font := CachedAdminSetting.SlateFontPath
	fontSize := "30"
	fontColor := "white"
	fade := "0.3"
	border := 50
	borderColor := "black"
	lut := fmt.Sprintf("lut3d=%s", scan.LutPath)
	scale := fmt.Sprintf("scale=%d:%d,setsar=1:1", scan.Width/scan.ProxyRatio, scan.Height/scan.ProxyRatio)
	topArea := fmt.Sprintf("drawbox=x=0:y=0:w=iw:h=%d:color=%s@%s:width=iw:height=%d:t=%d", border, borderColor, fade, border, border)
	bottomArea := fmt.Sprintf("drawbox=x=0:y=ih-50:w=iw:h=ih:color=%s@%s:width=iw:height=50:t=50", borderColor, fade)
	leftTop := fmt.Sprintf("drawtext=fontfile=%s: text='%s': fontsize=%s: x=10: y=((50-th)/2) : fontcolor=%s", font, scan.Project, fontSize, fontColor)
	centerTop := fmt.Sprintf("drawtext=fontfile=%s: text='%s': fontsize=%s: x=(w-tw)/2: y=((50-th)/2) : fontcolor=%s", font, scan.Name, fontSize, fontColor)
	rightTop := fmt.Sprintf("drawtext=fontfile=%s: text='%%{localtime\\:%%Y-%%m-%%d %%a %%T}': fontsize=%s: x=(w-tw-10): y=((50-th)/2) : fontcolor=%s", font, fontSize, fontColor)
	leftBottom := fmt.Sprintf("drawtext=fontfile=%s: text='Plate %s': fontsize=%s: x=10: y=h-50+((50-th)/2): fontcolor=%s", font, scan.Type, fontSize, fontColor)
	centerBottom := fmt.Sprintf("drawtext=fontfile=%s: timecode='%s': r=%s: fontsize=%s: x=(w-tw)/2: y=h-50+((50-th)/2): fontcolor=%s", font, strings.Replace(scan.TimecodeIn, ":", "\\:", -1), scan.Fps, fontSize, fontColor)
	rightBottom := fmt.Sprintf("drawtext=fontfile=%s: text='%%{n}': start_number=%d: r=%s: fontsize=%s: x=(w-tw-10): y=h-50+((50-th)/2): fontcolor=%s", font, scan.FrameIn, fontSize, scan.Fps, fontColor)
	divisibleByError := "pad=ceil(iw/2)*2:ceil(ih/2)*2"
	template := []string{
		lut, scale,
		topArea, bottomArea,
		leftTop, centerTop, rightTop,
		leftBottom, centerBottom, rightBottom,
		divisibleByError,
	}
	return strings.Join(template, ",")
}

func MakeEven(n int) int {
	if n == 0 {
		return 0
	}
	if n%2 == 0 {
		return n
	} else {
		return n - 1
	}
}

func TrimDotRight(path string) string {
	if !strings.Contains(path, ".") {
		return path
	}
	return strings.Split(path, ".")[0]
}
