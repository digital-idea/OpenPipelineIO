package main

import (
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

// GetTokenFromHeader 함수는 사용자가 전달한 Token 값을 가지고 온다.
func GetTokenFromHeader(r *http.Request) (string, error) {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return "", errors.New("authorization failed")
	}
	return auth[1], nil
}

// legacy
// TokenHandler 함수는 토큰으로 restAPI를 사용할 수 있는지 체크하고 아이디와 엑세스 레벨을 반환한다.
func TokenHandler(r *http.Request, session *mgo.Session) (string, AccessLevel, error) {
	key, err := GetTokenFromHeader(r)
	if err != nil {
		return "unknown", UnknownAccessLevel, err
	}
	token, err := validToken(session, key)
	if err != nil {
		return "unknown", UnknownAccessLevel, err
	}
	if token.AccessLevel < 2 {
		return token.ID, token.AccessLevel, errors.New("Insufficient authority levels")
	}
	return token.ID, token.AccessLevel, nil
}

// TokenHandlerV2 함수는 토큰으로 restAPI를 사용할 수 있는지 체크하고 아이디와 엑세스 레벨을 반환한다.
// mgo에서 mongo-driver로 바꾼 version이다.
func TokenHandlerV2(r *http.Request, client *mongo.Client) (string, AccessLevel, error) {
	key, err := GetTokenFromHeader(r)
	if err != nil {
		return "unknown", UnknownAccessLevel, err
	}
	token, err := validTokenV2(client, key)
	if err != nil {
		return "unknown", UnknownAccessLevel, err
	}
	if token.AccessLevel < 2 {
		return token.ID, token.AccessLevel, errors.New("Insufficient authority levels")
	}
	return token.ID, token.AccessLevel, nil
}
