package main

import (
	"net/http"
	"time"
)

// SetSessionID 는 SessionID를 가지고 온다.
func SetSessionID(w http.ResponseWriter, id string) {
	c := http.Cookie{
		Name:    "session",
		Value:   id,
		Expires: time.Now().Add(time.Duration(*flagCookieAge) * time.Hour),
	}
	http.SetCookie(w, &c)
}

// GetSessionID 는 SessionID를 가지고 온다.
func GetSessionID(r *http.Request) string {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "session" {
			return cookie.Value
		}
	}
	return ""
}

// ClearSessionID 는 SessionID를 제거한다.
func ClearSessionID(w http.ResponseWriter) {
	c := http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &c)
}
