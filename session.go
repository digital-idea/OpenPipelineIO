package main

import (
	"encoding/base64"
	"net/http"
	"time"
)

// SetSessionID 는 SessionID를 가지고 온다.
func SetSessionID(w http.ResponseWriter, id string) {
	c := http.Cookie{
		Name:    "session",
		Value:   base64.StdEncoding.EncodeToString([]byte(id)),
		Expires: time.Now().Add(time.Duration(*flagCookieAge) * time.Hour),
	}
	http.SetCookie(w, &c)
}

// GetSessionID 는 SessionID를 가지고 온다.
func GetSessionID(r *http.Request) string {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "session" {
			decoded, err := base64.StdEncoding.DecodeString(cookie.Value)
			if err != nil {
				return ""
			}
			return string(decoded)
		}
	}
	return ""
}

// RmSessionID 는 SessionID를 제거한다.
func RmSessionID(w http.ResponseWriter) {
	c := http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &c)
}
