package main

import (
	"log"

	"github.com/ashwanthkumar/slack-go-webhook"
	"gopkg.in/mgo.v2"
)

func slacklog(session *mgo.Session, project, logString string) error {
	p, err := getProject(session, project)
	if err != nil {
		return err
	}
	if p.SlackWebhookURL != "" {
		payload := slack.Payload{
			Text:    logString,
			Channel: "#" + project,
		}
		err := slack.Send(p.SlackWebhookURL, "", payload)
		if len(err) > 0 {
			for _, e := range err {
				log.Println(e) // slack에서 생성되는 로그는 출력만 한다. 자체 서비스가 아니기 때문에 프로세스에 영향을 주지 않는다.
			}
		}
	}
	return nil
}
