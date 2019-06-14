// 현재는 dilog에 로그를 남기지만,
// 추후 회사에서 CSI에 다른 로그서버를 사용한다면 이 파일을 편집할것
package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// setlog는 restAPI를 이용해서 dilog서버에 로그를 남긴다.
// 프로젝트, slug, 사용자명, log내용을 입력받아서..
// (POST에 대한 결과, 에러)를 반환한다.
func setlog(user, project, name, typ, log string) (string, error) {
	v := url.Values{}
	v.Set("tool", "csi3")
	v.Set("keep", "180") // MPAA 기준 로그보관일은 180일이다.
	v.Set("user", user)  // 유저자료구조가 들어가야 본격적으로 이 부분은 사용할 수 있다.
	v.Set("project", project)
	v.Set("slug", name+"_"+typ)
	v.Set("log", log)
	s := v.Encode()
	// endPoint 변수는 dilog를 남기기 위한 restAPI endpoint이다.
	endPoint := *flagDILOG + "/api/setlog"
	req, err := http.NewRequest("POST", endPoint, strings.NewReader(s))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
