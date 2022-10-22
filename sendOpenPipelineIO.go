package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
)

func SendItemOpenPipelineIO(scan ScanPlate, item Item) error {
	// 아이템을 전송한다.
	pbytes, err := json.Marshal(item)
	if err != nil {
		return err
	}
	buff := bytes.NewBuffer(pbytes)

	req, err := http.NewRequest("POST", scan.DNS+"/api/item", buff)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+scan.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		println(str)
		return err
	}
	return nil
}

func SendImageOpenPipelineIO(scan ScanPlate, item Item, path string) error {
	args := []string{
		"-X",
		"POST",
		"-H",
		"Authorization: Basic " + scan.Token,
		"-F",
		"project=" + item.Project,
		"-F",
		"name=" + item.Name,
		"-F",
		"image=@" + path,
		scan.DNS + "/api/uploadthumbnail",
	}

	err := exec.Command(CachedAdminSetting.Curl, args...).Run()
	if err != nil {
		return err
	}

	return nil
}
