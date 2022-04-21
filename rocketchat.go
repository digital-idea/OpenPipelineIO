package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type HookMessage struct {
	Text        string           `json:"text"`
	Channel     string           `json:"channel"` // 로켓챗은 #Chennel, @UserID 를 사용할 수 있다.
	Attachments []HookAttachment `json:"attachments"`
}

type HookAttachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	ImageURL  string `json:"image_url"`
	Color     string `json:"color"`
}

type HookResponse struct {
	Success bool `json:"success"`
}

func (msg *HookMessage) SendRocketChat() (*HookResponse, error) {

	// Check URL Validate
	_, err := url.ParseRequestURI(CachedAdminSetting.RocketChatWebHookURL)
	if err != nil {
		return nil, err
	}

	// Check Token Validate
	if !regexpRocketChatToken.MatchString(CachedAdminSetting.RocketChatToken) {
		return nil, errors.New("check rocketchat token string")
	}

	opt, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/hooks/%s", CachedAdminSetting.RocketChatWebHookURL, CachedAdminSetting.RocketChatToken)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(opt))
	if err != nil {

		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("no such host " + url)
	}
	defer resp.Body.Close()
	res := HookResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}
