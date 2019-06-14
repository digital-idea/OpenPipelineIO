package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

// setelliteCsvKey2dbKey 함수는 입력받은 문자에서 알파벳과 숫자만 남긴다.
func setelliteCsvKey2dbKey(str string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(str, "")
}

// CSV 경로를 이용해서 Setellite자료구조 리스트를생성한다.
func csv2SetelliteList(csvpath string) []Setellite {
	csvfile, err := os.Open(csvpath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvfile.Close()
	var result []Setellite

	scanner := bufio.NewScanner(csvfile)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	// title
	csvkeys := strings.Split(scanner.Text(), ",")
	var keys []string
	for _, k := range csvkeys {
		keys = append(keys, setelliteCsvKey2dbKey(k))
	}
	for scanner.Scan() {
		var item Setellite
		var notes []string
		for i, value := range strings.Split(scanner.Text(), ",") {
			key := keys[i]
			if key == "Notes" && value != "" {
				notes = append(notes, value)
			} else {
				setSetellite(&item, key, value)
			}
		}
		// Setellite는 여러정보값이 있다면, ";"로 처리한다.
		setSetellite(&item, "Notes", strings.Join(notes, ";"))
		item.setID()
		result = append(result, item)
	}
	return result
}
