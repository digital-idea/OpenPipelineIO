package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// userTemppath 함수는 id를 받아서 각 유저별 Temp 경로를 생성 반환한다.
func userTemppath(id string) (string, error) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	path := filepath.Dir(tempDir) + "/" + id
	err = os.MkdirAll(path, 0777)
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