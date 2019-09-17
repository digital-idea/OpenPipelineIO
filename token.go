package main

import (
	"errors"
	"net/http"
	"strings"

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

// TokenHandler 함수는 토큰으로 restAPI를 사용할 수 있는지 체크하고 아이디와 엑세스 레벨을 반환한다.
func TokenHandler(r *http.Request, session *mgo.Session) (string, AccessLevel, error) {
	key, err := GetTokenFromHeader(r)
	if err != nil {
		return "", 0, err
	}
	token, err := validToken(session, key)
	if err != nil {
		return "", 0, err
	}
	if token.AccessLevel < 2 {
		return "", 0, errors.New("Insufficient authority levels")
	}
	return token.ID, token.AccessLevel, nil
}
