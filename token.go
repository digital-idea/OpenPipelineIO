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

// TokenHandler 함수는 토큰으로 restAPI를 사용할 수 있는지 체크하는 함수이다.
func TokenHandler(r *http.Request, session *mgo.Session) error {
	if !*flagAuthmode {
		return nil
	}
	key, err := GetTokenFromHeader(r)
	if err != nil {
		return err
	}
	token, err := validToken(session, key)
	if err != nil {
		return err
	}
	if token.AccessLevel < 2 {
		return errors.New("Insufficient authority levels")
	}
	return nil
}
