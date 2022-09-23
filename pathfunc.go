package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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
				// 이전에 수집된 경로가 존재하지 않으면 처리되는 코드
				item := ScanPlate{
					Searchpath: searchpath,
					Dir:        filepath.Dir(path),
					Base:       filepath.Base(key),
					Ext:        ext,
					Length:     1,
					FrameIn:    num,
					FrameOut:   num,
					RenderIn:   num,
					ConvertExt: ext,
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
