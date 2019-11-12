package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// userTemppath 함수는 id를 받아서 각 유저별 Temp 경로를 반환한다.
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
