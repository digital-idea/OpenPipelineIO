package main

import (
	"net/http"
)

// handleErrorCaptcha 핸들러는 Captcha 코드가 다를 때 출력하는 웹페이지이다.
func handleErrorCaptcha(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "error-captcha", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
